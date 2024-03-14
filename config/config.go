// Package config
// @title
// @description
// @author njy
// @since 2022/12/12 16:45
package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	print2 "metalab-run/utils/print"
	"os"
)

type tomlConf struct {
	Client  clientConf  `toml:"client"`
	App     appConf     `toml:"app"`
	Service serviceConf `toml:"service"`
}

type clientConf struct {
	ClientCode  string `toml:"client_code"`
	AppKey      string `toml:"app_key"`
	AppSecret   string `toml:"app_secret"`
	MqttPwd     string `toml:"mqtt_pwd"`
	DockerPwd   string `toml:"docker_pwd"`
	ServerIp    string `toml:"server_ip"`
	WxAppKey    string `toml:"wx_app_key"`
	WxAppSecret string `toml:"wx_app_secret"`
}

type appConf struct {
	MetalabPage metalabPageConf `toml:"metalab_page"`
	Metalab     metalabConf     `toml:"metalab"`
	Ffmpeg      ffmpegConf      `toml:"ffmpeg"`
}

type metalabPageConf struct {
	ConfPath string `toml:"conf_path"`
}

type metalabConf struct {
	ConfPath string `toml:"conf_path"`
	DataPath string `toml:"data_path"`
}

type ffmpegConf struct {
	ConfPath   string `toml:"conf_path"`
	FfmpegPath string `toml:"ffmpeg_path"`
}

type serviceConf struct {
	Mysql mysqlConf `toml:"mysql"`
	Redis redisConf `toml:"redis"`
	Minio minioConf `toml:"minio"`
	Srs   SrsConf   `toml:"srs"`
}

type mysqlConf struct {
	Port     string `toml:"port"`
	Pwd      string `toml:"pwd"`
	ConfPath string `toml:"conf_path"`
	DataPath string `toml:"data_path"`
	//LogPath  string `toml:"log_path"`
}

type redisConf struct {
	Port     string `toml:"port"`
	Pwd      string `toml:"pwd"`
	ConfPath string `toml:"conf_path"`
	DataPath string `toml:"data_path"`
	LogPath  string `toml:"log_path"`
}

type minioConf struct {
	Port     string `toml:"port"`
	WebPort  string `toml:"web_port"`
	Pwd      string `toml:"pwd"`
	ConfPath string `toml:"conf_path"`
	DataPath string `toml:"data_path"`
}

type SrsConf struct {
	ConfPath string `toml:"conf_path"`
}

var (
	Cfg         *tomlConf
	TemplateMap map[string]string
)

func init() {
	Cfg = new(tomlConf)
	_, err := toml.DecodeFile("./all_config.toml", &Cfg)
	if err != nil {
		print2.Red(fmt.Sprintf("load Docker config file failed: %v\n", err))
		os.Exit(1)
	}

	TemplateMap = make(map[string]string)
	TemplateMap["@{mysql.path.conf}"] = checkEmpty(Cfg.Service.Mysql.ConfPath)
	TemplateMap["@{mysql.path.data}"] = checkEmpty(Cfg.Service.Mysql.DataPath)
	TemplateMap["@{mysql.pwd}"] = checkEmpty(Cfg.Service.Mysql.Pwd)
	TemplateMap["@{mysql.port}"] = checkEmpty(Cfg.Service.Mysql.Port)

	TemplateMap["@{minio.path.data}"] = checkEmpty(Cfg.Service.Minio.DataPath)
	TemplateMap["@{minio.path.conf}"] = checkEmpty(Cfg.Service.Minio.ConfPath)
	TemplateMap["@{minio.pwd}"] = checkEmpty(Cfg.Service.Minio.Pwd)
	TemplateMap["@{minio.port.endpoint}"] = checkEmpty(Cfg.Service.Minio.Port)
	TemplateMap["@{minio.port.web}"] = checkEmpty(Cfg.Service.Minio.WebPort)

	TemplateMap["@{redis.path.data}"] = checkEmpty(Cfg.Service.Redis.DataPath)
	TemplateMap["@{redis.path.conf}"] = checkEmpty(Cfg.Service.Redis.ConfPath)
	TemplateMap["@{redis.path.log}"] = checkEmpty(Cfg.Service.Redis.LogPath)
	TemplateMap["@{redis.port}"] = checkEmpty(Cfg.Service.Redis.Port)
	TemplateMap["@{redis.pwd}"] = checkEmpty(Cfg.Service.Redis.Pwd)

	TemplateMap["@{srs.path.conf}"] = checkEmpty(Cfg.Service.Srs.ConfPath)

	TemplateMap["@{client.code}"] = checkEmpty(Cfg.Client.ClientCode)
	TemplateMap["@{client.appKey}"] = checkEmpty(Cfg.Client.AppKey)
	TemplateMap["@{client.appSecret}"] = checkEmpty(Cfg.Client.AppSecret)
	TemplateMap["@{client.mqttPwd}"] = checkEmpty(Cfg.Client.MqttPwd)
	TemplateMap["@{client.serverIp}"] = checkEmpty(Cfg.Client.ServerIp)
	TemplateMap["@{client.wxAppKey}"] = checkEmpty(Cfg.Client.WxAppKey)
	TemplateMap["@{client.wxAppSecret}"] = checkEmpty(Cfg.Client.WxAppSecret)

	TemplateMap["@{back.data}"] = checkEmpty(Cfg.App.Metalab.DataPath)
	TemplateMap["@{back.conf}"] = checkEmpty(Cfg.App.Metalab.ConfPath)
	TemplateMap["@{back.data}"] = checkEmpty(Cfg.App.Metalab.DataPath)

	TemplateMap["@{front.conf}"] = checkEmpty(Cfg.App.MetalabPage.ConfPath)

	TemplateMap["@{ffmpeg.conf}"] = checkEmpty(Cfg.App.Ffmpeg.ConfPath)
}

func checkEmpty(property string) string {
	if property == "" {
		print2.Red(fmt.Sprintf("load Docker config file failed: property cannot be empty\n"))
		os.Exit(1)
	}
	return property
}
