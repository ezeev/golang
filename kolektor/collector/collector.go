package collector

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ezeev/golang/kolektor/lib"
	"github.com/ezeev/golang/kolektor/output"
	"gopkg.in/yaml.v2"
)

//CollectorConfig collector config
type CollectorConfig struct {
	Collector string
}

//Collector collector
type Collector interface {
	Name() string
	Collect() ([]output.Metric, error)
	Interval() float64
	SetLastCollectionTime(time.Time)
	LastCollectionTime() time.Time
}

//NewCollector Creates an instance of a collector based on the yaml config.
func NewCollector(strYaml string) (Collector, error) {

	genericConfig := CollectorConfig{}
	err := yaml.Unmarshal([]byte(strYaml), &genericConfig)
	if err != nil {
		return nil, err
	}

	if genericConfig.Collector == "cassandra" {
		cassConfig := CassandraConfig{}
		err2 := yaml.Unmarshal([]byte(strYaml), &cassConfig)
		if err2 != nil {
			return nil, err2
		}
		return NewCassandraCollector(cassConfig)
	}

	fmt.Println(genericConfig.Collector)
	return nil, err
}

// CollectAndFlush Executes the collector and sends the returned metrics to the backend
func CollectAndFlush(collector Collector, be output.Backend) {
	metrics, err := collector.Collect()
	if err != nil {
		fmt.Println("Error from collector:", collector.Name(), "Error:", err)
	} else {
		go be.Flush(metrics)
	}
}

//RunCollectors Runs each collector. Loads the collectors and iterates through each one until the process is cancelled.
func RunCollectors(pathToYamlFiles string, config kolektor.Configuration) {

	//load YAML files
	files, ferr := ioutil.ReadDir(pathToYamlFiles)
	if ferr != nil {
		fmt.Println(ferr)
	}
	//load collectors
	collectors := make([]Collector, len(files)) //hole the collectors in a slice
	for _, f := range files {
		buf, err := ioutil.ReadFile(pathToYamlFiles + "/" + f.Name())
		if err != nil {
			fmt.Println("Unable to read collector Yaml:", err)
			panic(err)
		}
		yamlstr := string(buf)
		collector, err := NewCollector(yamlstr)
		if err != nil {
			panic(err)
		}
		collector.SetLastCollectionTime(time.Now()) //initialize timer
		collectors = append(collectors, collector)  //add the collector to the slice
	}
	//end collector loading

	//create backend
	be, err := output.NewBackend(config.Backend, config.BackendArgs)
	if err != nil {
		fmt.Println("Error loading backend: ", err)
	}

	//time and loop collectors
	for {
		for _, collector := range collectors { //loop through the collectors forever
			if collector != nil {
				t := time.Now()
				d := t.Sub(collector.LastCollectionTime()) //get duration since last flush
				//fmt.Println(d.Seconds())
				if d.Seconds() >= collector.Interval() { //is it time to collect?
					now := time.Now()
					go CollectAndFlush(collector, be)
					collector.SetLastCollectionTime(now)
				}
			}
		}
	}

}
