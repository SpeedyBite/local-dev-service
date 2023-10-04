package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("63")).Width(80)
	yesNoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func AskForConfirmation(s string) bool {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s %s: ", titleStyle.Render(s), yesNoStyle.Render("[y/n]"))

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
