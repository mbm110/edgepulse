package task

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/mbm110/edgepulse/utils"

	"github.com/VividCortex/ewma"
)

const (
	bufferSize                     = 1024
	defaultURL                     = "https://speed.cloudflare.com/__down?bytes=200000000"
	defaultTimeout                 = 10 * time.Second
	defaultDisableDownload         = false
	defaultTestNum                 = 10
	defaultMinSpeed        float64 = 0.0
)

var (
	URL     = defaultURL
	Timeout = defaultTimeout
	Disable = defaultDisableDownload

	TestCount = defaultTestNum
	MinSpeed  = defaultMinSpeed
)

func checkDownloadDefault() {
	if URL == "" {
		URL = defaultURL
	}
	if Timeout <= 0 {
		Timeout = defaultTimeout
	}
	if TestCount <= 0 {
		TestCount = defaultTestNum
	}
	if MinSpeed <= 0.0 {
		MinSpeed = defaultMinSpeed
	}
}

func TestDownloadSpeed(ipSet utils.PingDelaySet) (speedSet utils.DownloadSpeedSet) {
	checkDownloadDefault()
	if Disable {
		return utils.DownloadSpeedSet(ipSet)
	}
	if len(ipSet) <= 0 {
		utils.Yellow.Println("[Info] No endpoints passed the latency test; skipping download tests.")
		return
	}
	testNum := TestCount
	if len(ipSet) < TestCount || MinSpeed > 0 {
		testNum = len(ipSet)
	}
	if testNum < TestCount {
		TestCount = testNum
	}

	utils.Cyan.Printf("Starting download test (minimum: %.2f MB/s, count: %d, queue: %d)\n", MinSpeed, TestCount, testNum)
	bar_a := len(strconv.Itoa(len(ipSet)))
	bar_b := "     "
	for i := 0; i < bar_a; i++ {
		bar_b += " "
	}
	bar := utils.NewBar(TestCount, bar_b, "")
	for i := 0; i < testNum; i++ {
		speed, colo := downloadHandler(ipSet[i].IP)
		ipSet[i].DownloadSpeed = speed
		if ipSet[i].Colo == "" {
			ipSet[i].Colo = colo
		}
		if speed >= MinSpeed*1024*1024 {
			bar.Grow(1, "")
			speedSet = append(speedSet, ipSet[i])
			if len(speedSet) == TestCount {
				break
			}
		}
	}
	bar.Done()
	if MinSpeed == 0.00 {
		speedSet = utils.DownloadSpeedSet(ipSet)
	} else if utils.Debug && len(speedSet) == 0 {
		utils.Yellow.Println("[Debug] No endpoint met the minimum speed; returning all measurements for troubleshooting.")
		speedSet = utils.DownloadSpeedSet(ipSet)
	}
	sort.Sort(speedSet)
	return
}

func getDialContext(ip *net.IPAddr) func(ctx context.Context, network, address string) (net.Conn, error) {
	var fakeSourceAddr string
	if isIPv4(ip.String()) {
		fakeSourceAddr = fmt.Sprintf("%s:%d", ip.String(), TCPPort)
	} else {
		fakeSourceAddr = fmt.Sprintf("[%s]:%d", ip.String(), TCPPort)
	}
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

func printDownloadDebugInfo(ip *net.IPAddr, err error, statusCode int, url, lastRedirectURL string, response *http.Response) {
	finalURL := url
	if lastRedirectURL != "" {
		finalURL = lastRedirectURL
	} else if response != nil && response.Request != nil && response.Request.URL != nil {
		finalURL = response.Request.URL.String()
	}
	if url != finalURL {
		if statusCode > 0 {
			utils.Red.Printf("Message", ip.String(), statusCode, url, finalURL)
		} else {
			utils.Red.Printf("Message", ip.String(), err, url, finalURL)
		}
	} else {
		if statusCode > 0 {
			utils.Red.Printf("[Debug] IP: %s, download test stopped，HTTP status: %d, download URL: %s\n", ip.String(), statusCode, url)
		} else {
			utils.Red.Printf("[Debug] IP: %s, download test failed，error: %v, download URL: %s\n", ip.String(), err, url)
		}
	}
}

// return download Speed
func downloadHandler(ip *net.IPAddr) (float64, string) {
	var lastRedirectURL string
	client := &http.Client{
		Transport: &http.Transport{DialContext: getDialContext(ip)},
		Timeout:   Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			lastRedirectURL = req.URL.String()
			if len(via) > 10 {
				if utils.Debug {
					utils.Red.Printf("Message", ip.String(), req.URL.String())
				}
				return http.ErrUseLastResponse
			}
			if req.Header.Get("Referer") == defaultURL {
				req.Header.Del("Referer")
			}
			return nil
		},
	}
	defer client.CloseIdleConnections()
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		if utils.Debug {
			utils.Red.Printf("[Debug] IP: %s, Download testrequest creation failed，error: %v, download URL: %s\n", ip.String(), err, URL)
		}
		return 0.0, ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.80 Safari/537.36")

	response, err := client.Do(req)
	if err != nil {
		if utils.Debug {
			printDownloadDebugInfo(ip, err, 0, URL, lastRedirectURL, response)
		}
		return 0.0, ""
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		if utils.Debug {
			printDownloadDebugInfo(ip, nil, response.StatusCode, URL, lastRedirectURL, response)
		}
		return 0.0, ""
	}

	colo := getHeaderColo(response.Header)

	timeStart := time.Now()
	timeEnd := timeStart.Add(Timeout)

	contentLength := response.ContentLength
	buffer := make([]byte, bufferSize)

	var (
		contentRead     int64 = 0
		timeSlice             = Timeout / 100
		timeCounter           = 1
		lastContentRead int64 = 0
	)

	var nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
	e := ewma.NewMovingAverage()

	for contentLength != contentRead {
		currentTime := time.Now()
		if currentTime.After(nextTime) {
			timeCounter++
			nextTime = timeStart.Add(timeSlice * time.Duration(timeCounter))
			e.Add(float64(contentRead - lastContentRead))
			lastContentRead = contentRead
		}
		if currentTime.After(timeEnd) {
			break
		}
		bufferRead, err := response.Body.Read(buffer)
		if err != nil {
			if err != io.EOF {
				break
			} else if contentLength == -1 {
				break
			}
			last_time_slice := timeStart.Add(timeSlice * time.Duration(timeCounter-1))
			e.Add(float64(contentRead-lastContentRead) / (float64(currentTime.Sub(last_time_slice)) / float64(timeSlice)))
		}
		contentRead += int64(bufferRead)
	}
	return e.Value() / (Timeout.Seconds() / 120), colo
}
