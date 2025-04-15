#  scaler

算网融合的k8s集群扩缩容机制

## hpa算法

​        默认HPA 算法必须等到 Pod 的负载达到指定的扩容阈值时，才会触发扩容机制，并且集群生成新的 Pod 需要一定时间，而如果此时遇到突发的大流量请求就很有可能会导致集群扩容不及时而造成响应时间过长。而在集群访问流量骤降时 HPA 算法又会迅速的进行缩容，当遇到第二次突发大流量时，刚缩容完的 Pod 又要马上进行扩容，如此频繁 的伸缩也会导致集群性能受损。

​       拟设计一种既能根据历史数据预测未来负载情况，又能在预测间隔期应对突发流量的算法。

### 算法思路

​       ![01](./images/01.png)

### 算法内容

#### 1.hpa指标采集

拟采集deployment的qps、入口带宽、出口带宽、磁盘读、写速率作为指标

```go
var (
	qpsQuery         = `sum(rate(istio_requests_total{destination_workload_namespace="%s", destination_workload="%s"}[1m]))`
	receivedQuery    = `sum(rate(container_network_receive_bytes_total{namespace="%s", pod=~"%s-.*"}[1m]))/1024`
	TransmittedQuery = `sum(rate(container_network_transmit_bytes_total{namespace="%s", pod=~"%s-.*"}[1m]))/1024`
	readQuery        = `sum(rate(container_fs_reads_bytes_total{namespace="%s", pod=~"%s-.*", container!="POD"}[1m]))`
	writeQuery       = `sum(rate(container_fs_writes_bytes_total{namespace='%s', pod=~'%s-.*', container!='POD'}[1m]))`
)
```

每10s采集一次，1m为一个预测周期，即每次预测1m内的数据,避免采集过于频繁增大Prometheus pod的压力

#### 2.基于二次移动平均法(DMA)预测

![02](./images/02.png)

代码中，n为窗口大小，T为预测未来T个周期(10s)内的数据大小。

```go
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
```

#### 3.动态阈值调整

​     基于动态阈值的容器云弹性伸缩主要是通过动态下降算法来调整扩缩容的阈值。

​                                           ![03](./images/03.png)

```go
// calculateDynamicThreshold 计算动态阈值
func calculateDynamicThreshold(fx, historicalAvg, loadGrowthRate float64) (float64, float64) {
	expansionThreshold := fx - historicalAvg*loadGrowthRate
	shrinkThreshold := fx + historicalAvg*loadGrowthRate
	return expansionThreshold, shrinkThreshold
}
```

#### 4 核心策略

​         通过二次移动平均法预测未来资源需求，动态下调扩容阈值 *T*expansion，在预测间隔期（默认 60 秒）内，当负载达到 *T*expansion 时提前创建新 Pod，避免突发流量导致的服务质量下降。



## 后续待解决

1.依据任务类型给微服务分类:计算敏感，存储敏感，网络敏感，不同类型的参数选择选择。

![04](./images/04.jpg)

​                                                                                    依据cpu使用率的扩缩容后，可看cpu消耗情况

![05](./images/05.jpg)

​                                                                                                              qps情况

2.多重指标加入算法的改进策略，综合得到期望副本数。