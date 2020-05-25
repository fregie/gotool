package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
)

// 服务信息
type ServiceInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type Service struct {
	ServiceInfo   ServiceInfo
	LeaseDruation time.Duration
	Client        *clientv3.Client
	stop          chan error
	leaseID       clientv3.LeaseID
}

// NewService 创建一个注册服务
func NewService(info ServiceInfo, client *clientv3.Client) *Service {
	service := &Service{
		ServiceInfo: info,
		Client:      client,
	}
	return service
}

// Start 注册服务启动
func (service *Service) Start() (err error) {
	for {
		ch, err := service.keepAlive()
		if err != nil {
			log.Fatal(err)
			return err
		}

		for {
			select {
			case err := <-service.stop:
				return err
			case <-service.Client.Ctx().Done():
				return errors.New("service closed")
			case _, ok := <-ch:
				// 监听租约
				if !ok {
					log.Println("keep alive channel closed")
					// service.revoke()
					ch, err = service.keepAlive()
					if err != nil {
						log.Print(err)
					}
					time.Sleep(time.Second)
				}
				// log.Printf("Recv reply from service: %s, ttl:%d", service.getKey(), resp.TTL)
			}
		}
	}

}

func (service *Service) Stop() {
	service.stop <- nil
}

func (service *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &service.ServiceInfo
	key := service.getKey()
	val, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}

	// 创建一个租约
	resp, err := service.Client.Grant(context.TODO(), int64(service.LeaseDruation.Seconds()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = service.Client.Put(context.TODO(), key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	service.leaseID = resp.ID
	return service.Client.KeepAlive(context.TODO(), resp.ID)
}

func (service *Service) revoke() error {
	_, err := service.Client.Revoke(context.TODO(), service.leaseID)
	if err != nil {
		return err
	}
	log.Printf("servide:%s stop\n", service.getKey())
	return err
}

func (service *Service) getKey() string {
	return fmt.Sprintf("/service/%s/%s:%d", service.ServiceInfo.Name, service.ServiceInfo.IP, service.ServiceInfo.Port)
}

func GetLocalIP(infName string) (string, error) {
	inf, err := net.InterfaceByName(infName)
	if err != nil {
		return "", err
	}
	addrs, err := inf.Addrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("interface have no addr")
}
