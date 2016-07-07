
package pairs
import (
        "bufio"
        "bytes"
        "fmt"
        "net"
        "net/http"
        "net/url"
        "strconv"
        "strings"
        "sync"
        "time"
        "github.com/influxdata/telegraf"
        "github.com/influxdata/telegraf/plugins/inputs"
)
type Pairs struct {
        Urls, Metrics []string
        AllMetrics    int
}
var sampleConfig = `
  ### An array of Urls containing metrics in key/value pair format.
  urls = ["http://localhost/stats"]
  Best practice is to create config file with url array and name_override for measurement name.
  Example:
`
func (n *Pairs) SampleConfig() string {
        return sampleConfig
}
func (n *Pairs) Description() string {
        return "Read web pages with metrics in key/value pair format"
}
func (n *Pairs) Gather(acc telegraf.Accumulator) error {
        var wg sync.WaitGroup
        var outerr error
        for _, u := range n.Urls {
                addr, err := url.Parse(u)
                if err != nil {
                        return fmt.Errorf("Unable to parse address '%s': %s", u, err)
                }
                wg.Add(1)
                go func(addr *url.URL) {
                        defer wg.Done()
                        outerr = n.gatherUrl(addr, acc)
                }(addr)
        }
        wg.Wait()
        return outerr
}
var tr = &http.Transport{
        ResponseHeaderTimeout: time.Duration(3 * time.Second),
}
var client = &http.Client{Transport: tr}
// Check to see if Metric names listed in conf are on webpage.  If AllMetrics = 1 bypass.
func (n *Pairs) stringInSlice(str string) bool {
        for _, v := range n.Metrics {
                if strings.Replace(v, " ", "_", -1) == str {
                        return true
                }
        }
        return false
}
// Loop through the URL Array
func (n *Pairs) gatherUrl(addr *url.URL, acc telegraf.Accumulator) error {
        resp, err := client.Get(addr.String())
        if err != nil {
                return fmt.Errorf("error making HTTP request to %s: %s", addr.String(), err)
        }
        defer resp.Body.Close()
        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
        }
        tags := getTags(addr)
        sc := bufio.NewScanner(resp.Body)
        fields := make(map[string]interface{})
        // Evaluate each line of the individual URLs
        for sc.Scan() {
                line := sc.Text()
                lineLen := bytes.TrimSpace(sc.Bytes())
                if len(lineLen) == 0 {
                        continue
                }
                //Look for Key Value pairs with an = as a seperator
                if strings.Contains(line, "=") {
                        // Separate the key from the value
                        parts := strings.SplitN(line, "=", 2)
                        key, part := strings.Replace(parts[0], " ", "_", -1), strings.TrimSpace(parts[1])
                        value, err := strconv.ParseFloat(part, 64)
                        if err != nil {
                                continue
                        }
                        // Evaluate wether to grab all metrics on the page, or look in the conf file for specific(default)
                        if n.AllMetrics == 0 {
                                if n.stringInSlice(key) {
                                        fields[key] = value
                                }
                        } else {
                                fields[key] = value
                        }
                        // Look for Key Value pairs with an : as a seperator
                } else if strings.Contains(line, ":") {
                        parts := strings.SplitN(line, ":", 2)
                        key, part := strings.Replace(parts[0], " ", "_", -1), strings.TrimSpace(parts[1])
                        value, err := strconv.ParseFloat(part, 64)
                        if err != nil {
                                continue
                        }
                        // Evaluate wether to grab all metrics on the page, or look in the conf file for specific(default)
                        if n.AllMetrics == 0 {
                                if n.stringInSlice(key) {
                                        fields[key] = value
                                }
                        } else {
                                fields[key] = value
                        }
                }
        }
        acc.AddFields("pairs", fields, tags)
        return nil
}

func getTags(addr *url.URL) map[string]string {
        h := addr.Host
        host, port, err := net.SplitHostPort(h)
        if err != nil {
                host = addr.Host
                if addr.Scheme == "http" {
                        port = "80"
                } else if addr.Scheme == "https" {
                        port = "443"
                } else {
                        port = ""
                }
        }
        return map[string]string{"host": host, "port": port}
}
func init() {
        inputs.Add("pairs", func() telegraf.Input {
                return &Pairs{}
        })
}
