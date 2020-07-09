package ip138

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MobileJsonInfo struct {
	Ret    string    `json:"ret"`
	Mobile string    `json:"mobile"`
	Data   [5]string `json:"data"`
}

type MobileInfo struct {
	Number   string
	Province string
	City     string
	Provider string
	Zone     string
}

func (i *MobileJsonInfo) MobileInfo() MobileInfo {
	return MobileInfo{
		Number:   i.Mobile,
		Province: i.Data[0],
		City:     i.Data[1],
		Provider: i.Data[2],
		Zone:     i.Data[3],
	}
}

func (i *Ip138) Mobile(number string) (*MobileInfo, error) {
	queryUrl := fmt.Sprintf("%s/mobile/?mobile=%s&datatype=%s", URL, number, "jsonp")
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
	var info MobileJsonInfo
	err = json.Unmarshal(bodyByte, &info)
	if err != nil {
		return nil, fmt.Errorf("json decode failed: %s", err)
	}
	m := info.MobileInfo()
	return &m, nil
}
