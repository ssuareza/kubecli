package get

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetObjects(t *testing.T) {
	// fake client
	clientset := fake.NewSimpleClientset()

	// create fake pod
	_, err := clientset.CoreV1().Pods("testing").Create(
		&v1.Pod{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Pod",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "fakepod",
				Labels: map[string]string{
					"app": "fakelabel",
				},
			},
		})

	table, err := getObjects(clientset, "testing", "fakelabel")
	if err != nil || len(table.Rows) != 2 {
		t.Fail()
	}
}
