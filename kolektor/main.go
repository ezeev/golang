package main

import (
	"flag"
	"fmt"

	"github.com/ezeev/golang/kolektor/collector"
	"github.com/ezeev/golang/kolektor/lib"
	"github.com/ezeev/golang/kolektor/output"
)

func main() {

	//Load configs
	flagConfigPath := flag.String("config", "config.json", "Specifies path to configuration file")
	flagCollectorsPath := flag.String("collectors", "/collectors", "Specifies path to collector YAML files")

	flag.Parse()

	config := kolektor.LoadConfig(*flagConfigPath)

	fmt.Println("Starting backend:", config.Backend)

	//Start backend
	be, err := output.NewBackend(config.Backend, config.BackendArgs)
	if err != nil {
		fmt.Println("Error loading backend: ", err)
	}
	//Create listener
	sc := kolektor.NewStatListener()
	fmt.Println("Starting stat listener on port:", config.ListenPort)
	//Start stats listener
	go sc.ListenForStats(config.ListenPort)
	//Start flusher
	go sc.FlushStats(config.Interval, be)
	//start collectors
	fmt.Println("Starting collectors")
	go collector.RunCollectors(*flagCollectorsPath, config)
	fmt.Println("To exit press ctrl+c")
	//TODO check process health here
	for {
		//to keep it running..
	}

}
