// Package main
// @title
// @description
// @author njy
// @since 2022/12/12 16:25
package main

import (
	"fmt"
	"metalab-run/deploy"
	"metalab-run/deploy/app"
	"metalab-run/deploy/service"
	print2 "metalab-run/utils/print"
	"os"
)

// curl --output all_config.toml https://minio.metalabs.top:20080/conf/run/all_config.toml && curl --output gometalab  https://minio.metalabs.top:20080/conf/run/gometalab
func main() {
	print2.Green(" ____    ____   _________     _________        __")
	print2.Cyan("         _____            __         ______     \n")
	print2.Green("|_   \\  /   _| |_   ___  |   |  _   _  |       /  \\       ")
	print2.Cyan("|_   _|          /  \\       |_   _ \\   \n")
	print2.Green("  |   \\/   |     | |_  \\_|   |_/ | | \\_|      / /\\ \\        ")
	print2.Cyan("| |           / /\\ \\        | |_) |  \n")
	print2.Green("  | |\\  /| |     |  _|  _        | |         / ____ \\       ")
	print2.Cyan("| |   _      / ____ \\       |  __'.  \n")
	print2.Green(" _| |_\\/_| |_   _| |___/ |      _| |_      _/ /    \\ \\_    ")
	print2.Cyan("_| |__/ |   _/ /    \\ \\_    _| |__) | \n")
	print2.Green("|_____||_____| |_________|     |_____|    |____|  |____|  ")
	print2.Cyan("|________|  |____|  |____|  |_______/   \n")
	fmt.Println("=============Please contact us if you have any questions=============")
	if len(os.Args) > 1 && os.Args[1] == "postern" {
		printOptions()
	} else {
		deploy.PreparationWork()
		deployAllService(1)
	}
}

func printOptions() {
	for true {
		fmt.Println("1. Deploy/Remove All of the Services;")
		fmt.Println("2. Deploy/Remove Only Applications.")
		print2.Green("Please enter your options (numbers only): \n")
		var option int
		_, err := fmt.Scanf("%d\n", &option)
		fmt.Println(option)
		if err != nil {
			print2.Red("Input error, please re-input!\n")
			continue
		}
		if option < 1 || option > 2 {
			print2.Red("Input error, please re-input!\n")
			continue
		}
		fmt.Print("\033[H\033[2J")
		deployOption := printDeployOptions()
		deploy.PreparationWork()
		switch option {
		case 1:
			deployAllService(deployOption)
			goto End
		case 2:
			deployApp(deployOption)
			goto End
		}
	}
End:
}

func printDeployOptions() (option int) {
	for true {
		fmt.Println("1. First Deploy;")
		fmt.Println("2. ReDeploy But Retain Data;")
		fmt.Println("3. Remove Deployment.")
		print2.Green("Please enter your options (numbers only): \n")
		_, err := fmt.Scanf("%d\n", &option)
		if err != nil {
			print2.Red("Input error, please re-input!\n")
			continue
		}
		if option < 1 || option > 3 {
			print2.Red("Input error, please re-input!\n")
			continue
		}
		break
	}
	return option
}

func deployAllService(option int) {
	switch option {
	case 1:
		service.FirstDeploy()
		app.FirstDeploy()
	case 2:
		service.ReDeploy()
		app.ReDeploy()
	case 3:
		app.RemoveDeployment()
		service.RemoveDeployment()
	}
}

func deployApp(option int) {
	switch option {
	case 1:
		app.FirstDeploy()
	case 2:
		app.ReDeploy()
	case 3:
		app.RemoveDeployment()
	}
}
