package main

import (
	"encoding/json"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	//"k8s.io/client-go/tools/clientcmd"
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

func getMetrics(clientset *kubernetes.Clientset, pods *MetricValueList) error {
	data, err := clientset.RESTClient().Get().AbsPath("/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/fs_usage_bytes").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func getKafkaMetrics(clientset *kubernetes.Clientset, pods *MetricValueList) error {
	data, err := clientset.RESTClient().Get().AbsPath("/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/kafka_server_brokertopicmetrics_total_messagesinpersec_count").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func getRdbMetrics(clientset *kubernetes.Clientset, pods *MetricValueList) error {
	data, err := clientset.RESTClient().Get().AbsPath("/apis/custom.metrics.k8s.io/v1beta1/namespaces/redis/pods/*/redis_commands_processed").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func getNginxMetrics(clientset *kubernetes.Clientset, pods *MetricValueList) error {
	data, err := clientset.RESTClient().Get().AbsPath("/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/nginx_connections_accepted").DoRaw()
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &pods)
	return err
}

func listFsMetrics(c *cli.Context) {
	//fmt.Println("listing running Pods")
	clientset := getKubeHandle()

	var pods MetricValueList
	err := getMetrics(clientset, &pods)
	if err != nil {
		panic(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, m := range pods.Items {
		data := []string{m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value}
		table.Append(data)
		//fmt.Println(m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value)
	}
	table.Render()
}

func listKafkaMetrics(c *cli.Context) {
	//fmt.Println("listing running Pods")
	clientset := getKubeHandle()

	var pods MetricValueList
	err := getKafkaMetrics(clientset, &pods)
	if err != nil {
		panic(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, m := range pods.Items {
		data := []string{m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value}
		table.Append(data)
		//fmt.Println(m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value)
	}
	table.Render()
}

func listRdbMetrics(c *cli.Context) {
	//fmt.Println("listing running Pods")
	clientset := getKubeHandle()

	var pods MetricValueList
	err := getRdbMetrics(clientset, &pods)
	if err != nil {
		panic(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, m := range pods.Items {
		data := []string{m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value}
		table.Append(data)
		//fmt.Println(m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value)
	}
	table.Render()
}

func listNginxMetrics(c *cli.Context) {
	//fmt.Println("listing running Pods")
	clientset := getKubeHandle()

	var pods MetricValueList
	err := getNginxMetrics(clientset, &pods)
	if err != nil {
		panic(err.Error())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, m := range pods.Items {
		data := []string{m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value}
		table.Append(data)
		//fmt.Println(m.DescribedObject.Name, m.DescribedObject.Namespace, m.MetricName, m.Value)
	}
	table.Render()
}
