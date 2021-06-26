package metrics

import (
	"flag"
	"metrics/ping"
	"net/http"

	log "github.com/gogap/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 定义 Exporter 的版本（Version）、监听地址（listenAddress）、采集 url（metricPath）以及⾸⻚（landingPage）
var (
	Version       = "  v0.1.0"
	listenAddress = flag.String("web.listen-address", ":8000", "Address to listen on for web interface and telemetry.")
	metricPath    = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	landingPage   = []byte("<html><head><title>Example Exporter" + Version + "</title></head><body><h1>Example Exporter" + Version + "</h1><p><ahref='" + *metricPath + "'>Metrics</a></p></body></html>")
)

// 定义 Exporter 结构体
type Exporter struct {
	error        prometheus.Gauge
	scrapeErrors *prometheus.CounterVec
}

// 定义结构体实例化的函数 NewExporter
func NewExporter() *Exporter {
	return &Exporter{}
}

// Describe 函数，传递指标描述符到 channel，这个函数不⽤动，直接使⽤即可，⽤来⽣成采集指标的描述信息。
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})
	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()
	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
}

//collect 函数，采集数据的⼊⼝
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	var err error
	// 每个指标值的采集逻辑，在对应的采集函数中
	if err = ScrapePing(ch); err != nil {
		e.scrapeErrors.WithLabelValues("ping").Inc()
	}
	/*
		if err = ScrapeDisk(ch); err != nil {
			e.scrapeErrors.WithLabelValues("disk").Inc()
		}*/
}

// Collect 函数将执⾏抓取函数并返回数据，返回的数据传递到 channel 中
// 并且传递的同时绑定原先的指标描述符，以及指标的类型（Guage）
// 需要将所有的指标获取函数在这⾥写⼊。

// 指标仅有单条数据，不带维度信息示例如下：
func ScrapePing(ch chan<- prometheus.Metric) error {
	// 指标获取逻辑，此处不做具体操作，仅仅赋值进⾏示例
	//mem_info, _ := mem.VirtualMemory()
	var metrics prometheus.Metric
	infos := ping.ReturnInfos()
	for _, info := range infos {
		if info.Type == "not_found" {
			log.Error("Error to scrape,beacuse domain type not found...")
		}
		// ⽣成 NewDesc 类型的数据格式，该指标⽆维度，[] string {} 为空
		new_desc := prometheus.NewDesc(info.MetricsName, info.Help, info.LabelsKey, nil)
		// ⽣成具体的采集信息并写⼊ ch 通道
		metrics = prometheus.MustNewConstMetric(new_desc,
			info.MetricsType, float64(info.ErrorsTotal), info.Domain, info.Type)

		ch <- metrics
		log.Info("Succeed to scrape metrics ", info.Domain, " total ping is ", info.Total)
	}

	return nil
}

/*
// 指标有多条数据，带维度信息示例如下：
func ScrapeDisk(ch chan<- prometheus.Metric) error {
	fs, _ := disk.Partitions(false)
	for _, val := range fs {
		d, _ := disk.Usage(val.Mountpoint)
		metric_name := prometheus.BuildFQName("sys", "", "disk_size")
		new_desc := prometheus.NewDesc(metric_name, "Gauge metric with disk_usage", []string{"mountpoint"}, nil)
		metric_mes := prometheus.MustNewConstMetric(new_desc,
			prometheus.GaugeValue, float64(d.UsedPercent), val.Mountpoint)
		ch <- metric_mes
	}
	return nil
}*/

func Start() {
	// 解析定义的监听端⼝等信息
	flag.Parse()
	// ⽣成⼀个 Exporter 类型的对象，该 exporter 需具有 collect 和 Describe ⽅法
	exporter := NewExporter()
	// 将 exporter 注册⼊ prometheus，prometheus 将定期从 exporter 拉取数据
	prometheus.MustRegister(exporter)
	// 接收 http 请求时，触发 collect 函数，采集数据
	http.Handle(*metricPath, promhttp.Handler())
	// Exporter 界面
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})
	log.Error(http.ListenAndServe(":8080", nil))
}
