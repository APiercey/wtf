package client

import (
	"os/exec"
	"fmt"
	"strings"
)

func ProjectList() []string {
	stdout, err := exec.Command("tmuxinator", "list").Output()

	if err != nil {
			fmt.Println(err.Error())
			return []string{}
	}

	output := strings.Split(string(stdout), ":")
	projectList := strings.Split(output[1], " ")
	
	var cleanProjectList []string

	for _, str := range projectList {
		cleanString := strings.ReplaceAll(str, "\n", "")

		if cleanString != "" {
			cleanProjectList = append(cleanProjectList, cleanString)
		}
	}

	return cleanProjectList
}

func StartProject(projectName string) {
	_, err := exec.Command("tmuxinator", "start", projectName).Output()

	if err != nil {
			fmt.Println(err.Error())
	}
}

func EditProject(projectName string) {
	subcommand := fmt.Sprintf("tmuxinator edit %s", projectName)
	_, err := exec.Command("tmux", "new-window", subcommand).Output()

	if err != nil {
			fmt.Println(err.Error())
	}
}

func DeleteProject(projectName string) {
	_, err := exec.Command("tmuxinator", "delete", projectName).Output()

	if err != nil {
			fmt.Println(err.Error())
	}
}

func CopyProject(leftProj, rightProj string) {
	subcommand := fmt.Sprintf("tmuxinator copy %s %s", leftProj, rightProj)
	_, err := exec.Command("tmux", "new-window", subcommand).Output()

	if err != nil {
			fmt.Println(err.Error())
	}
}
