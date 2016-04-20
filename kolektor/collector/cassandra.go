package collector

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/ezeev/golang/kolektor/output"
)

// CassandraCollector Struct for the apache Cassandra collector. The Cassandra collector depends on golokia
type CassandraCollector struct {
	LastCollection time.Time
	Config         CassandraConfig
}

//CassandraConfig this is some config
type CassandraConfig struct {
	Collector string
	Host      string
	Interval  float64
	Metrics   []string
	Tags      map[string]string
}

// NewCassandraCollector new CassandraCollector
func NewCassandraCollector(config CassandraConfig) (*CassandraCollector, error) {
	c := &CassandraCollector{Config: config}
	return c, nil
}

// Interval get interval o
func (c *CassandraCollector) Interval() float64 {
	return c.Config.Interval
}

// Name Returns the name of the collector
func (c *CassandraCollector) Name() string {
	return c.Config.Collector
}

// SetLastCollectionTime Sets the last collection time of the current collector
func (c *CassandraCollector) SetLastCollectionTime(t time.Time) {
	c.LastCollection = t
}

// LastCollectionTime Returns the last collection time of the last collector
func (c *CassandraCollector) LastCollectionTime() time.Time {
	return c.LastCollection
}

// Collect Collects metrics from Cassandra and returns a slice of Metrics
func (c *CassandraCollector) Collect() ([]output.Metric, error) {
	//url := "http://localhost:8778"
	url := c.Config.Host + "/jolokia/read/"
	var metrics []output.Metric
	timestamp := time.Now().Unix()

	for _, v := range c.Config.Metrics {
		//first part is path to metric
		metricPath := strings.Split(v, "|")[0]
		metricRename := strings.Split(v, "|")[1]
		//second part is metric name

		res, err := http.Get(url + metricPath)
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		// fmt.Printf("%s\n", string(body))

		js, err := simplejson.NewJson(body)
		if err != nil {
			return nil, err
		}

		value := js.Get("value").MustFloat64()
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, output.Metric{Name: metricRename, Value: value, Timestamp: timestamp, Tags: c.Config.Tags})

	}
	return metrics, nil
}
