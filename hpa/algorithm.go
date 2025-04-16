package hpa

import (
	"fmt"
	"time"
)

const (
	windowSize       = 10 //取10组数据预测
	predictionPeriod = 60 //预测60s之后的数据
	CollectInterval  = 6 * time.Second
)

var HistoricalData = make([]float64, 0, windowSize+1)

// 一次移动的平均值
func FirstMovingAverage(data []float64) float64 {
	sum := float64(0)
	for _, v := range data {
		sum += v
	}
	fmt.Printf("一次平均移动值为:%v\n", sum/float64(len(data)))
	return sum / float64(len(data))
}

// 计算二次移动的平均值
func SecondMovingAverage(data []float64) float64 {
	firstMovingAverages := make([]float64, len(data)-windowSize+1)
	for i := 0; i <= len(data)-windowSize; i++ {
		firstMovingAverages[i] = FirstMovingAverage(data[i : i+windowSize])
	}

	s := FirstMovingAverage(firstMovingAverages)
	fmt.Printf("二次平均移动值为:%v\n", s)
	return s
}

// 预测负载
func PredictLoad() float64 {
	if len(HistoricalData) < windowSize {
		return 0
	}

	m1 := FirstMovingAverage(HistoricalData)
	m2 := SecondMovingAverage(HistoricalData)

	a := 2*m1 - m2
	b := (2 / float64(windowSize-1)) * (m1 - m2)

	return a + b*float64(predictionPeriod/CollectInterval.Seconds())
}
