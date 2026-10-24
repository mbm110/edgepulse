package task

import (
	//"crypto/tls"

	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mbm110/edgepulse/utils"
)

var (
	Httping               bool
	HttpingStatusCode     int
	HttpingCFColo         string
	HttpingCFColomap      *sync.Map
	RegexpColoIATACode    = regexp.MustCompile(`[A-Z]{3}`)
	RegexpColoCountryCode = regexp.MustCompile(`[A-Z]{2}`)
	RegexpColoGcore       = regexp.MustCompile(`^[a-z]{2}`)
)

// pingReceived pingTotalTime
func (p *Ping) httping(ip *net.IPAddr) (int, time.Duration, string) {
	hc := http.Client{
		Timeout: time.Second * 2,
		Transport: &http.Transport{
			DialContext: getDialContext(ip),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	defer hc.CloseIdleConnections()

	var colo string
	{
		request, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			if utils.Debug {
				utils.Red.Printf("[Debug] IP: %s, Latency testrequest creation failed，error: %v, Test URL: %s\n", ip.String(), err, URL)
			}
			return 0, 0, ""
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		response, err := hc.Do(request)
		if err != nil {
			if utils.Debug {
				utils.Red.Printf("[Debug] IP: %s, latency test failed，error: %v, Test URL: %s\n", ip.String(), err, URL)
			}
			return 0, 0, ""
		}
		defer response.Body.Close()

		//fmt.Println("IP:", ip, "StatusCode:", response.StatusCode, response.Request.URL)
		if HttpingStatusCode == 0 ||  HttpingStatusCode < 100 || HttpingStatusCode > 599 {
			if response.StatusCode != 200 && response.StatusCode != 301 && response.StatusCode != 302 {
				if utils.Debug {
					utils.Red.Printf("[Debug] IP: %s, latency test stopped，HTTP status: %d, Test URL: %s\n", ip.String(), response.StatusCode, URL)
				}
				return 0, 0, ""
			}
		} else {
			if response.StatusCode != HttpingStatusCode {
				if utils.Debug {
					utils.Red.Printf("Message", ip.String(), response.StatusCode, HttpingStatusCode, URL)
				}
				return 0, 0, ""
			}
		}

		io.Copy(io.Discard, response.Body)

		colo = getHeaderColo(response.Header)

		if HttpingCFColo != "" {
			colo = p.filterColo(colo)
			if colo == "" {
				if utils.Debug {
					utils.Red.Printf("[Debug] IP: %s, Regiondid not match: %s\n", ip.String(), colo)
				}
				return 0, 0, ""
			}
		}
	}

	success := 0
	var delay time.Duration
	for i := 0; i < PingTimes; i++ {
		request, err := http.NewRequest(http.MethodHead, URL, nil)
		if err != nil {
			log.Fatal("Message", err)
			return 0, 0, ""
		}
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")
		if i == PingTimes-1 {
			request.Header.Set("Connection", "close")
		}
		startTime := time.Now()
		response, err := hc.Do(request)
		if err != nil {
			continue
		}
		success++
		io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
		duration := time.Since(startTime)
		delay += duration
	}

	return success, delay, colo
}

func MapColoMap() *sync.Map {
	if HttpingCFColo == "" {
		return nil
	}
	colos := strings.Split(strings.ToUpper(HttpingCFColo), ",")
	colomap := &sync.Map{}
	for _, colo := range colos {
		colomap.Store(colo, colo)
	}
	return colomap
}

func getHeaderColo(header http.Header) (colo string) {
	if header.Get("server") != "" {
		// server: cloudflare
		// cf-ray: 7bd32409eda7b020-SJC
		if header.Get("server") == "cloudflare" {
			if colo = header.Get("cf-ray"); colo != "" {
				return RegexpColoIATACode.FindString(colo)
			}
		}
		// server: CDN77-Turbo
		// x-77-pop: frankfurtDE
		// x-77-pop: amsterdamNL
		// x-77-pop: singaporeSG
		if header.Get("server") == "CDN77-Turbo" {
			if colo = header.Get("x-77-pop"); colo != "" {
				return RegexpColoCountryCode.FindString(colo)
			}
		}
		// server: BunnyCDN-TW1-1121
		if colo = header.Get("server"); strings.Contains(colo, "BunnyCDN-") {
			return RegexpColoCountryCode.FindString(strings.TrimPrefix(colo, "BunnyCDN-"))
		}
	}
	// x-amz-cf-pop: SIN52-P1
	if colo = header.Get("x-amz-cf-pop"); colo != "" {
		return RegexpColoIATACode.FindString(colo)
	}
	// x-served-by: cache-qpg1275-QPG
	if colo = header.Get("x-served-by"); colo != "" {
		if matches := RegexpColoIATACode.FindAllString(colo, -1); len(matches) > 0 {
			return matches[len(matches)-1]
		}
	}
	// x-id-fe: fr5-hw-edge-gc17
	// x-shard: fr5-shard0-default
	// x-id: fr5-hw-edge-gc28
	if colo = header.Get("x-id-fe"); colo != "" {
		if colo = RegexpColoGcore.FindString(colo); colo != "" {
			return strings.ToUpper(colo)
		}
	}

	return ""
}

func (p *Ping) filterColo(colo string) string {
	if colo == "" {
		return ""
	}
	if HttpingCFColomap == nil {
		return colo
	}
	_, ok := HttpingCFColomap.Load(colo)
	if ok {
		return colo
	}
	return ""
}
