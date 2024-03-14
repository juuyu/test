package service

import (
	"fmt"
	"metalab-run/config"
	"metalab-run/deploy"
	"metalab-run/deploy/docker"
	"metalab-run/utils/cmd"
	"metalab-run/utils/http"
	print2 "metalab-run/utils/print"
	"os"
	"path"
	"time"
)

func FirstDeploy() {
	print2.Blue("=======================Start Install Service=========================\n")
	composeFilename := path.Base(config.Url.ServerComposeTemplate)
	print2.Yellow("start download and replace server-docker-compose.yml\n")
	deploy.DownloadAndReplaceTemplate(config.Url.ServerComposeTemplate)

	print2.Yellow("start download and replace srs conf file\n")
	deploy.DownloadAndReplaceTemplate(config.Url.SrsConfPath)
	confPath := config.Cfg.Service.Srs.ConfPath
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		if err = os.MkdirAll(confPath, os.ModePerm); err != nil {
			print2.Red("make srs conf path failed\n")
			os.Exit(1)
		}
	}
	SrsConfFilename := path.Base(config.Url.SrsConfPath)
	cmd.Exec("mv", "./"+SrsConfFilename, confPath)
	print2.Yellow("start download SSL Certificate\n")
	if KeyFilename, err := http.Download(fmt.Sprintf(config.Url.FrontKeyPath, config.Cfg.Client.ClientCode)); err != nil {
		print2.Red("download SSL key failed\n")
	} else if err := cmd.ExecNotExit("mv", "./"+KeyFilename, confPath); err != nil {
		print2.Red("move front SSL key to conf path failed\n")
	} else if err := cmd.ExecNotExit("mv", confPath+"/"+KeyFilename, confPath+"/"+"i-metalab.com.key"); err != nil {
		print2.Red("rename SSL key failed\n")
	}
	if PemFilename, err := http.Download(fmt.Sprintf(config.Url.FrontPemPath, config.Cfg.Client.ClientCode)); err != nil {
		print2.Red("download SSL pem failed\n")
	} else if err := cmd.ExecNotExit("mv", "./"+PemFilename, confPath); err != nil {
		print2.Red("move front SSL pem to conf path failed\n")
	} else if err := cmd.ExecNotExit("mv", confPath+"/"+PemFilename, confPath+"/"+"i-metalab.com.pem"); err != nil {
		print2.Red("rename SSL pem failed\n")
	}

	print2.Yellow("start download redis.conf\n")
	if filename, err := http.Download(config.Url.RedisConf); err != nil {
		print2.Red("download redis.conf failed\n")
	} else {
		confPath := config.Cfg.Service.Redis.ConfPath
		if _, err := os.Stat(confPath); os.IsNotExist(err) {
			if err = os.MkdirAll(confPath, os.ModePerm); err != nil {
				print2.Red("make redis conf path failed\n")
				os.Exit(1)
			}
		}
		if err := cmd.ExecNotExit("mv", "./"+filename, confPath); err != nil {
			print2.Red("move redis.conf to conf path failed\n")
		}
	}

	print2.Yellow("start download SQL file\n")
	if _, err := os.Stat(config.Cfg.Service.Mysql.ConfPath); os.IsNotExist(err) {
		if err = os.MkdirAll(config.Cfg.Service.Mysql.ConfPath+"/init", os.ModePerm); err != nil {
			print2.Red("make mysql conf path failed\n")
			os.Exit(1)
		}
	}
	if SQLFilename, err := http.Download(config.Url.SQLPath); err != nil {
		print2.Red("download SQL file failed\n")
	} else {
		cmd.Exec("mv", "./"+SQLFilename, config.Cfg.Service.Mysql.ConfPath)
	}

	docker.Up("./" + composeFilename)

	go func() {
		time.Sleep(time.Duration(10) * time.Second)
		cmd.ExecBashPipeNotExit("docker exec mysql mysql -uroot -p" + config.Cfg.Service.Mysql.Pwd + " -e\"source /etc/mysql/conf.d/metalabs-teaching.sql\"")
	}()
	print2.Blue("===================Service Installation Complete=====================\n")
}

func ReDeploy() {
	print2.Blue("=======================Start ReDeploy Service========================\n")
	composeFilePath := "./" + path.Base(config.Url.ServerComposeTemplate)
	if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
		print2.Yellow("start download and replace server-docker-compose.yml\n")
		deploy.DownloadAndReplaceTemplate(config.Url.ServerComposeTemplate)
	}
	docker.Down(composeFilePath)
	if err := docker.RmImage("mysql"); err != nil {
		print2.Red("delete docker image of mysql failed\n")
	}
	if err := docker.RmImage("redis"); err != nil {
		print2.Red("delete docker image of redis failed\n")
	}
	if err := docker.RmImage("minio"); err != nil {
		print2.Red("delete docker image of minio failed\n")
	}
	if err := docker.RmImage("srs"); err != nil {
		print2.Red("delete docker image of srs failed\n")
	}
	docker.Up(composeFilePath)
	print2.Blue("======================ReDeploy Service Complete======================\n")
}

func RemoveDeployment() {
	print2.Blue("========================Start Remove Service=========================\n")
	if err := docker.RmContainer("mysql"); err != nil {
		print2.Red("delete docker container of mysql failed\n")
	}
	if err := docker.RmContainer("redis"); err != nil {
		print2.Red("delete docker container of redis failed\n")
	}
	if err := docker.RmContainer("minio"); err != nil {
		print2.Red("delete docker container of minio failed\n")
	}
	if err := docker.RmContainer("srs"); err != nil {
		print2.Red("delete docker container of srs failed\n")
	}
	if err := docker.RmImage("mysql"); err != nil {
		print2.Red("delete docker image of mysql failed\n")
	}
	if err := docker.RmImage("redis"); err != nil {
		print2.Red("delete docker image of redis failed\n")
	}
	if err := docker.RmImage("minio"); err != nil {
		print2.Red("delete docker image of minio failed\n")
	}
	if err := docker.RmImage("srs"); err != nil {
		print2.Red("delete docker image of srs failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Mysql.DataPath); err != nil {
		print2.Red("delete mysql data path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Mysql.ConfPath); err != nil {
		print2.Red("delete mysql conf path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Redis.DataPath); err != nil {
		print2.Red("delete redis dataPath failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Redis.ConfPath); err != nil {
		print2.Red("delete redis conf path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Redis.LogPath); err != nil {
		print2.Red("delete redis log path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Minio.DataPath); err != nil {
		print2.Red("delete minio data path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Minio.ConfPath); err != nil {
		print2.Red("delete minio conf path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.Service.Srs.ConfPath); err != nil {
		print2.Red("delete srs conf path failed\n")
	}
	if err := cmd.ExecNotExit("docker", "network", "rm", "metalabs"); err != nil {
		print2.Red("delete docker network failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", "/usr/local/bin/docker-compose"); err != nil {
		print2.Red("delete docker-compose failed\n")
	}
	print2.Blue("=======================Remove Service Complete=======================\n")
}
