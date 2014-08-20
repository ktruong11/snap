package collection

import (
	"encoding/json"
	"os/exec"
	"time"
)

var (
	facter_cache metricCache = newMetricCache()
)

type facterCollector struct {
	Collector
	FacterPath string
}

type fact interface {}

func NewFacterCollector(fp string, caching bool, cache_ttl float64) collector{
	c := new(facterCollector)
	c.FacterPath = fp
	c.Caching = caching
	c.CachingTTL = cache_ttl
	return c
}

func (f *facterCollector) GetMetricList() []Metric{
	if !f.Caching || (f.Caching && facter_cache.IsExpired(f.CachingTTL)) {
		out, _ := exec.Command("sh", "-c", f.FacterPath+" -j").Output()
		update_time := time.Now()
		facter_map := new(map[string]metricType)
		json.Unmarshal(out, facter_map)

		hostname := (*facter_map)["fqdn"].(string)

		metric := Metric{hostname, []string{"facter", "facts"}, update_time, map[string]metricType{}, "facter", Polling}
		for k, v := range *facter_map {
			metric.Values[k] = v
		}
		facter_cache.Metrics = []Metric{metric}
		facter_cache.LastPull = time.Now()
		facter_cache.New = false
	}

	return facter_cache.Metrics
}

func (f *facterCollector) GetMetricValues(metrics []Metric, things...interface {}) []Metric{
	return f.GetMetricList()
}