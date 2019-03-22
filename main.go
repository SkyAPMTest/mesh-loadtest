package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"math"
	"math/rand"
	"os"
	"runtime"
)

var (
	insecure = flag.Bool("insecure", true, "Specify for non TLS connection")

	h = flag.String("h", "oap:11800", "Host of target.")

	sn = flag.Uint("sn", 10, "Number of services.")
	in = flag.Uint("in", 10, "Number of service instances.")
	en = flag.Uint("en", 100, "Number of service endpoints.")

	c = flag.Uint("c", 10, "Number of requests to run concurrently.")
	n = flag.Uint("n", 100, "Number of requests to run. Default is 200.")
	q = flag.Uint("q", 0, "Rate limit, in queries per second (QPS). Default is no rate limit.")
	z = flag.Duration("z", 0, "Duration of application to send requests.")

	ct = flag.Uint("T", 10, "Connection timeout in seconds for the initial connection dial.")
	kt = flag.Uint("L", 0, "Keepalive time in seconds.")

	cpus = flag.Uint("cpus", uint(runtime.GOMAXPROCS(-1)), "")

	)

func main() {

	flag.Parse()
	total := *sn * *en * 2
	data := make([]map[string]interface{}, total)
	num := 0
	oldInsId := 0
	insId := 0
	for s := 0; s < int(*sn); s++ {
		for e := 0; e < int(*en); e++ {
			data[num] = map[string]interface{}{
				"startTime": "{{.Starttime}}",
				"endTime": "{{.Endtime}}",
				"sourceServiceName": fmt.Sprintf("/service/%d", s - 1),
				"sourceServiceInstance": fmt.Sprintf("/service/%d/%d", s - 1, oldInsId),
				"destServiceName": fmt.Sprintf("/service/%d", s),
				"destServiceInstance": fmt.Sprintf("/service/%d/%d", s, insId),
				"latency": 500,
				"responseCode": 200,
				"status": true,
				"detectPoint": 1,
				"protocol": 1,
				"endpoint": fmt.Sprintf("/endpoint/%d/%d", s, e),
			}
			num++
			oldInsId = insId
			insId = rand.Intn(int(*in))
			data[num] = map[string]interface{}{
				"startTime": "{{.Starttime}}",
				"endTime": "{{.Endtime}}",
				"sourceServiceName": fmt.Sprintf("/service/%d", s),
				"sourceServiceInstance": fmt.Sprintf("/service/%d/%d", s, oldInsId),
				"destServiceName": fmt.Sprintf("/service/%d", s + 1),
				"destServiceInstance": fmt.Sprintf("/service/%d/%d", s + 1, insId),
				"latency": 500,
				"responseCode": 200,
				"status": true,
				"detectPoint": 0,
				"protocol": 1,
				"endpoint": fmt.Sprintf("/endpoint/%d/%d", s, e),
			}
			num++
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	reqNum := *n
	if *z > 0 {
		reqNum = math.MaxInt32
	}
	report, err := runner.Run(
		"ServiceMeshMetricService.collect",
		*h,
		runner.WithProtoset("bundle.protoset"),
		runner.WithDataFromJSON(string(b)),
		runner.WithTotalRequests(reqNum),
		runner.WithConcurrency(*c),
		runner.WithRunDuration(*z),
		runner.WithCPUs(*cpus),
		runner.WithQPS(*q),
		runner.WithInsecure(*insecure),
	)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	printer := printer.ReportPrinter{
		Out:    os.Stdout,
		Report: report,
	}

	printer.Print("pretty")
}
