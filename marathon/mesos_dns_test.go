package marathon

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type (
	dnsResponses struct {
		jsonResponse string
		service      string
		host         string
		ip           string
		port         string
		hostPort     string
		valid        bool
	}
)

var (
	dnsTestData = []dnsResponses{
		{`[
                {
                   "service": "_api._kafka._tcp.marathon.mesos",
                   "host": "kafka-94yoo-s30.marathon.mesos.",
                   "ip": "10.0.6.180",
                   "port": "31567"
                }
            ]`, "_api._kafka._tcp.marathon.mesos", "kafka-94yoo-s30.marathon.mesos.", "10.0.6.180", "31567", "10.0.6.180:31567", true},
		{`[
                {
                   "service": "",
                   "host": "",
                   "ip": "",
                   "port": ""
                }
            ]`, "", "", "", "", "", false},
		{`[
                {
                   "service": "_api._kafka._tcp.marathon.mesos",
                   "host": "kafka-94yoo-s30.marathon.mesos.",
                   "ip": "10.0.6.180",
                   "port": "31567"
               },
               {
                  "service": "",
                  "host": "",
                  "ip": "",
                  "port": ""
               }
            ]`, "_api._kafka._tcp.marathon.mesos", "kafka-94yoo-s30.marathon.mesos.", "10.0.6.180", "31567", "10.0.6.180:31567", true},
		{`[]`, "", "", "", "", "", false},
	}
)

func Test_MesosDNSLoadFromServer(t *testing.T) {

	for _, dnsRequest := range dnsTestData {
		// serve that json chunk from a local loopback webserver
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, dnsRequest.jsonResponse)
		}))
		defer ts.Close()

		if dnsRequest.valid {
			service, err := mesosDNSLookup(ts.URL, "foo")
			if err != nil {
				t.Errorf("loading from dns: %v", err)
			}
			if !reflect.DeepEqual(service.Service, dnsRequest.service) {
				t.Errorf("loaded from mesos-dns: got %v, expected: %v", service.Service, dnsRequest.service)
			}
			if !reflect.DeepEqual(service.Host, dnsRequest.host) {
				t.Errorf("loaded from mesos-dns: got %v, expected: %v", service.Host, dnsRequest.host)
			}
			if !reflect.DeepEqual(service.IP, dnsRequest.ip) {
				t.Errorf("loaded from mesos-dns: got %v, expected: %v", service.IP, dnsRequest.ip)
			}
			if !reflect.DeepEqual(service.Port, dnsRequest.port) {
				t.Errorf("loaded from mesos-dns: got %v, expected: %v", service.Port, dnsRequest.port)
			}

			hostPort1, err := mesosDNSHostPort(ts.URL, "foo")
			if err != nil {
				t.Errorf("loading from dns: %v", err)
			}

			hostPort2, err := mesosDNSHostPort(ts.URL, "foo")
			if err != nil {
				t.Errorf("loading repeat from dns: %v", err)
			}

			if !reflect.DeepEqual(hostPort1, dnsRequest.hostPort) {
				t.Errorf("loaded from mesos-dns: got %v, expected: %v", hostPort1, dnsRequest.hostPort)
			}

			if !reflect.DeepEqual(hostPort1, hostPort2) {
				t.Errorf("repeated lookups should match 1: %v, 2: %v", hostPort1, hostPort2)
			}
		} else {
			service, err := mesosDNSLookup(ts.URL, "foo")
			if err == nil {
				t.Errorf("Expected and error got: %v", service)
			}
		}
	}
}
