package metrics

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Metrics interface {
	Get()
}

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

type MetricName struct {
	Name string
}

func NewMetrics(metricFlag string) Metrics {
	return MetricName{
		Name: metricFlag,
	}
}

func (m MetricName) Get() {
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

	podz, err := getPods(clientset)

	var re = regexp.MustCompile(str1)
	var nspace []string
	for _, pod := range podz.Items {
		match := re.MatchString(pod.Name)
		if match {
			nspace = append(nspace, pod.Namespace)
		} else {
			nspace = getNS(clientset)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Pod Name", "Namespace", "Metric", "Value"})

	for _, value := range nspace {
		path := "/apis/custom.metrics.k8s.io/v1beta1/namespaces/" + value + "/pods/*/" + mflag
		data, err := clientset.RESTClient().Get().AbsPath(path).DoRaw()
		if err != nil {
			fmt.Println("No resource found on namspace:", value)
		}
		var pods MetricValueList
		err = json.Unmarshal(data, &pods)
		if err != nil {
			fmt.Println(err)
		}

		for _, l := range pods.Items {
			data1 := []string{l.DescribedObject.Name, l.DescribedObject.Namespace, l.MetricName, l.Value}
			table.Append(data1)
		}

	}
	table.Render()
}

func getPods(clientset *kubernetes.Clientset) (podz *v1.PodList, err error) {
	podz, err = clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		fatal(fmt.Sprintf("error getting list of pods: %v", err))
	}
	return podz, err
}

func getNS(clientset *kubernetes.Clientset) (nsdata []string) {
	ns, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		fatal(fmt.Sprintf("error getting list of namespaces: %v", err))
	}
	//var nsdata []string
	for _, nsv := range ns.Items {
		nsdata = append(nsdata, nsv.Name)
	}
	return nsdata
}

func fatal(msg string) {
	os.Stderr.WriteString(msg + "\n")
	os.Exit(1)
}
