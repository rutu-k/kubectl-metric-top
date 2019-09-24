package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type Metrics struct {
	Name string
}

func NewMetrics(metricFlag string) Metrics {
	return Metrics{
		Name: metricFlag,
	}
}

func (m Metrics) Get() {
	mflag := m.Name
	str := strings.Split(mflag, "_")
	str1 := str[0]

	var config *rest.Config
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		fatal(fmt.Sprintf("error in getting Kubeconfig: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fatal(fmt.Sprintf("error in getting clientset from Kubeconfig: %v", err))
	}

	podz, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		fatal(fmt.Sprintf("error getting list of pods: %v", err))
	}

	var re = regexp.MustCompile(str1)
	var nspace string
	for _, pod := range podz.Items {
		for _, match := range re.FindAllString(pod.Name, -1) {
			if match != "" {
				nspace = pod.Namespace
			} else {
				nspace = "default"
			}
		}
	}

	path := "/apis/custom.metrics.k8s.io/v1beta1/namespaces/" + nspace + "/pods/*/" + mflag

	data, err := clientset.RESTClient().Get().AbsPath(path).DoRaw()
	if err != nil {
		return
	}
	var pods MetricValueList
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
