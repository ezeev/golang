package collector

import (
  "time"
  "github.com/bitly/go-simplejson"
  "github.com/ezeev/golang/kolektor/output"
  "log"
  "net/http"
  "io/ioutil"
  "strings"
)


type CassandraCollector struct {
  //Host string
  //Interval float64
  LastCollection time.Time
  //Metrics []string
  //Tags map[string]string
  Config CassandraConfig
}

type CassandraConfig struct {
  Collector string
  Host string
  Interval float64
  Metrics []string
  Tags map[string]string
}


func NewCassandraCollector(config CassandraConfig) (*CassandraCollector,error) {
  /*c := &CassandraCollector{
    Host: config.Host,
    Interval: config.Interval,
    Metrics: config.Metrics,
    Tags: config.Tags,
  }*/
  c := &CassandraCollector{ Config: config }
	return c, nil
}

func (c *CassandraCollector) GetInterval() float64 {
  return c.Config.Interval
}

func (c *CassandraCollector) SetLastCollectionTime(t time.Time) {
  c.LastCollection = t
}

func (c *CassandraCollector) GetLastCollectionTime() time.Time {
  return c.LastCollection
}

func (c *CassandraCollector) Collect() []output.Metric {
  //url := "http://localhost:8778"
  url := c.Config.Host + "/jolokia/read/"
  metrics := make([]output.Metric,0)
  timestamp := time.Now().Unix()


  for _, v := range c.Config.Metrics {
    //first part is path to metric
    metricPath := strings.Split(v,"|")[0]
    metricRename := strings.Split(v,"|")[1]
    //second part is metric name

    res, err := http.Get(url + metricPath)
    if err != nil {
        log.Fatalln(err)
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatalln(err)
    }

    // fmt.Printf("%s\n", string(body))

    js, err := simplejson.NewJson(body)
    if err != nil {
        log.Fatalln(err)
    }

    value := js.Get("value").MustFloat64()
    if err != nil {
        log.Fatalln(err)
    }
    //metrics[metricRename] = value
    metrics = append(metrics,output.Metric{ Name: metricRename, Value: value, Timestamp: timestamp, Tags: c.Config.Tags, })

  }
  return metrics
}
