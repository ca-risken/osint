package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type privateExposeAPI interface {
	run(string, uint32) (*[]host, error)
}

type privateExposeClient struct {
	config privateExposeConfig
}

type privateExposeConfig struct {
	ResultPath    string `required:"true" split_words:"true"`
	HarvesterPath string `required:"true" split_words:"true"`
}

func newPrivateExposeClient() *privateExposeClient {
	var conf privateExposeConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	return &privateExposeClient{config: conf}
}

func (h *privateExposeClient) run(domain string, relAlertFindingID uint32) (*[]host, error) {
	now := time.Now().Unix()
	filePath := fmt.Sprintf("%s/%v_%v.xml", h.config.ResultPath, relAlertFindingID, now)
	harvesterPath := fmt.Sprintf("%s/theHarvester.py", h.config.HarvesterPath)
	cmd := exec.Command("python3", harvesterPath, "-d", domain, "-b", "all", "-f", filePath)
	cmd.Dir = h.config.HarvesterPath
	err := cmd.Run()
	if err != nil {
		appLogger.Errorf("Failed exec theHarvester. error: %v", err)
		return nil, nil
	}

	bytes, err := readAndDeleteFile(filePath)
	if err != nil {
		return nil, err
	}
	harvesterStruct := theHarvester{}
	xml.Unmarshal(bytes, &harvesterStruct)
	if err != nil {
		return nil, err
	}
	return &harvesterStruct.Hosts, nil
}

func tmpRun() (*[]host, error) {
	harvesterStruct := theHarvester{}
	bytes, err := readFile("/tmp/1001_1603703153.xml")
	if err != nil {
		return nil, err
	}
	xml.Unmarshal(bytes, &harvesterStruct)
	if err != nil {
		return nil, err
	}

	return &harvesterStruct.Hosts, nil
}

func detectHost(hosts *[]host, domain string, detectList *[]string) *[]string {
	retList := []string{}

	for _, host := range *hosts {
		if !zero.IsZeroVal(host.IP) && !zero.IsZeroVal(host.HostName) && isDetected(host.HostName, detectList) && strings.Index(host.HostName, "."+domain) > -1 {
			retList = append(retList, host.HostName)
		}
	}
	return &retList
}

func getHTTPStatus(detectedHosts *[]string) *[]httpStatus {
	retList := []httpStatus{}

	for _, host := range *detectedHosts {
		http, urlHTTP, _ := getStatusCodeHost(host, "http")
		https, urlHTTPS, _ := getStatusCodeHost(host, "https")
		if !zero.IsZeroVal(http) && !zero.IsZeroVal(https) {
			retList = append(retList, httpStatus{HostName: host, HTTP: http, UrlHTTP: urlHTTP, HTTPS: https, UrlHTTPS: urlHTTPS})
		}
	}
	return &retList
}

func getStatusCodeHost(host, protocol string) (int, string, error) {
	url := fmt.Sprintf("%s://%s", protocol, host)
	req, _ := http.NewRequest("GET", url, nil)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)

	// Timeoutもエラーに入るので、特にログも出さないでスルー(ドメインを見つけてもHTTPで使われているとは限らないため)
	if err != nil {
		return 0, "", err
	}

	defer res.Body.Close()
	return res.StatusCode, res.Request.URL.String(), nil
}

func isDetected(host string, detectList *[]string) bool {
	for _, detectWord := range *detectList {
		if strings.Index(host, detectWord) > -1 {
			return true
		}
	}
	return false
}

func readAndDeleteFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(fileName); err != nil {
		return nil, err
	}
	if err := os.Remove(fileName + ".html"); err != nil {
		return nil, err
	}
	return bytes, nil
}

func readFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

type theHarvester struct {
	Hosts []host `xml:"host"`
}

type host struct {
	IP       string `xml:"ip"`
	HostName string `xml:"hostname"`
}

type httpStatus struct {
	HostName string `json:"hostname"`
	HTTP     int    `json:"http"`
	UrlHTTP  string `json:"url_http"`
	HTTPS    int    `json:"https"`
	UrlHTTPS string `json:"url_https"`
}
