package ping

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/gogap/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

var domains = []string{
	"www.baidu.com",
	"www.sohu.com",
	"www.google.com",
}

const (
	Baidu  = "www.baidu.com"
	Sohu   = "www.sohu.com"
	Google = "www.google.com"
)

type PingInfo struct {
	Domain       string
	Type         string
	Total        int
	SucceedTotal int
	ErrorsTotal  int

	LabelsKey   []string
	Help        string
	MetricsName string
	MetricsType prometheus.ValueType
}

var infos []PingInfo

// Prometheus Sql
// 过去一分钟探测失败总数: (friendly_u_ping_proxy_error_total offset 1m)-(friendly_u_ping_proxy_error_total offset 2m)
func Start() {
	var wg sync.WaitGroup

	infos = make([]PingInfo, len(domains))
	for k, _ := range infos {
		infos[k].Domain = domains[k]
		initInfo(k)
		fmt.Println(k, infos[k])
	}

	for {
		wg.Add(len(domains))
		for k, _ := range infos {
			go checkServer(&infos[k])
			wg.Done()
		}
		wg.Wait()
		time.Sleep(5 * time.Second)
	}
}

func checkServer(info *PingInfo) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get("http://" + info.Domain)
	info.Total++
	if err != nil {
		info.ErrorsTotal++
		log.Error("Error to check ", info.Domain, ",error: ", err)
		return
	}

	log.Info("Succeed to check ", info.Domain)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		info.SucceedTotal++
	} else {
		info.ErrorsTotal++
	}
}

func ReturnInfos() []PingInfo {
	return infos
}

func initInfo(k int) {
	if infos[k].Domain == Baidu ||
		infos[k].Domain == Sohu {
		infos[k].Type = "nat"
		infos[k].Help = "Gauge metric with nat64"
		infos[k].MetricsName = "friendly_u_ping_nat64_error_total"
		infos[k].MetricsType = prometheus.GaugeValue
		infos[k].LabelsKey = []string{"host", "type"}
	} else if infos[k].Domain == Google {
		infos[k].Type = "proxy"
		infos[k].Help = "Gauge metric with proxy"
		infos[k].MetricsName = "friendly_u_ping_proxy_error_total"
		infos[k].MetricsType = prometheus.GaugeValue
		infos[k].LabelsKey = []string{"host", "type"}
	} else {
		infos[k].Type = "not_found"
	}
}
