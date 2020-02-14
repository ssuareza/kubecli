package get

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func getObjects(clientset kubernetes.Interface, namespace string, label string) (*uitable.Table, error) {
	selector := "app=" + label

	// table
	table := uitable.New()
	table.MaxColWidth = 150
	table.Wrap = true
	table.AddRow(color.YellowString("Name"), color.YellowString("Labels"), color.YellowString("Status"))

	// list services
	var deploy string
	services, err := clientset.CoreV1().Services(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}
	if len(services.Items) != 0 {
		for _, service := range services.Items {
			table.AddRow(color.GreenString("service/"+service.Name), mapToString(service.Labels), mapToString(service.Spec.Selector))
			deploy = service.Spec.Selector["deploy"]
		}
	}

	// list jobs
	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}

	if len(jobs.Items) != 0 {
		for _, job := range jobs.Items {
			if job.Labels["deploy"] == deploy {
				var jobStatus string
				if job.Status.Succeeded == 1 {
					jobStatus = "Succeeded"
				} else if job.Status.Failed == 1 {
					jobStatus = "Failed"
				} else if job.Status.Active == 1 {
					jobStatus = "Active"
				}
				table.AddRow(color.GreenString("job/"+job.Name), mapToString(job.Labels), jobStatus)
				continue
			}
			table.AddRow("job/"+job.Name, mapToString(job.Labels), "")
		}
	}

	// list configmaps
	configmaps, err := clientset.CoreV1().ConfigMaps(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}

	if len(configmaps.Items) != 0 {
		for _, configmap := range configmaps.Items {
			if configmap.Labels["deploy"] == deploy {
				table.AddRow(color.GreenString("configmap/"+configmap.Name), mapToString(configmap.Labels), "")
				continue
			}
			table.AddRow("configmap/"+configmap.Name, mapToString(configmap.Labels), "")
		}
	}

	// list secrets
	secrets, err := clientset.CoreV1().Secrets(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}

	if len(secrets.Items) != 0 {
		for _, secret := range secrets.Items {
			if secret.Labels["deploy"] == deploy {
				table.AddRow(color.GreenString("secret/"+secret.Name), mapToString(secret.Labels), "")
				continue
			}
			table.AddRow("secret/"+secret.Name, mapToString(secret.Labels), "")
		}
	}

	// list pods
	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: selector})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) != 0 {
		for _, pod := range pods.Items {
			if pod.Labels["deploy"] == deploy {
				if pod.Status.Phase == "Running" || pod.Status.Phase == "Succeeded" {
					table.AddRow(color.GreenString("pod/"+pod.Name), mapToString(pod.Labels), pod.Status.Phase)
				} else {
					table.AddRow(color.RedString("pod/"+pod.Name), mapToString(pod.Labels), pod.Status.Phase)
				}
				continue
			}
			fmt.Printf("pod/%s %sstatus=%s\n", pod.Name, mapToString(pod.Labels), pod.Status.Phase)
		}
	}

	return table, nil
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
		if strings.HasPrefix(v, "controller-uid") || strings.HasPrefix(v, "pod-template-hash") {
			continue
		}
		str += v + " "
	}
	return str
}
