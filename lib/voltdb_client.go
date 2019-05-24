package lib

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type VoltDBClient struct {
	username      string
	password      string
	databases     []string
	useHttp       bool
	insecureHttps bool
}

type Stats struct {
	database string
	state    []byte
	cpu      []byte
	txns     []byte
	latency  []byte
	ram      []byte
	dr_role  []byte
	dr_state []byte
}

func NewVoltDBClient(user string, pass string, dbs []string, useHttp bool, insecureHttps bool) *VoltDBClient {
	initializeClient(user, pass, dbs, useHttp, insecureHttps)

	return &VoltDBClient{
		username:      user,
		password:      pass,
		databases:     dbs,
		useHttp:       useHttp,
		insecureHttps: insecureHttps,
	}
}

func initializeClient(user string, pass string, addrs []string, useHttp bool, insecureHttps bool) {
	var proto string
	for _, addr := range addrs {
		if !useHttp {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: insecureHttps}
			proto = "https"
		} else {
			proto = "http"
		}

		request := fmt.Sprintf("%s://%s/api/1.0/?Procedure=@Ping&admin=false&User=%s&Password=%s", proto, addr, user, pass)
		resp, err := http.Get(request)
		if err != nil {
			log.Fatal(err)
		} else if resp.StatusCode != http.StatusOK {
			re := regexp.MustCompile(`(Password)=(.*?:)`)
			errMsg := re.ReplaceAllString(addr, `$1=**********`)
			log.Fatal(fmt.Sprintf("Failed to connect to client at %s", errMsg))
		}
	}
	log.Print("Successfully connected to client(s)")
}

func get(path string) ([]byte, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func setPaths(addr string, client *VoltDBClient) []string {
	var proto string
	if client.useHttp {
		proto = "http"
	} else {
		proto = "https"
	}
	basepath := proto + "://%s/api/1.0/?Procedure=%s&Parameters=%s&admin=false&User=%s&Password=%s"
	return []string{
		fmt.Sprintf(basepath, addr, "@SystemInformation", "['OVERVIEW']", client.username, client.password),
		fmt.Sprintf(basepath, addr, "@Statistics", "['CPU',0]", client.username, client.password),
		fmt.Sprintf(basepath, addr, "@Statistics", "['LATENCY',0]", client.username, client.password),
		fmt.Sprintf(basepath, addr, "@Statistics", "['MEMORY',0]", client.username, client.password),
		fmt.Sprintf(basepath, addr, "@Statistics", "['DRROLE']", client.username, client.password),
	}
}

func scrapeData(paths []string) (data [][]byte, err error) {
	for _, path := range paths {
		buf, err := get(path)
		if err != nil {
			return nil, err
		}
		data = append(data, buf)
	}
	return data, nil
}

func getStats(addrs []string, client *VoltDBClient) (*[]Stats, error) {
	var stats []Stats
	for _, addr := range addrs {
		paths := setPaths(addr, client)

		data, err := scrapeData(paths)
		if err != nil {
			return nil, err
		} else {
			stats = append(stats, Stats{
				database: addr,
				state:    data[0],
				cpu:      data[1],
				txns:     data[2],
				latency:  data[2],
				ram:      data[3],
				dr_role:  data[4],
				dr_state: data[4],
			})
		}
	}
	return &stats, nil
}
