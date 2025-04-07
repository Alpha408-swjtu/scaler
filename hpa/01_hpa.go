package hpa

import (
	"context"
	"scaler/log"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type MicroserviceData struct {
	AppName         string // 微服务名称
	ScalingAction   string // 伸缩决策
	DesiredReplicas int    // 期望副本数
	CurrentReplicas int    // 当前副本数
	//CPURequest      float64 // CPU 请求
	MaxReplicas int // 最大副本数决策
	CurrentQPS  float64
	//CpuPercentage   float64 //cpu使用率
}

type Standards struct {
	QPS         float64 //连接数
	Recived     float64 //入口带宽
	Transmitted float64 //出口带宽
	Iops        float64 //磁盘读写
}

type AppInfo struct {
	Namespace string
	AppName   string
}

type Hpa struct {
	Client     *kubernetes.Clientset
	Deployment *appsv1.Deployment
	AppInfo    AppInfo
	Standard   Standards
}

func NewHpa(client *kubernetes.Clientset, namespace, appName string) *Hpa {
	deploy, err := client.AppsV1().Deployments(namespace).Get(context.Background(), appName, metav1.GetOptions{})
	if err != nil {
		log.Logger.Panicf("获取deployment失败:%v", err)
	}

	return &Hpa{
		Client:     client,
		Deployment: deploy,
		AppInfo: AppInfo{
			Namespace: namespace,
			AppName:   appName,
		},
		Standard: Standards{
			QPS:         getQuery(appName, namespace, qpsQuery),
			Recived:     getQuery(appName, namespace, receivedQuery),
			Transmitted: getQuery(appName, namespace, transmittedQuery),
			Iops:        0,
		},
	}
}

func (h *Hpa) MicroData(testTime int) MicroserviceData {
	desiredReplicas, currentReplicas, parameter := monitor(h.Client, h.AppInfo.AppName, h.AppInfo.Namespace)
	scalingAction, desiredReplica := analyse(h.AppInfo.AppName, desiredReplicas, currentReplicas, h.Standard.QPS, parameter.TARGET_QPS, parameter.DEFAULT_MIN_REPLICAS)
	return MicroserviceData{
		AppName:         h.AppInfo.AppName,
		ScalingAction:   scalingAction,
		DesiredReplicas: desiredReplica,
		CurrentReplicas: currentReplicas,
		MaxReplicas:     parameter.MAX_REPLICAS,
		CurrentQPS:      h.Standard.QPS,
	}
}

func (h *Hpa) Scale(desiredReplicas int) error {
	*h.Deployment.Spec.Replicas = int32(desiredReplicas)
	if _, err := h.Client.AppsV1().Deployments(h.AppInfo.Namespace).Update(context.TODO(), h.Deployment, metav1.UpdateOptions{}); err != nil {
		return err
	}
	return nil
}
