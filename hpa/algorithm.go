package hpa

import (
	"time"
)

const (
	windowSize       = 10 //取10组数据预测
	predictionPeriod = 60 //预测60s之后的数据
	CollectInterval  = 6 * time.Second
)

var HistoricalData = make([]float64, 0, windowSize+1)

// 一次移动的平均值
func MovingAverage(data []float64) float64 {
	sum := float64(0)
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// 一次移动后的时间序列
func FirstMovingSequence(data []float64) []float64 {
	n := len(data)
	if n == 0 {
		return []float64{}
	}

	result := make([]float64, n)
	sum := float64(0)

	// 从数组末尾开始计算
	for i := n - 1; i >= 0; i-- {
		sum += data[i]
		result[i] = sum / float64(n-i)
	}

	return result
}

// 预测负载
func PreditLoad(datas []float64, T int) float64 {
	m1 := MovingAverage(datas)
	firstDatas := FirstMovingSequence(datas)
	m2 := MovingAverage(firstDatas)

	a := 2*m1 - m2
	b := 2 * (m1 - m2) / float64(len(datas)-1)

	return a + b*float64(T)
}
