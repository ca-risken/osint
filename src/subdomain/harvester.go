package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/vikyd/zero"
)

type harvesterConfig struct {
	ResultPath    string `required:"true" split_words:"true"`
	HarvesterPath string `required:"true" split_words:"true"`
}

func newHarvesterConfig() harvesterConfig {
	var conf harvesterConfig
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}
	return conf
}

func (h *harvesterConfig) run(domain string, relAlertFindingID uint32) (*[]host, error) {
	now := time.Now().Unix()
	filePath := fmt.Sprintf("%s/%v_%v.xml", h.ResultPath, relAlertFindingID, now)
	harvesterPath := fmt.Sprintf("%s/theHarvester.py", h.HarvesterPath)
	cmd := exec.Command("python3", harvesterPath, "-d", domain, "-b", "all", "-f", filePath)
	cmd.Dir = h.HarvesterPath
	err := cmd.Run()
	if err != nil {
		appLogger.Errorf("Failed exec theHarvester. error: %v", err)
		return nil, err
	}

	bytes, err := readAndDeleteFile(filePath)
	if err != nil {
		return nil, err
	}
	hostsWithIP := hostsWithIP{}
	hostsWithoutIP := hostsWithoutIP{}
	xml.Unmarshal(bytes, &hostsWithIP)
	xml.Unmarshal(bytes, &hostsWithoutIP)

	return makeHosts(&hostsWithIP, &hostsWithoutIP, domain), nil
}

func tmpRun() (*[]host, error) {
	hostsWithIP := hostsWithIP{}
	hostsWithoutIP := hostsWithoutIP{}
	bytes, err := tmpReadFile("/tmp/bbb.xml")
	if err != nil {
		return nil, err
	}
	xml.Unmarshal(bytes, &hostsWithIP)
	xml.Unmarshal(bytes, &hostsWithoutIP)

	return makeHosts(&hostsWithIP, &hostsWithoutIP, ""), nil
}

func makeHosts(hostsWithIP *hostsWithIP, hostsWithoutIP *hostsWithoutIP, domain string) *[]host {
	arrHost := []host{}
	for _, hostWithIP := range hostsWithIP.Hosts {
		if strings.Index(hostWithIP.HostName, "."+domain) > -1 {
			arrHost = append(arrHost, hostWithIP)
		}
	}
	for _, hostWithoutIP := range hostsWithoutIP.Hosts {
		if strings.Index(hostWithoutIP, "."+domain) > -1 {
			arrHost = append(arrHost, host{IP: getIPAddr(hostWithoutIP), HostName: hostWithoutIP})
		}
	}
	ret := sliceUnique(&arrHost)
	return &ret
}

func getIPAddr(domain string) string {
	ips, _ := net.LookupIP(domain)
	for _, ip := range ips {
		return ip.String()
	}
	return ""
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

func tmpReadFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func sliceUnique(target *[]host) []host {
	m := map[host]struct{}{}
	ret := []host{}
	for _, t := range *target {
		if _, ok := m[t]; !ok {
			ret = append(ret, t)
			m[t] = struct{}{}
		}
	}
	return ret
}

type osintResults struct {
	OsintResults *[]osintResult
}
type osintResult struct {
	Host          host
	PrivateExpose privateExpose
	Takeover      takeover
}

type hostsWithIP struct {
	Hosts []host `xml:"host"`
}

type hostsWithoutIP struct {
	Hosts []string `xml:"host"`
}

type host struct {
	IP       string `xml:"ip"`
	HostName string `xml:"hostname"`
}

func (h *host) isDown() bool {
	if zero.IsZeroVal(h.IP) {
		return true
	}
	return false
}
