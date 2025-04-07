package config

import (
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 存储的全局变量
var (
	//k8s配置文件路径
	FilePath string
	//项目根目录
	RootPath string
	//扩容冷却时间
	COOLDOWN_UP int64
	//缩容冷却时间
	COOLDOWN_DOWN int64
	//普罗米修斯路径
	PrometheusUrl string
	//微服务应用名称
	Apps []string
	//微服务应用的默认设置
	Parameters map[string]Parameter

	Client *kubernetes.Clientset
)

type Config struct {
	File       FileConfig
	Root       RootConfig
	ColdUp     ColdupConfig
	Coldown    ColdownConfig
	Url        UrlConfig
	Apps       AppsConfig
	Parameters map[string]Parameter
}

type FileConfig struct {
	ConfigPath string
}

type RootConfig struct {
	RootPath string
}

type ColdupConfig struct {
	COOLDOWN_UP int64
}

type ColdownConfig struct {
	COOLDOWN_DOWN int64
}

type UrlConfig struct {
	PrometheusUrl string
}

type AppsConfig struct {
	AppNames []string
}

type Parameter struct {
	DEFAULT_MIN_REPLICAS int
	DEFAULT_CPU_REQUEST  float64
	TARGET_QPS           float64
	MAX_REPLICAS         int
}

// 初始化配置，解析配置文件的过程
func init() {
	//error := 0
	var config Config
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.Unmarshal(&config)

	FilePath = config.File.ConfigPath
	RootPath = config.Root.RootPath
	COOLDOWN_UP = config.ColdUp.COOLDOWN_UP
	COOLDOWN_DOWN = config.Coldown.COOLDOWN_DOWN
	PrometheusUrl = config.Url.PrometheusUrl
	Apps = config.Apps.AppNames
	Parameters = config.Parameters

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", FilePath)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err)
	}

	Client = client
}
