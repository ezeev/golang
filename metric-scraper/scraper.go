package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strings"
	"time"
)

type Scraper struct {
	ScrapeUrl      string
	Interval       int
	SourceTagFile  string
	SourceTagParts []string
}

func NewScraper(scrapeUrl string, interval int, sourceTagFile string) Scraper {
	scraper := Scraper{
		ScrapeUrl:     scrapeUrl,
		Interval:      interval,
		SourceTagFile: sourceTagFile,
	}
	log.Printf("Initialized scraper with:\nURL: %s\nSource Tag List: %s\nInterval: %d\n", scraper.ScrapeUrl, scraper.SourceTagFile, scraper.Interval)

	scraper.LoadSourceTagParts()

	return scraper
}

func (this *Scraper) Run() {
	tick := time.Tick(1 * time.Second)
	tickCount := 0
	for {
		select {
		case <-tick:
			tickCount++
			if tickCount == this.Interval {
				log.Print("Starting Crawl")
				this.CrawlPage()
				tickCount = 0
			}
		}
	}
}

func (this *Scraper) ParseMetric(metric string) (string, string, string) {

	//extract the value first
	metricNameAndValue := strings.Split(metric, "=")
	metricName := metricNameAndValue[0]
	metricValue := metricNameAndValue[1]

	metricParts := strings.Split(metricName, ".")
	//iterate through possible source tags
	for v := range this.SourceTagParts {
		for part := range metricParts {
			if strings.Contains(metricParts[part], this.SourceTagParts[v]) {
				//there's a match
				sourceTag := metricParts[part]
				//delete this part from metric name
				metricParts = append(metricParts[:part], metricParts[part+1:]...)
				newMetricName := strings.Join(metricParts, ".")
				return newMetricName, metricValue, sourceTag
			}
		}
	}
	return metric, metricValue, "unknown"
}

func (this *Scraper) CrawlPage() error {
	path := this.ScrapeUrl
	file, err := os.Open(path)
	//current timestamp
	timestamp := time.Now().Unix()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	var metrics []string
	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		m1 := scanner.Text()
		//see if it contains any source tags
		m2, value, sourceTag := this.ParseMetric(m1)
		log.Printf("%s %s %d source=%s", m2, value, timestamp, sourceTag)

		metrics = append(metrics, scanner.Text())
		count = count + 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Scraped %d metrics", count)

	return nil
}

func (this *Scraper) LoadSourceTagParts() error {
	path := this.SourceTagFile
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	var parts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts = append(parts, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	this.SourceTagParts = parts
	return nil
}

func main() {
	var scrape_url = flag.String("url", "", "The URL you wish to scrape")
	var interval = flag.Int("interval", 10, "The number of seconds between scrapes")
	var sourcetag_prefix_file = flag.String("sourcetag_prefix_file", "", "A file containing a list of strings of potential source tags")
	flag.Parse()
	scraper := NewScraper(*scrape_url, *interval, *sourcetag_prefix_file)
	scraper.Run()
}
