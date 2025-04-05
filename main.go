package main

import (
	"fmt"
	"scaler/config"
	"scaler/hpa"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", config.FilePath)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	s := hpa.NewHpa(client, "boutique", "frontend")
	fmt.Println(s)
}
