package exec

import (
	"fmt"
	"os"
	"log"

	"golang.org/x/crypto/ssh/terminal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/tools/clientcmd"
)

func exec(namespace string, label string) {
	log.SetFlags(0)

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

	podName := pods.Items[0].Name

	// exec
	cmd := []string{
		"/bin/sh",
	}

	req := clientset.CoreV1().RESTClient().
		Post().Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: label,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	// Put the terminal into raw mode to prevent it echoing characters twice.
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(0, oldState)

	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
	if err != nil {
		fmt.Println(err)
	}
}
