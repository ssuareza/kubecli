package ckill

import (
	"fmt"
	"log"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func kill(namespace string, label string) {
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

	// list pods = podLabel
	podLabel := "app=" + label

	pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: podLabel})
	if err != nil {
		log.Fatal(err)
	}

	// save list of pods
	if len(pods.Items) == 0 {
		fmt.Println("No pods found")
		os.Exit(0)
	}

	// get pods with some container in CrashLoopBackOff
	podList := []string{}
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if !containerStatus.Ready {
				if containerStatus.State.Waiting != nil {
					if containerStatus.State.Waiting.Reason == "CrashLoopBackOff" {
						podList = append(podList, pod.Name)
						break
					}
				}
			}
		}
	}

	if len(podList) == 0 {
		fmt.Println("No CrashLoopBackOff pods found")
		os.Exit(0)
	}

	for _, pod := range podList {
		if err := delete(clientset, namespace, pod); err != nil {
			log.Fatal(err)
		}
	}
}

// delete deletes a pod
func delete(client *kubernetes.Clientset, environment string, pod string) error {
	if err := client.CoreV1().Pods(environment).Delete(pod, &metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("error deleting /%s", pod)
	}

	fmt.Printf("pod/%s deleted\n", pod)

	return nil
}
