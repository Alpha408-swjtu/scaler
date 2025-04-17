package file

import (
	"bufio"
	"fmt"
	"os"
	"scaler/config"
	"strconv"
	"strings"
)

func ReadFile(name string) []float64 {
	file, err := os.Open(config.RootPath + "/file/" + name)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return nil
	}
	defer file.Close()

	// 创建一个切片来存储数据
	var data []float64

	// 使用 bufio 逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text()) // 去掉行首尾的空白字符
		if line == "" {
			continue // 跳过空行
		}

		// 将字符串转换为浮点数
		value, err := strconv.ParseFloat(line, 64)
		if err != nil {
			fmt.Printf("无法解析数据: %v\n", err)
			continue
		}

		// 将解析后的数据添加到切片中
		data = append(data, value)
	}

	// 检查读取过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时出错: %v\n", err)
		return nil
	}

	return data
}

func SaveFile(data []float64, name string) {
	// 打开文件，如果文件不存在则创建
	file, err := os.Create(config.RootPath + "/file/" + name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 写入数据
	for _, value := range data {
		_, err := file.WriteString(fmt.Sprintf("%f\n", value))
		if err != nil {
			panic(err)
		}
	}
}
