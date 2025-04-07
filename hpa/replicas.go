package hpa

import (
	"context"
	"fmt"
	mylog "scaler/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// 获取指定deploy的就绪副本数
func getReadyReplicas(client *kubernetes.Clientset, appName, namespace string) int {
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app=%s", appName),
	})
	if err != nil {
		mylog.Logger.Errorf("获取pod:%s失败:%v\n", appName, err)
		return 0
	}

	readyCount := 0
	for _, pod := range pods.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == appName && containerStatus.Ready {
				readyCount++
			}
		}
	}

	mylog.Logger.Infof("%s的就绪副本数为:%v", appName, readyCount)
	return readyCount
}

// 获取当前期望副本数
func getDesiredReplicas(client *kubernetes.Clientset, appName, namespace string) int {
	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), appName, metav1.GetOptions{})
	if err != nil {
		mylog.Logger.Errorf("获取deploy期望副本:%s失败:%v\n", appName, err)
		return 1
	}

	desiredReplicas := deployment.Spec.Replicas
	mylog.LogEntry.Logger.Infof("deployment:%s的正确副本数为:%v", appName, *desiredReplicas)
	return int(*desiredReplicas)
}
