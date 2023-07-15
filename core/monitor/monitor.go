package monitor

import (
	"fmt"
	"strings"
	"time"

	"github.com/nyelonong/boilerplate-go/core/monitor/interfaces"
	"github.com/nyelonong/boilerplate-go/core/monitor/prometheus"
)

const (
	EnginePrometheus = "prometheus"
)

type Config struct {
	Engine       string
	Prefix       string
	DefaultLabel map[string]string
}

type Client struct {
	config Config
	engine interfaces.Engine
}

var globalClient *Client

func Init(conf Config) error {
	// Currently only Prometheus available
	if conf.Engine != EnginePrometheus {
		return fmt.Errorf("engine %s not found", conf.Engine)
	}

	if !strings.HasSuffix(conf.Prefix, "_") {
		conf.Prefix = conf.Prefix + "_"
	}

	promeClient := prometheus.New()
	globalClient = &Client{
		config: conf,
		engine: promeClient,
	}

	return nil
}

func addPrefix(name string) string {
	return globalClient.config.Prefix + name
}

func addLabel(labels []string) []string {
	for k, _ := range globalClient.config.DefaultLabel {
		// Check Duplicate Keys
		exist := false
		for _, v := range labels {
			if k == v {
				exist = true
			}
		}

		if !exist {
			labels = append(labels, k)
		}
	}
	return labels
}

func addParam(param map[string]string) map[string]string {
	for key, value := range globalClient.config.DefaultLabel {
		if _, ok := param[key]; !ok {
			param[key] = value
		}
	}
	return param
}

func NewCounter(name string, labels []string) {
	globalClient.engine.NewCounter(addPrefix(name), addLabel(labels))
}

func Increment(name string, param map[string]string) error {
	return globalClient.engine.Increment(addPrefix(name), addParam(param))
}

func Add(name string, number float64, param map[string]string) error {
	return globalClient.engine.Add(addPrefix(name), number, addParam(param))
}

func NewGauge(name string, labels []string) {
	globalClient.engine.NewGauge(addPrefix(name), addLabel(labels))
}

func Set(name string, number float64, param map[string]string) error {
	return globalClient.engine.Set(addPrefix(name), number, addParam(param))
}

func NewHistogram(name string, labels []string) {
	globalClient.engine.NewHistogram(addPrefix(name), addLabel(labels))
}

func Record(name string, start time.Time, param map[string]string) error {
	diff := time.Since(start).Milliseconds()
	return globalClient.engine.Record(addPrefix(name), float64(diff), addParam(param))
}
