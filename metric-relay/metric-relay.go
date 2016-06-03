package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront"
)

func sendMetrics(relays string, ch <-chan string) {
	for msg := range ch {
		fmt.Printf("%s\n", msg)
		if doMetrics == true {
			numLinesSent.Inc(1)
		}
	}
	if doMetrics == true {
		numGoRoutinesGauge.Update(int64(runtime.NumGoroutine()))
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn, ch chan string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		if strings.HasPrefix(msg, "version") {
			conn.Write([]byte("ack"))
		} else {
			if doMetrics == true {
				numLinesReceived.Inc(1)
			}
			ch <- msg
		}
	}
	conn.Close()
}

func runServer(host string, port int, ch chan string) {
	l, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Listening on " + host + ":" + strconv.Itoa(port))
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, ch)
	}
}

var doMetrics bool
var numGoRoutinesGauge metrics.Gauge
var numLinesReceived metrics.Counter
var numLinesSent metrics.Counter

func initMetrics(proxyAddr string) {
	hostTags := map[string]string{
		"source": "ubuntudev",
	}

	//func initRelays(relays string) {
	//arrRelays := strings.Split(relays,",")
	//}

	//num go routines
	numGoRoutinesGauge = metrics.NewGauge()
	metrics.Register("num-go-routines", numGoRoutinesGauge)

	numLinesReceived = metrics.NewCounter()
	metrics.Register("lines-received", numLinesReceived)

	numLinesSent = metrics.NewCounter()
	metrics.Register("lines-sent", numLinesSent)

	addr, err := net.ResolveTCPAddr("tcp", proxyAddr)
	if err != nil {
		doMetrics = false
		log.Fatalf("Unable to connect to Wavefront proxy: %s", err.Error())
	} else {
		doMetrics = true
		go wavefront.Wavefront(metrics.DefaultRegistry, 1*time.Second, hostTags, "metric-relay", addr)
	}
}

func main() {

	flags := flag.NewFlagSet("stream-splitter-args", flag.ExitOnError)
	host := flags.String("host", "", "Hostname/IP address to listen on")
	proxyAddr := flags.String("proxy-addr", "", "Address and port of Wavefront proxy (optional)")
	port := flags.Int("port", 9991, "Port to listen on")
	relays := flags.String("relays", "", "Comma separated list of addresses to relay to")
	bufferSize := flags.Int("buffer-size", 100, "Buffer size for the main channel")
	flags.Parse(os.Args[1:])

	if *proxyAddr != "" {
		initMetrics(*proxyAddr)
	}

	//initialize channel
	ch := make(chan string, *bufferSize)
	//start the server, it will pass the channel to handleRequest,
	//which will read data into the channel
	go runServer(*host, *port, ch)

	//will send any data read into the channel
	go sendMetrics(*relays, ch)

	/*go func() {
		log.Println(http.ListenAndServe("127.0.0.1:9992", nil))
	}()*/

	runtime.Goexit()

}
