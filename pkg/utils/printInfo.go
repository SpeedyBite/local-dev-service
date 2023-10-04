package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Width(80)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("27")).Width(80)
)

func PrintError(s string) {
	str := fmt.Sprintf("%s %s", "[Error]", s)
	fmt.Println(errorStyle.Render(str))
}

func PrintInfo(s string) {
	str := fmt.Sprintf("%s %s", "[Info]", s)
	fmt.Println(infoStyle.Render(str))
}
