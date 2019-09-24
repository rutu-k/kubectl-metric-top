package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	"github.com/rutu-k/kubectl-metric-top/pkg/metrics"
)

func main() {
	var metricFlag string

	parser := argparse.NewParser("metrics", "provides custom metrics info")
	fsCmd := parser.NewCommand("fsusage", "provides fs_usage_bytes")
	kafkaCmd := parser.NewCommand("kafka", "provides kafka_server_brokertopicmetrics_total_messagesinpersec_count")
	rdbCmd := parser.NewCommand("rdb", "provides redis_commands_processed")
	nginxCmd := parser.NewCommand("nginx", "provides nginx_connections_accepted")

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if fsCmd.Happened() {
		metricFlag = "fs_usage_bytes"
	} else if kafkaCmd.Happened() {
		metricFlag = "kafka_server_brokertopicmetrics_total_messagesinpersec_count"
	} else if rdbCmd.Happened() {
		metricFlag = "redis_commands_processed"
	} else if nginxCmd.Happened() {
		metricFlag = "nginx_connections_accepted"
	} else {
		fmt.Println("Unknown flag")
		os.Exit(1)
	}
	met := metrics.NewMetrics(metricFlag)
	met.Get()
}
