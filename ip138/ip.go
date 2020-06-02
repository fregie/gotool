package ip138

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	URL = "https://api.ip138.com/query/"
)

type Ip138 struct {
	URL   string
	Token string
	cli   *http.Client
}

// {中国 四川 成都 电信 610000 028}
type LocationInfo struct {
	Country string
	Region  string
	City    string
	Isp     string
	Zip     string
	Zone    string
}

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

func (i JsonInfo) LocationInfo() LocationInfo {
	return LocationInfo{
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
	i.cli = &http.Client{Transport: &http.Transport{}}
	return i
}

func (i *Ip138) IpLocation(ip string) (*LocationInfo, error) {
	queryUrl := fmt.Sprintf("%s?ip=%s&datatype=%s", URL, ip, "jsonp")
	reqest, err := http.NewRequest("GET", queryUrl, nil)
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
