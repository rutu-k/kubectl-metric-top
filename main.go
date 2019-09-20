package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "kubectl-metric-top"
	app.Usage = "kube-show-metrics"

	app.Commands = []cli.Command{
		{Name: "fsusage", Usage: "inform pod fsusage", Action: listFsMetrics},
		{Name: "kafka", Usage: "inform pod kafka outgoing bytes", Action: listKafkaMetrics},
		{Name: "rdb", Usage: "inform pod kafka outgoing bytes", Action: listRdbMetrics},
		{Name: "nginx", Usage: "inform pod kafka outgoing bytes", Action: listNginxMetrics},
	}

	app.Run(os.Args)
}
