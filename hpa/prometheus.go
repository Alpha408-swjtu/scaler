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
	"time"

	"github.com/sirupsen/logrus"
)

var (
	qpsQuery         = `sum(rate(istio_requests_total{destination_workload_namespace="%s", destination_workload="%s"}[20s]))`
	receivedQuery    = `sum(rate(container_network_receive_bytes_total{namespace="%s", pod=~"%s-.*"}[100s]))/1024`
	TransmittedQuery = `sum(rate(container_network_transmit_bytes_total{namespace="%s", pod=~"%s-.*"}[100s]))/1024`
	readQuery        = `sum(rate(container_fs_reads_bytes_total{namespace="%s", pod=~"%s-.*", container!="POD"}[100s]))`
	writeQuery       = `sum(rate(container_fs_writes_bytes_total{namespace='%s', pod=~'%s-.*', container!='POD'}[100s]))`
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

// 查询qps
func GetQps(appName, namespace string, query string) []float64 {
	// 构造查询参数
	queryParams := url.Values{}
	queryParams.Add("query", fmt.Sprintf(query, namespace, appName))

	// 构造完整的 URL
	baseURL := config.PrometheusUrl + "/api/v1/query"
	url := baseURL + "?" + queryParams.Encode()

	fmt.Println(url)
	// 存储查询结果的数组
	results := make([]float64, 0, 4)

	// 查询 5 次
	for i := 0; i < 4; i++ {
		resp, err := http.Get(url)
		if err != nil {
			logrus.Errorf("Request error for %s (attempt %d): %v\n", appName, i+1, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			logrus.Errorf("Unexpected status code for %s (attempt %d): %d\n", appName, i+1, resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Failed to read response body for %s (attempt %d): %v\n", appName, i+1, err)
			continue
		}

		var qpsData Data
		if err := json.Unmarshal(body, &qpsData); err != nil {
			logrus.Errorf("Failed to parse JSON response for %s (attempt %d): %v\n", appName, i+1, err)
			continue
		}

		if len(qpsData.Data.Result) == 0 {
			logrus.Warnf("No data found for %s in namespace %s.\n", appName, namespace)
			results = append(results, 0.0)
			continue
		}

		qpsValues := make([]float64, 0, len(qpsData.Data.Result))

		for _, result := range qpsData.Data.Result {
			if len(result.Value) == 2 {
				if qpsValue, ok := result.Value[1].(string); ok {
					qps, _ := strconv.ParseFloat(qpsValue, 64)
					qpsValues = append(qpsValues, qps)
				}
			}
		}

		if len(qpsValues) == 0 {
			logrus.Warnf("No valid data for %s.\n", appName)
			results = append(results, 0.0)
		} else {
			// 计算平均值
			totalQPS := 0.0
			for _, value := range qpsValues {
				totalQPS += value
			}
			averageQPS := totalQPS / float64(len(qpsValues))
			results = append(results, averageQPS)
		}

		time.Sleep(15 * time.Second)
	}

	if len(results) == 0 {
		return nil
	}

	return results
}

// 从普罗米修斯接口获取其他数据
func getQuery(appName, namespace string, Query string) float64 {
	// 构造查询参数
	query := url.Values{}
	query.Add("query", fmt.Sprintf(Query, namespace, appName))

	// 构造完整的 URL
	baseURL := config.PrometheusUrl + "/api/v1/query"
	url := baseURL + "?" + query.Encode()

	//最多发送三次请求
	for attempt := 0; attempt < 3; attempt++ {
		resp, err := http.Get(url)
		if err != nil {
			mylog.Logger.Errorf("Request error for %s (attempt %d): %v\n", appName, attempt+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			mylog.Logger.Errorf("Unexpected status code for %s (attempt %d): %d\n", appName, attempt+1, resp.StatusCode)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(body))
			time.Sleep(2 * time.Second)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			mylog.Logger.Errorf("Failed to read response body for %s (attempt %d): %v\n", appName, attempt+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		var qpsData Data
		if err := json.Unmarshal(body, &qpsData); err != nil {
			mylog.Logger.Errorf("Failed to parse JSON response for %s (attempt %d): %v\n", appName, attempt+1, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(qpsData.Data.Result) == 0 {
			mylog.Logger.Warnf("No data found for %s in namespace %s.\n", appName, namespace)
		}

		qpsValues := make([]float64, 0, len(qpsData.Data.Result))

		for _, result := range qpsData.Data.Result {
			if len(result.Value) == 2 {
				if qpsValue, ok := result.Value[1].(string); ok {
					qps, _ := strconv.ParseFloat(qpsValue, 64)
					qpsValues = append(qpsValues, qps)
				}

			}
		}

		if len(qpsValues) == 0 {
			mylog.Logger.Warnf("No valid CPU data for %s.\n", appName)
			return 0
		}

		// 计算平均值
		totalQPS := 0.0
		for _, value := range qpsValues {
			totalQPS += value
		}

		averageQPS := totalQPS / float64(len(qpsValues))
		return averageQPS
	}

	mylog.Logger.Errorf("Failed to fetch QPS usage for %s after 3 attempts.\n", appName)
	return 0
}
