package marathon

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type mesosDNSService struct {
	Service string `json:"service"`
	Host    string `json:"host"`
	IP      string `json:"ip"`
	Port    string `json:"port"`
}

func mesosDNSLookup(resolverURL string, serviceName string) (*mesosDNSService, error) {
	url := resolverURL + "/v1/services/_api._" + serviceName + "._tcp.marathon.mesos"

	req, err := http.Get(url)
	if nil != err {
		return nil, err
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if nil != err {
		return nil, err
	}

	var data []mesosDNSService

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		return nil, err
	}

	if len(data) < 1 {
		err = errors.New("Lookup returned zero services")
		return nil, err
	}

	if len(data) > 1 {
		log.Println("[WARN] Lookup returned multiple services, expected 1, got: " + strconv.Itoa(len(data)))
	}

	if len(data[0].Service) == 0 {
		err = errors.New("Lookup returned empty service")
		return nil, err
	}

	return &data[0], nil
}

func mesosDNSHostPort(resolverURL string, serviceName string) (string, error) {
	service, err := mesosDNSLookup(resolverURL, serviceName)
	if err != nil {
		return "", err
	}
	hostport := service.IP + ":" + service.Port
	return hostport, nil
}
