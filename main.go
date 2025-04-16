package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"scaler/config"
	"scaler/hpa"
	"time"
)

func clearDir() {
	dir := config.RootPath + "/hpa/data"
	// 确保文件夹存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("文件夹 %s 不存在\n", dir)
		return
	}

	// 遍历文件夹中的所有文件和子文件夹
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Printf("读取文件夹时出错: %v\n", err)
		return
	}

	// 遍历并删除每个文件或子文件夹
	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())

		// 如果是文件，直接删除
		if !file.IsDir() {
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("删除文件 %s 时出错: %v\n", filePath, err)
				continue
			}
			fmt.Printf("已删除文件: %s\n", filePath)
		} else {
			// 如果是目录，递归删除
			err := os.RemoveAll(filePath)
			if err != nil {
				fmt.Printf("删除目录 %s 时出错: %v\n", filePath, err)
				continue
			}
			fmt.Printf("已删除目录: %s\n", filePath)
		}
	}

	time.Sleep(5 * time.Second)
}

func main() {
	timer := time.After(59 * time.Second)
	fmt.Println(time.Now())
	s := hpa.GetHistoryMetrics(hpa.TransmittedQuery, "frontend", "boutique", 60, 1)
	fmt.Printf("原始时间序列为：%v\n", s)

	a := hpa.PreditLoad(s, 60)
	fmt.Printf("预测1m后的负载为:%v\n", a)

	<-timer
	f := hpa.GetQuery("frontend", "boutique", hpa.CurrentTransmittedQuery)
	fmt.Printf("一分钟后的实际负载是:%v", f)

}
