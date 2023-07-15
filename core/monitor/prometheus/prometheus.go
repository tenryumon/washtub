package prometheus

import (
	"fmt"

	prome "github.com/prometheus/client_golang/prometheus"

	"github.com/nyelonong/boilerplate-go/core/monitor/interfaces"
)

type Counter struct {
	data   *prome.CounterVec
	labels []string
}

type Gauge struct {
	data   *prome.GaugeVec
	labels []string
}

type Summary struct {
	data   *prome.SummaryVec
	labels []string
}

type Client struct {
	counters  map[string]Counter
	gauges    map[string]Gauge
	summaries map[string]Summary
}

func New() interfaces.Engine {
	return &Client{
		counters:  map[string]Counter{},
		gauges:    map[string]Gauge{},
		summaries: map[string]Summary{},
	}
}

func getLabels(config []string, param map[string]string) prome.Labels {
	labels := prome.Labels{}
	for _, v := range config {
		labels[v] = param[v]
	}
	return labels
}

func (c *Client) NewCounter(name string, labels []string) {
	counter := prome.NewCounterVec(prome.CounterOpts{Name: name}, labels)
	prome.MustRegister(counter)

	c.counters[name] = Counter{data: counter, labels: labels}
}

func (c *Client) Increment(name string, param map[string]string) error {
	counter, ok := c.counters[name]
	if !ok {
		return fmt.Errorf("Metrics %s not found", name)
	}

	labels := getLabels(counter.labels, param)
	promeCounter, err := counter.data.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("Failed GetMetricsWith %s because %s", name, err)
	}

	promeCounter.Inc()
	return nil
}

func (c *Client) Add(name string, number float64, param map[string]string) error {
	counter, ok := c.counters[name]
	if !ok {
		return fmt.Errorf("Metrics %s not found", name)
	}

	labels := getLabels(counter.labels, param)
	promeCounter, err := counter.data.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("Failed GetMetricsWith %s because %s", name, err)
	}

	promeCounter.Add(number)
	return nil
}

func (c *Client) NewGauge(name string, labels []string) {
	gauge := prome.NewGaugeVec(prome.GaugeOpts{Name: name}, labels)
	prome.MustRegister(gauge)

	c.gauges[name] = Gauge{data: gauge, labels: labels}
}

func (c *Client) Set(name string, number float64, param map[string]string) error {
	gauge, ok := c.gauges[name]
	if !ok {
		return fmt.Errorf("Metrics %s not found", name)
	}

	labels := getLabels(gauge.labels, param)
	promeGauge, err := gauge.data.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("Failed GetMetricsWith %s because %s", name, err)
	}

	promeGauge.Set(number)
	return nil
}

func (c *Client) NewHistogram(name string, labels []string) {
	summary := prome.NewSummaryVec(prome.SummaryOpts{
		Name:       name,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01},
	}, labels)
	prome.MustRegister(summary)

	c.summaries[name] = Summary{data: summary, labels: labels}
}

func (c *Client) Record(name string, number float64, param map[string]string) error {
	summary, ok := c.summaries[name]
	if !ok {
		return fmt.Errorf("Metrics %s not found", name)
	}

	labels := getLabels(summary.labels, param)
	promeSummary, err := summary.data.GetMetricWith(labels)
	if err != nil {
		return fmt.Errorf("Failed GetMetricsWith %s because %s", name, err)
	}

	promeSummary.Observe(number)
	return nil
}
