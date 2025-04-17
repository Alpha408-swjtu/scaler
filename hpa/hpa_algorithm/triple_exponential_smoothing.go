package hpaAlgorithm

func TripleExponentialSmoothing(data []float64, alpha float64) []float64 {
	n := len(data)
	if n == 0 || alpha <= 0 || alpha >= 1 {
		return []float64{}
	}

	// 初始化三层平滑值数组
	s1 := make([]float64, n)
	s2 := make([]float64, n)
	s3 := make([]float64, n)
	s1[0] = data[0]
	s2[0] = data[0]
	s3[0] = data[0]

	// 计算各时刻的平滑值
	for t := 1; t < n; t++ {
		s1[t] = alpha*data[t] + (1-alpha)*s1[t-1]
		s2[t] = alpha*s1[t] + (1-alpha)*s2[t-1]
		s3[t] = alpha*s2[t] + (1-alpha)*s3[t-1]
	}

	// 生成平滑后的结果
	smoothed := make([]float64, n)
	for t := 0; t < n; t++ {
		a := 3*s1[t] - 3*s2[t] + s3[t]
		bNumerator := alpha * ((6-5*alpha)*s1[t] - 2*(5-4*alpha)*s2[t] + (4-3*alpha)*s3[t])
		bDenominator := 2 * (1 - alpha) * (1 - alpha)
		b := bNumerator / bDenominator
		cNumerator := alpha * (s1[t] - 2*s2[t] + s3[t])
		c := cNumerator / bDenominator
		// 这里 T 取 0 表示当前时刻的预测值
		smoothed[t] = a + b*0 + c*0
	}

	return smoothed
}
