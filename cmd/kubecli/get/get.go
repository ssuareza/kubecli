package get

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/fatih/color"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getObjects(namespace string, label string) {
	log.SetFlags(0)
	selector := "app=" + label

	// kubeconfig
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// list services
	var deploy string
	services, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		log.Fatal(err)
	}
	if len(services.Items) != 0 {
		for _, service := range services.Items {
			fmt.Printf(color.GreenString("service/%s\n"), service.Name)
			deploy = service.Spec.Selector["deploy"]
		}
	}

	// list jobs
	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		log.Fatal(err)
	}

	if len(jobs.Items) != 0 {
		for _, job := range jobs.Items {
			if job.Labels["deploy"] == deploy {
				fmt.Printf(color.GreenString("job/%s\n"), job.Name)
				continue
			}
			fmt.Printf(color.RedString("job/%s\n"), job.Name)
		}
	}

	// list pods
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		log.Fatal(err)
	}

	if len(pods.Items) != 0 {
		fmt.Printf("\n")
		for _, pod := range pods.Items {
			if pod.Labels["deploy"] == deploy {
				if pod.Status.Phase == "Running" || pod.Status.Phase == "Succeeded" {
					fmt.Printf(color.GreenString("pod/%s")+" %sstatus=%s\n", pod.Name, mapToString(pod.Labels), pod.Status.Phase)
				} else {
					fmt.Printf(color.RedString("pod/%s")+" %sstatus=%s\n", pod.Name, mapToString(pod.Labels), pod.Status.Phase)
				}
				continue
			}
			fmt.Printf("pod/%s %sstatus=%s\n", pod.Name, mapToString(pod.Labels), pod.Status.Phase)
		}
	}

}

func mapToString(m map[string]string) string {
	var str string
	mk := make([]string, len(m))
	i := 0
	for k, v := range m {
		mk[i] = k + ":" + v
		i++
	}
	sort.Strings(mk)
	for _, v := range mk {
		str += v + " "
	}
	return str
}
