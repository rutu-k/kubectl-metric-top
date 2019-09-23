package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type MetricValueList struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
	Metadata   struct {
		SelfLink string `json:"selfLink"`
	}
	Items []struct {
		DescribedObject struct {
			Kind       string `json:"kind"`
			Name       string `json:"name"`
			Namespace  string `json:"namespace"`
			APIVersion string `json:"apiVersion"`
		} `json:"describedObject"`
		MetricName string    `json:"metricName"`
		Timestamp  time.Time `json:"timestamp"`
		Value      string    `json:"value"`
	} `json:"items"`
}

func listMetrics(c *cli.Context) {
	var config *rest.Config

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		fatal(fmt.Sprintf("error in getting Kubeconfig: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fatal(fmt.Sprintf("error in getting clientset from Kubeconfig: %v", err))
	}
	var metricFlag string
	if c.Bool("fsusage") {
		metricFlag = "fs_usage_bytes"
	} else if c.Bool("kafka") {
		metricFlag = "kafka_server_brokertopicmetrics_total_messagesinpersec_count"
	} else if c.Bool("rdb") {
		metricFlag = "redis_commands_processed"
	} else if c.Bool("nginx") {
		metricFlag = "nginx_connections_accepted"
	} else {
		fmt.Println("Unknown Flag")
		os.Exit(1)
	}

	var pods MetricValueList
	path := "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/" + metricFlag
	data, err := clientset.RESTClient().Get().AbsPath(path).DoRaw()
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &pods)
	if err != nil {
		fmt.Println(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, l := range pods.Items {
		data1 := []string{l.DescribedObject.Name, l.DescribedObject.Namespace, l.MetricName, l.Value}
		table.Append(data1)
	}
	table.Render()
}

func fatal(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}
