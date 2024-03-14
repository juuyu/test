package deploy

import (
	"bufio"
	"fmt"
	"metalab-run/config"
	"metalab-run/utils/http"
	print2 "metalab-run/utils/print"
	"os"
	"strings"
)

func DownloadAndReplaceTemplate(url string) {
	filename, err := http.Download(url)
	if err != nil {
		print2.Red(filename + "download failed\n")
		os.Exit(1)
	}
	err = Replace(filename)
	if err != nil {
		print2.Red("replace template file failed: " + filename)
		os.Exit(1)
	}
}

func Replace(filename string) error {
	res := ""
	err := func() error {
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			res += processTemplate(line)
		}
		if err := scanner.Err(); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		print2.Red("open/read " + filename + " failed\n")
		return err
	}
	fw, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer fw.Close()
	w := bufio.NewWriter(fw)
	fmt.Fprint(w, res)
	w.Flush()
	return nil
}

func processTemplate(s string) string {
	i := strings.Index(s, "@")
	if i == -1 {
		return s + "\n"
	}
	if s[i+1] == '{' {
		j := i + 1
		for j < len(s) {
			if s[j] == '}' {
				break
			}
			j++
		}
		if j < len(s) && s[j] == '}' {
			t := s[i : j+1]
			replace, ok := config.TemplateMap[t]
			if ok {
				s = strings.Replace(s, t, replace, -1)
			}
		}
	}
	return s + "\n"
}
