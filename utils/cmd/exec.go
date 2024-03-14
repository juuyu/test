package cmd

import (
	"errors"
	"fmt"
	"io"
	print2 "metalab-run/utils/print"
	"os"
	"os/exec"
)

// Exec 执行shell命令，失败就退出程序
func Exec(cmdName string, args ...string) {
	err := ExecNotExit(cmdName, args...)
	if err != nil {
		os.Exit(1)
	}
}

// ExecNotExit 执行shell命令，失败返回err
func ExecNotExit(cmdName string, args ...string) error {
	cmd := exec.Command(cmdName, args...)
	stdErrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		print2.Red(fmt.Sprintf("cmd execution: \"%s\" failed: %v\n", getCmdString(cmdName, args...), err))
		return err
	}
	defer cmd.Wait()
	errChan := make(chan bool)
	go syncLog(stdErrPipe, true, errChan)
	go syncLog(stdoutPipe, false, errChan)
	for i := 0; i < 2; i++ {
		if <-errChan {
			print2.Red(fmt.Sprintf("cmd execution: \"%s\" failed\n", getCmdString(cmdName, args...)))
			return errors.New("")
		}
	}
	return nil
}

// 通过管道同步获取日志的函数
func syncLog(reader io.ReadCloser, flag bool, errChain chan bool) {
	defer reader.Close()
	hasOut := false
	buf := make([]byte, 1024, 1024)
	for {
		strNum, err := reader.Read(buf)
		if strNum > 0 {
			hasOut = true
			outputByte := buf[:strNum]
			if flag {
				print2.Red(string(outputByte))
			} else {
				fmt.Print(string(outputByte))
			}
		}
		if err != nil {
			//读到结尾
			if err == io.EOF {
				err = nil
				if hasOut {
					fmt.Println()
					if flag {
						errChain <- true
					} else {
						errChain <- false
					}
				} else {
					errChain <- false
				}
				return
			}
		}
	}
}

func getCmdString(cmdName string, args ...string) string {
	res := cmdName
	for _, s := range args {
		res += " " + s
	}
	return res
}

// ExecBashPipe 执行shell管道命令（可以直接执行非管道shell命令），失败就退出程序
func ExecBashPipe(pipeCmd string) {
	if err := ExecBashPipeNotExit(pipeCmd); err != nil {
		os.Exit(1)
	}
}

// ExecBashPipeNotExit 执行shell管道命令（可以直接执行非管道shell命令），失败返回err
func ExecBashPipeNotExit(pipeCmd string) error {
	cmd := exec.Command("bash", "-c", pipeCmd)
	stdErrPipe, _ := cmd.StderrPipe()
	stdoutPipe, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		print2.Red(fmt.Sprintf("pipe cmd : \"%s\" execution failed: %v\n", pipeCmd, err))
		return err
	}
	defer cmd.Wait()
	errChan := make(chan bool)
	go syncLog(stdErrPipe, true, errChan)
	go syncLog(stdoutPipe, false, errChan)
	for i := 0; i < 2; i++ {
		if <-errChan {
			print2.Red(fmt.Sprintf("pipe cmd: \"%s\" execution failed\n", pipeCmd))
			return errors.New("")
		}
	}
	return nil
}

// ExecPipeCommands 使用两个cmd对象来执行管道命令，失败就退出程序
func ExecPipeCommands(commands ...*exec.Cmd) {
	if err := ExecPipeCommandsNotExit(commands...); err != nil {
		os.Exit(1)
	}
}

// ExecPipeCommandsNotExit 使用两个cmd对象来执行管道命令，失败返回err
func ExecPipeCommandsNotExit(commands ...*exec.Cmd) error {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		if err != nil {
			print2.Red(fmt.Sprintf("pipe cmd execution failed: %v\n", err))
			return err
		}
		command.Start()
		commands[i+1].Stdin = out
	}
	final, err := commands[len(commands)-1].Output()
	if err != nil {
		print2.Red(fmt.Sprintf("pipe cmd execution failed: %v\n", err))
		return err
	}
	fmt.Println(string(final))
	return nil
}
