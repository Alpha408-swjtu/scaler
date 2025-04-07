package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"scaler/config"
	"scaler/hpa"
	"scaler/log"
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
	clearDir()

	desiredTime := 196400 //运行时间
	startTime := time.Now()

	for time.Since(startTime) < time.Duration(desiredTime)*time.Second {
		testTime := int(time.Since(startTime).Seconds())
		log.LogEntry.Infof("在%d秒执行操作", testTime)

		hpaExecuter := hpa.NewExecuterTmp(config.Client, "boutique", config.Apps, testTime)
		hpaExecuter.ExecuteAndSave()

		fmt.Println("本轮操作完成-----------------------------------------------------------------------------------------------")
		time.Sleep(30 * time.Second)
	}
}
