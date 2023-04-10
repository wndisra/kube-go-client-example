package main

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeClient, err := newClient()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Creating pod ...")

	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-created-by-client",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}
	opts := metav1.CreateOptions{}

	_, err = createPod(kubeClient, "default", pod, opts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Pod created successfully")

	log.Println("Getting pod ...")

	getOpts := metav1.GetOptions{}

	podDetail, err := getPod(kubeClient, "default", "pod-created-by-client", getOpts)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Pod namespace:", podDetail.Namespace)
	log.Println("Pod name:", podDetail.Name)
	log.Println("Pod read successfully")
}

func newClient() (*kubernetes.Clientset, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	k8sClientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	return k8sClientset, nil
}

func createPod(kubeClient *kubernetes.Clientset, namespace string, pod *corev1.Pod, opts metav1.CreateOptions) (*corev1.Pod, error) {
	if pod == nil {
		return nil, errors.New("pod can't be nil")
	}

	createdPod, err := kubeClient.CoreV1().Pods(namespace).Create(context.Background(), pod, opts)
	if err != nil {
		return nil, err
	}

	return createdPod, nil
}

func getPod(kubeClient *kubernetes.Clientset, namespace string, podName string, opts metav1.GetOptions) (*corev1.Pod, error) {
	pod, err := kubeClient.CoreV1().Pods(namespace).Get(context.Background(), podName, opts)
	if err != nil {
		return nil, err
	}

	return pod, nil
}
