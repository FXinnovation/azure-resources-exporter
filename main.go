package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
)

var (
	configFile    = kingpin.Flag("config.file", "Exporter configuration file.").Default("config.yml").String()
	listenAddress = kingpin.Flag("web.listen-address", "The address to listen on for HTTP requests.").Default(":9259").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	config        Config
)

// Config of the exporter
type Config struct {
}

func init() {
	prometheus.MustRegister(version.NewCollector("azure_resources_exporter"))
}

func main() {
	kingpin.Version(version.Print("azure-resources-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.Println("Starting exporter", version.Info())
	log.Println("Build context", version.BuildContext())

	config = loadConfig(*configFile)

	collector := NewCollector()
	prometheus.MustRegister(collector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>azure-resources-exporter</title></head>
			<body>
			<h1>azure-resources-exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	log.Println("Beginning to serve on address ", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func loadConfig(configFile string) Config {
	config := Config{}

	// Load the config from the file
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	errYAML := yaml.Unmarshal([]byte(configData), &config)
	if errYAML != nil {
		log.Fatalf("Error: %v", errYAML)
	}

	return config
}
