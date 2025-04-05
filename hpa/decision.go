package hpa

import (
	"math"
	"scaler/config"
	"scaler/log"

	"k8s.io/client-go/kubernetes"
)

// Monitor组件 收集微服务监控数据，以及返回设置的默认参数
func monitor(client *kubernetes.Clientset, appName, namespace string) (int, int, config.Parameter) {
	//获取微服务监控数据
	currentReplicas := getReadyReplicas(client, appName, namespace)
	desiredReplicas := getDesiredReplicas(client, appName, namespace)
	//根据应用名称 从解析的配置文件中获取默认参数
	parameter := config.Parameters[appName]

	return desiredReplicas, currentReplicas, parameter
}

// 计算出微服务的期望副本数，并生成扩缩容决策
func analyse(appName string, desiredReplicas int, currentReplicas int, currentQPS float64, targetQPS float64, minReplicas int) (string, int) {
	// 计算期望副本数
	if desiredReplicas != currentReplicas {
		log.Logger.Warnf("目标副本和就绪副本不一致!")
	}

	desiredReplica := int(math.Ceil(currentQPS / targetQPS))
	if desiredReplica < minReplicas {
		desiredReplica = minReplicas
	}
	log.Logger.Debugf("计算出期%s望副本数为: %d", appName, desiredReplica)

	// 生成扩缩容决策
	var scalingAction string
	if desiredReplica > currentReplicas {
		scalingAction = "扩容"
	} else if desiredReplica < currentReplicas {
		scalingAction = "缩容"
	} else {
		scalingAction = "不变"
	}
	log.Logger.Debugf("%s的伸缩策略: %s", appName, scalingAction)

	return scalingAction, desiredReplica
}
