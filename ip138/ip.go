package ip138

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	URL = "https://api.ip138.com"
)

type Ip138 struct {
	URL   string
	Token string
	cli   *http.Client
}

// {中国 四川 成都 电信 610000 028}
type LocationInfo struct {
	IP      string
	Country string
	Region  string
	City    string
	Isp     string
	Zip     string
	Zone    string
}

func (l *LocationInfo) GetIP() string      { return l.IP }
func (l *LocationInfo) GetCountry() string { return l.Country }
func (l *LocationInfo) GetRegion() string  { return l.Region }
func (l *LocationInfo) GetCity() string    { return l.City }
func (l *LocationInfo) GetIsp() string     { return l.Isp }
func (l *LocationInfo) GetZip() string     { return l.Zip }
func (l *LocationInfo) GetZone() string    { return l.Zone }

//json struct
type JsonInfo struct {
	Ret  string    `json:"ret"`
	Ip   string    `json:"ip"`
	Data [6]string `json:"data"`
}

func (i *JsonInfo) LocationInfo() LocationInfo {
	return LocationInfo{
		IP:      i.Ip,
		Country: i.Data[0],
		Region:  i.Data[1],
		City:    i.Data[2],
		Isp:     i.Data[3],
		Zip:     i.Data[4],
		Zone:    i.Data[5],
	}
}

func NewIP138(token string) *Ip138 {
	i := &Ip138{
		URL:   URL,
		Token: token,
	}
	i.cli = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	return i
}

func (i *Ip138) IpLocation(ip string) (*LocationInfo, error) {
	return i.IpLocationWithContext(context.Background(), ip)
}

func (i *Ip138) IpLocationWithContext(ctx context.Context, ip string) (*LocationInfo, error) {
	queryUrl := fmt.Sprintf("%s/query/?ip=%s&datatype=%s", URL, ip, "jsonp")
	reqest, err := http.NewRequestWithContext(ctx, "GET", queryUrl, nil)
	if err != nil {
		return nil, err
	}
	reqest.Header.Add("token", i.Token)
	response, err := i.cli.Do(reqest)
	if err != nil {
		return nil, fmt.Errorf("request failed:%s", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %d", response.StatusCode)
	}
	bodyByte, _ := ioutil.ReadAll(response.Body)
	var info JsonInfo
	err = json.Unmarshal(bodyByte, &info)
	if err != nil {
		return nil, fmt.Errorf("json decode failed: %s", err)
	}
	l := info.LocationInfo()
	return &l, nil
}
