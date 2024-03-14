package print

import "fmt"

const (
	textBlack = iota + 30
	textRed
	textGreen
	textYellow
	textBlue
	textPurple
	textCyan
	textWhite
)

func Black(str string) {
	textColor(textBlack, str)
}

func Red(str string) {
	textColor(textRed, str)
}
func Yellow(str string) {
	textColor(textYellow, str)
}
func Green(str string) {
	textColor(textGreen, str)
}
func Cyan(str string) {
	textColor(textCyan, str)
}
func Blue(str string) {
	textColor(textBlue, str)
}
func Purple(str string) {
	textColor(textPurple, str)
}
func White(str string) {
	textColor(textWhite, str)
}

func textColor(color int, str string) {
	fmt.Printf("\x1b[0;%dm%s\x1b[0m", color, str)
}
