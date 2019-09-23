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
		{Name: "top", Usage: "inform pod fsusage", Flags: []cli.Flag{cli.BoolFlag{Name: "fsusage"}, cli.BoolFlag{Name: "kafka"}, cli.BoolFlag{Name: "rdb"}, cli.BoolFlag{Name: "nginx"}}, Action: listMetrics},
	}

	app.Run(os.Args)
}
