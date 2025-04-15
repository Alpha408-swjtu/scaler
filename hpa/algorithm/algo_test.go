package hpa_algorithm

import (
	"fmt"
	"testing"
)

// singleMovingAverage 计算一次移动平均值
func singleMovingAverage(data []float64, n int) []float64 {
	result := make([]float64, len(data))
	for i := n - 1; i < len(data); i++ {
		sum := 0.0
		for j := 0; j < n; j++ {
			sum += data[i-j]
		}
		result[i] = sum / float64(n)
	}
	return result
}

// doubleMovingAverage 计算二次移动平均值
func doubleMovingAverage(data []float64, n int) []float64 {
	singleMA := singleMovingAverage(data, n)
	return singleMovingAverage(singleMA, n)
}

// predict 进行预测
func predict(data []float64, n int, T int) float64 {
	singleMA := singleMovingAverage(data, n)
	doubleMA := doubleMovingAverage(data, n)
	t := len(data) - 1
	a := 2*singleMA[t] - doubleMA[t]
	b := (2 / float64(n-1)) * (singleMA[t] - doubleMA[t])
	return a + b*float64(T)
}

// calculateDynamicThreshold 计算动态阈值
func calculateDynamicThreshold(fx, historicalAvg, loadGrowthRate float64) (float64, float64) {
	expansionThreshold := fx - historicalAvg*loadGrowthRate
	shrinkThreshold := fx + historicalAvg*loadGrowthRate
	return expansionThreshold, shrinkThreshold
}

func TestA(t *testing.T) {
	// 模拟每 15 秒采集一次指标，一分钟获取 4 个数据
	cpuLoadData := []float64{0.1, 0.7, 0.2, 0.1}
	// 期数
	n := 3
	// 预测未来 1 个周期（1 分钟）的资源需求
	T := 1
	prediction := predict(cpuLoadData, n, T)
	fmt.Printf("预测未来 %d 个周期的 CPU 负载为: %.2f\n", T, prediction)

	// 静态扩容阈值和缩容阈值
	alpha := 0.3
	beta := 0.7
	// 假设 CPU 负载增长率和历史平均值
	loadGrowthRate := 0.1
	historicalAvg := 0.3
	// 当前负载
	fx := 0.5

	// 计算动态阈值
	expansionThreshold, shrinkThreshold := calculateDynamicThreshold(fx, historicalAvg, loadGrowthRate)
	fmt.Printf("动态扩容阈值: %.2f, 动态缩容阈值: %.2f\n", expansionThreshold, shrinkThreshold)

	// 模拟当前负载情况，判断是否需要扩缩容
	currentLoad := 0.6
	if currentLoad >= expansionThreshold && currentLoad < beta {
		fmt.Println("触发扩容操作")
	} else if currentLoad <= shrinkThreshold && currentLoad > alpha {
		fmt.Println("触发缩容操作")
	} else {
		fmt.Println("当前负载正常，无需扩缩容")
	}
}
