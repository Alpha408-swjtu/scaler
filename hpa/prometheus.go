package hpa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"scaler/config"
	mylog "scaler/log"
	"strconv"
)

var (
	QpsQuery      = `sum(rate(istio_requests_total{destination_workload_namespace="%s", destination_workload="%s"} [1m] offset %ds))`
	CurrentQps    = `sum(rate(istio_requests_total{destination_workload_namespace="%s", destination_workload="%s"} [1m] ))`
	receivedQuery = `sum(rate(container_network_receive_bytes_total{namespace="%s", pod=~"%s-.*"}[100s]))/1024`
	//TransmittedQuery = `sum(container_network_transmit_bytes_total{namespace="%s", pod=~"%s-.*"})/1024/1024`
	CurrentTransmittedQuery = `sum(rate(container_network_transmit_bytes_total{namespace="%s", pod=~"%s-.*"}[1m] ))/1024`
	TransmittedQuery        = `sum(rate(container_network_transmit_bytes_total{namespace="%s", pod=~"%s-.*"}[1m] offset %ds))/1024`
	readQuery               = `sum(rate(container_fs_reads_bytes_total{namespace="%s", pod=~"%s-.*", container!="POD"}[100s]))`
	writeQuery              = `sum(rate(container_fs_writes_bytes_total{namespace='%s', pod=~'%s-.*', container!='POD'}[100s]))`
)

type Data struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func GetQuery(appName, namespace, query string) float64 {
	// 构造查询参数
	queryParams := url.Values{}
	queryParams.Add("query", fmt.Sprintf(query, namespace, appName))

	// 构造完整的 URL
	baseURL := config.PrometheusUrl + "/api/v1/query"
	url := baseURL + "?" + queryParams.Encode()

	resp, err := http.Get(url)
	if err != nil {
		mylog.Logger.Errorf("获取资源有误:%s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Logger.Errorf("解析有误:%s", err)
	}

	var data Data
	err = json.Unmarshal(body, &data)
	if err != nil {
		mylog.Logger.Errorf("反序列化有误:%s", err)
	}

	if len(data.Data.Result) > 0 {
		s := data.Data.Result[0].Value[1].(string)
		a, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Println("获取指标有误:", err)
		}
		return a
	}

	return -1
}

// 获取普罗米修斯的历史数据
func GetHistoryMetrics(query, appName, namespace string, duration int, step int) []float64 {
	if duration%step != 0 {
		panic("参数非法")
	}

	result := make([]float64, 0, duration/step)

	for i := step; i < duration; i += step {
		queryParams := url.Values{}
		queryParams.Add("query", fmt.Sprintf(query, namespace, appName, i))
		baseURL := config.PrometheusUrl + "/api/v1/query"
		url := baseURL + "?" + queryParams.Encode()

		resp, err := http.Get(url)
		if err != nil {
			mylog.Logger.Errorf("获取资源有误:%s", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mylog.Logger.Errorf("解析有误:%s", err)
		}

		var data Data
		err = json.Unmarshal(body, &data)
		if err != nil {
			mylog.Logger.Errorf("反序列化有误:%s", err)
		}

		//获取信息
		if len(data.Data.Result) > 0 {
			s := data.Data.Result[0].Value[1].(string)
			a, err := strconv.ParseFloat(s, 64)
			if err != nil {
				fmt.Println("获取指标有误:", err)
			}
			//fmt.Printf("时间:%v之前获取到的指标:%v\n", i, a)
			result = append(result, a)
		}

	}
	return result
}
