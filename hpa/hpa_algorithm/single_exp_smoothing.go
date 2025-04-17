package hpaAlgorithm

// exponentialSmoothing 函数实现指数平滑法，返回预测值数组
func SingleExpSmoothing(data []float64, alpha float64) []float64 {
	if len(data) == 0 {
		return []float64{}
	}
	// 初始化平滑值数组
	smoothed := make([]float64, len(data))
	// 初始平滑值取第一个数据点
	smoothed[0] = data[0]

	for i := 1; i < len(data); i++ {
		smoothed[i] = alpha*data[i] + (1-alpha)*smoothed[i-1]
	}
	return smoothed
}
