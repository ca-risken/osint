package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gassara-kys/envconfig"
	"github.com/vikyd/zero"
)

type harvesterConfig struct {
	ResultPath    string `required:"true" split_words:"true" default:"/results"`
	HarvesterPath string `required:"true" split_words:"true" default:"/theHarvester"`
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
	filePath := fmt.Sprintf("%s/%v_%v", h.ResultPath, relAlertFindingID, now)
	harvesterPath := fmt.Sprintf("%s/theHarvester.py", h.HarvesterPath)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "python3", harvesterPath, "-d", domain, "-b", "all", "-f", filePath)
	cmd.Dir = h.HarvesterPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		appLogger.Errorf("Failed to execute theHarvester. error: %v", stderr.String())
		return nil, err
	}

	bytes, err := readAndDeleteFile(filePath)
	if err != nil {
		return nil, err
	}
	hostsWithIP := hostsWithIP{}
	hostsWithoutIP := hostsWithoutIP{}
	if err = xml.Unmarshal(bytes, &hostsWithIP); err != nil {
		appLogger.Errorf("Failed to unmarshal result. error: %v", err)
		return nil, err
	}
	if err = xml.Unmarshal(bytes, &hostsWithoutIP); err != nil {
		appLogger.Errorf("Failed to unmarshal result. error: %v", err)
		return nil, err
	}

	return makeHosts(&hostsWithIP, &hostsWithoutIP, domain), nil
}

func makeHosts(hostsWithIP *hostsWithIP, hostsWithoutIP *hostsWithoutIP, domain string) *[]host {
	arrHost := []host{}
	for _, hostWithIP := range hostsWithIP.Hosts {
		if strings.HasSuffix(hostWithIP.HostName, "."+domain) {
			arrHost = append(arrHost, hostWithIP)
		}
	}
	for _, hostWithoutIP := range hostsWithoutIP.Hosts {
		if strings.HasSuffix(hostWithoutIP, "."+domain) {
			arrHost = append(arrHost, host{IP: getIPAddr(hostWithoutIP), HostName: hostWithoutIP})
		}
	}
	// Add domain
	arrHost = append(arrHost, host{IP: getIPAddr(domain), HostName: domain})
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
	xmlFileName := fileName + ".xml"
	jsonFileName := fileName + ".json"
	bytes, err := ioutil.ReadFile(xmlFileName)
	if err != nil {
		return nil, err
	}
	if err := os.Remove(xmlFileName); err != nil {
		return nil, err
	}
	if err := os.Remove(jsonFileName); err != nil {
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
	Host                  host
	PrivateExpose         privateExpose
	Takeover              takeover
	CertificateExpiration certificateExpiration
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
	return zero.IsZeroVal(h.IP)
}
