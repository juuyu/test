package docker

import "metalab-run/utils/cmd"

func Up(composeFilePath string) {
	cmd.ExecNotExit("docker-compose", "-f", composeFilePath, "up", "-d")
}

func Down(composeFilePath string) {
	cmd.ExecNotExit("docker-compose", "-f", composeFilePath, "down")
}

func RmContainer(containerName string) error {
	if err := cmd.ExecBashPipeNotExit("docker rm -f $(docker ps -a | grep " + containerName + " | awk '{print $1}')"); err != nil {
		return err
	}
	return nil
}

func RmImage(imageName string) error {
	if err := cmd.ExecBashPipeNotExit("docker rmi -f $(docker images | grep " + imageName + " | awk '{print $3}')"); err != nil {
		return err
	}
	return nil
}
