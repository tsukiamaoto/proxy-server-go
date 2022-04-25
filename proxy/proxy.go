package proxy

import (
	"tsukiamaoto/proxy-server-go/redis"

	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	log "github.com/sirupsen/logrus"
)

type Proxy struct {
	IP        string
	Port      string
	Anonymity string
}

type QueryResponse struct {
	Addr string
	Ok   bool
}

var redisDB *redis.Redis

func init() {
	redisDB = redis.New()
}

func FetchTask() {
	var queryResCh = make(chan QueryResponse, 10)
	var wg sync.WaitGroup

	url := "https://free-proxy-list.net/"
	proxies := ScrapeProxy(url)

	go func() {
		for _, proxy := range proxies {
			wg.Add(1)
			go ValidateProxy(&wg, proxy, queryResCh)
		}
		wg.Wait()
		close(queryResCh)
	}()

	var aliveProxies = make([]string, 0)
	for ch := range queryResCh {
		if ch.Ok {
			aliveProxies = append(aliveProxies, ch.Addr)
		}
	}

	// if proxy is ok, set value to redis for caching
	redisDB.JSONSet("proxy", ".", aliveProxies)
}

func ScrapeProxy(url string) []Proxy {
	var proxies = make([]Proxy, 0)
	c := colly.NewCollector()

	c.OnHTML(".table-striped > tbody > tr", func(e *colly.HTMLElement) {
		ip := e.ChildText("td:first-child")
		port := e.ChildText("td:nth-child(2)")
		anonymity := e.ChildText("td:nth-child(4)")

		proxy := Proxy{
			IP:        ip,
			Port:      port,
			Anonymity: anonymity,
		}

		proxies = append(proxies, proxy)
	})

	c.Visit(url)
	c.Wait()

	return proxies
}

func ValidateProxy(wg *sync.WaitGroup, proxy Proxy, queryResCh chan QueryResponse) {
	defer wg.Done()
	proxyStr := "http://" + proxy.IP + ":" + proxy.Port
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Println(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxyURL),
		},
	}

	url := "https://api.ipify.org?format=json"
	res, err := client.Get(url)
	if err != nil {
		queryResCh <- QueryResponse{
			Addr: proxy.IP,
			Ok:   false,
		}
		return
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	r, _ := regexp.Compile(proxy.IP)
	if r.MatchString(string(data)) {
		queryResCh <- QueryResponse{
			Addr: proxyStr,
			Ok:   true,
		}
	}
}
