package hpaAlgorithm

// Das 结构体
type Das struct{}

// 预测负载
func (d *Das) PredictLoad(data []float64, T int) float64 {
	m1 := movingAverage(data)
	firstDatas := firstMovingSequence(data)
	m2 := movingAverage(firstDatas)

	a := 2*m1 - m2
	b := 2 * (m1 - m2) / float64(len(data)-1)

	return a + b*float64(T)
}

// 一次移动的平均值
func movingAverage(data []float64) float64 {
	sum := float64(0)
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

// 一次移动后的时间序列
func firstMovingSequence(data []float64) []float64 {
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

// 二次移动平均法，返回完整的平滑序列
func DoubleMovingAverage(data []float64) []float64 {
	n := len(data)
	if n == 0 {
		return []float64{}
	}

	// 一次移动平均
	firstDatas := firstMovingSequence(data)

	// 二次移动平均
	secondDatas := firstMovingSequence(firstDatas)

	// 计算 a 和 b
	a := make([]float64, n)
	b := make([]float64, n)

	for i := 0; i < n; i++ {
		a[i] = 2*firstDatas[i] - secondDatas[i]
		b[i] = 2 * (firstDatas[i] - secondDatas[i]) / float64(n-1)
	}

	// 计算预测值
	predicted := make([]float64, n)
	for i := 0; i < n; i++ {
		predicted[i] = a[i] + b[i]*float64(i)
	}

	return predicted
}
