package hpa

import (
	"fmt"
	"os"
	"scaler/config"
	"scaler/log"
	"sync"

	"github.com/xuri/excelize/v2"
	"k8s.io/client-go/kubernetes"
)

type HpaExecuter struct {
	TestTime int
	HpaSlice []*Hpa
	DataMp   map[string]MicroserviceData
	lock     sync.RWMutex
}

func NewExecuter(client *kubernetes.Clientset, namespace string, appNames []string, testTime int) *HpaExecuter {
	if len(appNames) <= 0 {
		log.Logger.Panicln("待操作的微服务数量有误!")
	}

	hpaSlice := make([]*Hpa, 0, len(appNames))
	dataMp := make(map[string]MicroserviceData, len(appNames))

	for _, app := range appNames {
		hpa := NewHpa(client, namespace, app)
		microData := hpa.MicroData(testTime)

		hpaSlice = append(hpaSlice, hpa)
		dataMp[app] = microData
	}

	return &HpaExecuter{
		TestTime: testTime,
		HpaSlice: hpaSlice,
		DataMp:   dataMp,
		lock:     sync.RWMutex{},
	}
}

// 并发执行操作，保存文件
func (he *HpaExecuter) ExecuteAndSave() {
	wg := sync.WaitGroup{}
	wg.Add(len(he.HpaSlice))
	for _, hpa := range he.HpaSlice {
		go func() {
			defer wg.Done()
			//加读锁
			he.lock.RLock()
			data, ok := he.DataMp[hpa.AppInfo.AppName]
			//解锁
			he.lock.RUnlock()
			if !ok {
				return
			}
			saveDataFile(data, he.TestTime)

			if data.ScalingAction != "不变" {
				he.lock.Lock()
				defer he.lock.Unlock()
				hpa.Scale(data.DesiredReplicas)
			}

			log.Logger.Infof("应用:%s操作执行完毕", hpa.AppInfo.AppName)
		}()
	}
	wg.Wait()
}

// 保存应用hpa信息数据到本地文件
func saveDataFile(data MicroserviceData, testTime int) error {
	appName := data.AppName
	desiredReplicas := data.DesiredReplicas
	currentReplicas := data.CurrentReplicas
	scalingAction := data.ScalingAction
	currentQPS := data.CurrentQPS

	filePath := config.RootPath + "/hpa/data/" + appName + ".xlsx"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，创建新文件
		f := excelize.NewFile()
		// 设置表头
		f.SetCellValue("Sheet1", "A1", "Test Time (sec)")
		f.SetCellValue("Sheet1", "B1", "QPS")
		f.SetCellValue("Sheet1", "C1", "Current Replicas")
		f.SetCellValue("Sheet1", "D1", "Desired Replicas")
		f.SetCellValue("Sheet1", "E1", "Max.Replicas")
		f.SetCellValue("Sheet1", "F1", "Scaling Action")
		// 保存新文件
		if err := f.SaveAs(filePath); err != nil {
			log.Logger.Errorf("保存文件失败:%s", err)
			return err
		}
	} else {
		// 文件存在，打开文件
		f, err := excelize.OpenFile(filePath)
		if err != nil {
			return err
		}

		// 获取当前行数
		rows, _ := f.GetRows("Sheet1")
		rowCount := len(rows)

		// 写入数据
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowCount+1), testTime)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowCount+1), currentQPS)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowCount+1), currentReplicas)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowCount+1), desiredReplicas)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowCount+1), data.MaxReplicas)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowCount+1), scalingAction)

		// 保存文件
		if err := f.Save(); err != nil {
			return err
		}
	}

	return nil
}
