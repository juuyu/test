package app

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
)

func FirstDeploy() {
	print2.Cyan("=========================Start Install App===========================\n")

	print2.Yellow("start download and replace application.yml\n")
	deploy.DownloadAndReplaceTemplate(config.Url.BackConfPath)
	ymlFilename := path.Base(config.Url.BackConfPath)
	ymlPath := config.Cfg.App.Metalab.ConfPath
	if _, err := os.Stat(ymlPath); os.IsNotExist(err) {
		err = os.MkdirAll(ymlPath, os.ModePerm)
		if err != nil {
			print2.Red("make application.yml conf path failed\n")
		}
	}
	if err := cmd.ExecNotExit("mv", "./"+ymlFilename, ymlPath); err != nil {
		print2.Red("move application.yml to conf path failed\n")
	}

	confPath := config.Cfg.App.MetalabPage.ConfPath
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		err = os.MkdirAll(confPath, os.ModePerm)
		if err != nil {
			print2.Red("make front conf path failed\n")
			os.Exit(1)
		}
	}

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

	print2.Yellow("start download front conf file\n")
	frontConfFilename := path.Base(config.Url.FrontConfPath)
	deploy.DownloadAndReplaceTemplate(config.Url.FrontConfPath)
	if err := cmd.ExecNotExit("mv", "./"+frontConfFilename, confPath); err != nil {
		print2.Red("move front conf file to conf path failed\n")
	}

	print2.Yellow("start download ffmpeg exec file\n")
	if _, err := os.Stat(config.Cfg.App.Metalab.DataPath); os.IsNotExist(err) {
		if err := cmd.ExecNotExit("mkdir", "-p", config.Cfg.App.Metalab.DataPath); err != nil {
			print2.Red("make ffmpeg path failed\n")
		}
	}
	if ffmpegFile, err := http.Download(config.Url.FfmpegPath); err != nil {
		print2.Red("download ffmpeg exec file failed\n")
	} else if err := cmd.ExecNotExit("mv", "./"+ffmpegFile, config.Cfg.App.Metalab.DataPath); err != nil {
		print2.Red("move ffmpeg exec file to data path failed\n")
	} else {
		cmd.ExecNotExit("chmod", "777", config.Cfg.App.Metalab.DataPath+"/"+ffmpegFile)
	}

	print2.Yellow("start download and replace ffmpeg conf file\n")
	deploy.DownloadAndReplaceTemplate(config.Url.FfmpegConfPath)
	if _, err := os.Stat(config.Cfg.App.Ffmpeg.ConfPath); os.IsNotExist(err) {
		err = os.MkdirAll(config.Cfg.App.Ffmpeg.ConfPath, os.ModePerm)
		if err != nil {
			print2.Red("make ffmpeg conf path failed\n")
			os.Exit(1)
		}
	}
	ffmpegConfFilename := path.Base(config.Url.FfmpegConfPath)
	if err := cmd.ExecNotExit("mv", "./"+ffmpegConfFilename, config.Cfg.App.Ffmpeg.ConfPath); err != nil {
		print2.Red("move ffmpeg conf file to conf path failed\n")
	}

	print2.Yellow("start download and replace app-docker-compose.yml\n")
	deploy.DownloadAndReplaceTemplate(config.Url.AppComposeTemplate)
	composeFilename := path.Base(config.Url.AppComposeTemplate)
	docker.Up("./" + composeFilename)
	print2.Cyan("=====================App Installation Complete=======================\n")
}

func ReDeploy() {
	print2.Cyan("=========================Start ReDeploy App==========================\n")
	composeFilePath := "./" + path.Base(config.Url.AppComposeTemplate)
	if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
		print2.Yellow("start download and replace app-docker-compose.yml\n")
		deploy.DownloadAndReplaceTemplate(config.Url.AppComposeTemplate)
	}
	docker.Down(composeFilePath)
	if err := docker.RmImage("metalabs"); err != nil {
		print2.Red("delete docker image of metalabs failed\n")
	}
	if err := docker.RmImage("metalabs-page"); err != nil {
		print2.Red("delete docker image of metalabs-page failed\n")
	}
	if err := docker.RmImage("ffmpeg-go"); err != nil {
		print2.Red("delete docker image of ffmpeg-go failed\n")
	}
	docker.Up(composeFilePath)
	print2.Cyan("========================ReDeploy App Complete========================\n")
}

func RemoveDeployment() {
	print2.Cyan("==========================Start Remove App===========================\n")
	if err := docker.RmContainer("metalabs"); err != nil {
		print2.Red("delete docker container of metalabs failed\n")
	}
	if err := docker.RmContainer("metalabs-page"); err != nil {
		print2.Red("delete docker container of metalabs-page failed\n")
	}
	if err := docker.RmContainer("ffmpeg-go"); err != nil {
		print2.Red("delete docker container of ffmpeg-go failed\n")
	}
	if err := docker.RmImage("metalabs"); err != nil {
		print2.Red("delete docker image of metalabs failed\n")
	}
	if err := docker.RmImage("metalabs-page"); err != nil {
		print2.Red("delete docker image of metalabs-page failed\n")
	}
	if err := docker.RmImage("ffmpeg-go"); err != nil {
		print2.Red("delete docker image of ffmpeg-go failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.App.Metalab.DataPath); err != nil {
		print2.Red("delete metalab data path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.App.Metalab.ConfPath); err != nil {
		print2.Red("delete metalab conf path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.App.MetalabPage.ConfPath); err != nil {
		print2.Red("delete metalab-page conf path failed\n")
	}
	if err := cmd.ExecNotExit("rm", "-rf", config.Cfg.App.Ffmpeg.ConfPath); err != nil {
		print2.Red("delete ffmpeg-go conf path failed\n")
	}
	print2.Cyan("=========================Remove App Complete=========================\n")
}
