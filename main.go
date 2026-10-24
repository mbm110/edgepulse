package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/mbm110/edgepulse/task"
	"github.com/mbm110/edgepulse/utils"
)

var (
	version, versionNew string
)

func init() {
	var printVersion bool
	var help = `
EdgePulse ` + version + `
Measure latency and throughput across CDN or website IP ranges to find the strongest IPv4 and IPv6 endpoints.
https://github.com/mbm110/edgepulse

Options:
    -n 200
        Latency workers; more workers finish sooner, but use fewer on low-power devices (default: 200, maximum: 1000).
    -t 4
        Latency probes per endpoint (default: 4).
    -dn 10
        Number of endpoints checked for download speed after latency sorting (default: 10).
    -dt 10
        Maximum download test time per endpoint (default: 10 seconds).
    -tp 443
        Port used for latency and download tests (default: 443).
    -url https://speed.cloudflare.com/__down?bytes=200000000
        URL used for HTTP and download tests; hosting your own endpoint is recommended.

    -httping
        Use HTTP latency probes with the URL supplied by -url (default: TCP).
    -httping-code 200
        Accepted HTTP status code for HTTP probes (default: 200, 301, 302).
    -cfcolo HKG,KHH,NRT,LAX,SEA,SJC,FRA,MAD
        Filter by region code; available in HTTP mode (default: all regions).

    -tl 200
        Maximum average latency; only endpoints below this value are kept (default: 9999 ms).
    -tll 40
        Minimum average latency (default: 0 ms).
    -tlr 0.2
        Maximum loss ratio from 0.00 to 1.00 (default: 1.00).
    -sl 5
        Minimum download speed; testing stops after enough matching endpoints are found (default: 0.00 MB/s).

    -p 10
        Number of displayed results; zero hides the table (default: 10).
    -f ip.txt
        IP range input file; other CDN ranges are supported (default: ip.txt).
    -ip 1.1.1.1,2.2.2.2/24,2606:4700::/32
        Comma-separated IP ranges supplied directly (default: empty).
    -o result.csv
        CSV output path; use an empty value to disable output (default: result.csv).

    -dd
        Disable download tests and sort by latency (default: enabled).
    -allip
        Test every address in each range; otherwise one address is selected per range.

    -debug
        Verbose diagnostic output for unexpected conditions (default: off).

    -v
        Print the program version and check for updates
    -h
        Print help
`
	var minDelay, maxDelay, downloadTime int
	var maxLossRate float64
	flag.IntVar(&task.Routines, "n", 200, "Latency test workers")
	flag.IntVar(&task.PingTimes, "t", 4, "Latency probes per endpoint")
	flag.IntVar(&task.TestCount, "dn", 10, "Download test count")
	flag.IntVar(&downloadTime, "dt", 10, "Download test duration")
	flag.IntVar(&task.TCPPort, "tp", 443, "Test port")
	flag.StringVar(&task.URL, "url", "https://speed.cloudflare.com/__down?bytes=200000000", "Test URL")

	flag.BoolVar(&task.Httping, "httping", false, "Switch test mode")
	flag.IntVar(&task.HttpingStatusCode, "httping-code", 0, "Accepted HTTP status")
	flag.StringVar(&task.HttpingCFColo, "cfcolo", "", "Filter by region")

	flag.IntVar(&maxDelay, "tl", 9999, "Maximum average latency")
	flag.IntVar(&minDelay, "tll", 0, "Minimum average latency")
	flag.Float64Var(&maxLossRate, "tlr", 1, "Maximum loss ratio")
	flag.Float64Var(&task.MinSpeed, "sl", 0, "Minimum download speed")

	flag.IntVar(&utils.PrintNum, "p", 10, "Displayed result count")
	flag.StringVar(&task.IPFile, "f", "ip.txt", "IP range file")
	flag.StringVar(&task.IPText, "ip", "", "Inline IP ranges")
	flag.StringVar(&utils.Output, "o", "result.csv", "Output file")

	flag.BoolVar(&task.Disable, "dd", false, "Disable download tests")
	flag.BoolVar(&task.TestAll, "allip", false, "Test every IP")

	flag.BoolVar(&utils.Debug, "debug", false, "Debug output")

	flag.BoolVar(&printVersion, "v", false, "Print version")
	flag.Usage = func() { fmt.Print(help) }
	flag.Parse()

	if task.MinSpeed > 0 && time.Duration(maxDelay)*time.Millisecond == utils.InputMaxDelay {
		utils.Yellow.Println("[Tip] Pair -sl with -tl to avoid testing indefinitely when too few endpoints meet the speed target.")
	}
	utils.InputMaxDelay = time.Duration(maxDelay) * time.Millisecond
	utils.InputMinDelay = time.Duration(minDelay) * time.Millisecond
	utils.InputMaxLossRate = float32(maxLossRate)
	task.Timeout = time.Duration(downloadTime) * time.Second
	task.HttpingCFColomap = task.MapColoMap()

	if printVersion {
		println(version)
		fmt.Println("Checking for updates...")
		checkUpdate()
		if versionNew != "" {
			utils.Yellow.Printf("*** A newer version [%s] is available. Visit the EdgePulse release page to update. ***", versionNew)
		} else {
			utils.Green.Println("Already running the latest version [" + version + "]!")
		}
		os.Exit(0)
	}
}

func main() {
	task.InitRandSeed() // Seed random generator

	fmt.Printf("# mbm110/edgepulse %s \n\n", version)

	pingData := task.NewPing().Run().FilterDelay().FilterLossRate()
	speedData := task.TestDownloadSpeed(pingData)
	utils.ExportCsv(speedData) // Write output file
	speedData.Print()          // Print results
	endPrint()                 // Select Windows exit behavior
}

func endPrint() {
	if utils.NoPrintResult() { // exit when result printing is disabled
		return
	}
	if runtime.GOOS == "windows" { // Windows waits for Enter or Ctrl+C after a double-click run
		fmt.Printf("Press Enter or Ctrl+C to exit.")
		fmt.Scanln()
	}
}

func checkUpdate() {
	timeout := 10 * time.Second
	client := http.Client{Timeout: timeout}
	res, err := client.Get("https://api.github.com/repos/mbm110/edgepulse/releases/latest")
	if err != nil {
		return
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var release struct {
		TagName string `json:"tag_name"`
	}
	if json.Unmarshal(body, &release) == nil && release.TagName != "" && release.TagName != version {
		versionNew = release.TagName
	}
}
