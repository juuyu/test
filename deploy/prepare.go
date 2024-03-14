package deploy

import (
	"metalab-run/config"
	"metalab-run/utils/cmd"
	"metalab-run/utils/http"
	print2 "metalab-run/utils/print"
	"os"
	"os/exec"
)

func PreparationWork() {
	print2.Green("==========Start System Check to Prepare Deploy Environment===========\n")
	if _, err := exec.LookPath("docker"); err != nil {
		print2.Yellow("docker is not installed, start install docker\n")
		cmd.Exec("yum", "-y", "install", "docker")
	}
	if err := cmd.ExecNotExit("docker", "ps"); err != nil {
		cmd.Exec("systemctl", "start", "docker")
		cmd.Exec("systemctl", "enable", "docker")
	}
	if err := cmd.ExecBashPipeNotExit("docker login -u " + config.Cfg.Client.ClientCode + " -p " + "'" + config.Cfg.Client.DockerPwd + "' " + config.DockerRegistry); err != nil {
		cmd.ExecBashPipe("echo   '{\n  \"insecure-registries\": [\n    \"60.188.46.2:20012\"\n  ]\n}' > /etc/docker/daemon.json")
		cmd.Exec("systemctl", "restart", "docker")
		cmd.ExecBashPipeNotExit("docker login -u " + config.Cfg.Client.ClientCode + " -p " + "'" + config.Cfg.Client.DockerPwd + "' " + config.DockerRegistry)
	}
	cmd.ExecNotExit("docker", "network", "rm", "metalabs")
	cmd.ExecNotExit("docker", "network", "create", "metalabs")
	if _, err := exec.LookPath("docker-compose"); err != nil {
		print2.Yellow("docker-compose is not installed, start install docker-compose\n")
		filename, err := http.Download(config.Url.DockerCompose)
		if err != nil {
			print2.Red("docker-compose file download failed\n")
			os.Exit(1)
		}
		cmd.Exec("mv", "./"+filename, "/usr/local/bin/docker-compose")
		cmd.Exec("chmod", "777", "/usr/local/bin/docker-compose")
	}
	print2.Green("=======================System Check Finished=========================\n")
}
