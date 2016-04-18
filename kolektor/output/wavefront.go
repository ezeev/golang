
package output

import (
  "fmt"
  "net"
  "strconv"
)

type WavefrontBackend struct {
  DefaultSource string
  ProxyAddress string
  ProxyPort string
  Prefix string
  PrintStats bool
}

func NewWavefrontBackend(beArgs map[string]string) (*WavefrontBackend,error) {

  wf := &WavefrontBackend{
    DefaultSource: beArgs["defaultSource"],
    ProxyAddress: beArgs["proxyAddress"],
    ProxyPort: beArgs["proxyPort"],
    Prefix: beArgs["prefix"],
  }

  printStats, err := strconv.ParseBool(beArgs["printStats"])
  if (err != nil) {
    fmt.Println("Error parsing printStats, should be \"true\" or \"false\"",err)
  }
  wf.PrintStats = printStats
	return wf, nil
}

func (b *WavefrontBackend) tagString(tags map[string]string) string {
  tagStr := ""
  _, ok := tags["source"];
  if (!ok) {
    tagStr = "source="+b.DefaultSource
  }
  for k,v := range tags {
    tagStr += k+"="+v+" "
  }
  return tagStr
}

func (b *WavefrontBackend) Flush(metrics []Metric) {
  //open connection
  conn, err := net.Dial("tcp", b.ProxyAddress+":"+b.ProxyPort)
  if (err != nil) {
    fmt.Println("Error connecting:",err)
    conn.Close()
  }

  for _,v := range metrics {
    metricName := b.Prefix+v.Name;
    value := v.Value
    ts := v.Timestamp
    tagStr := b.tagString(v.Tags)
    //name value timestamp tags
    if (b.PrintStats) {
      fmt.Printf("%s %v %v %s \n", metricName,value,ts,tagStr)
    }
    fmt.Fprintf(conn, fmt.Sprintf("%s %v %v %s \n", metricName,value,ts,tagStr))
  }

  //fmt.Println(metrics)
  /*
  if (b.PrintStats) {
    fmt.Println("Printing Metric Data:")
    fmt.Println(metrics)
  }
  */
  //close connection
  conn.Close()

}
