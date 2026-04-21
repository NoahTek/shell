package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strings"
)

var history []string

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go func() {
		for {
			<-sigChan
			fmt.Println()
			customPrompt()
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		customPrompt()
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			continue
		}

		if err = processInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func customPrompt() {
	path, _ := os.Getwd()
	u, _ := user.Current()
	host, _ := os.Hostname()
	fmt.Printf("[%s@%s %s]> ", u.Username, host, path)
}

func processInput(input string) error {
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		return nil
	}

	history = append(history, input)

	return execInput(input)
}

func execInput(input string) error {
	commands := strings.Split(input, "|")
	var pipeIn io.ReadCloser
	var cmds []*exec.Cmd

	for i, cmdStr := range commands {
		args := strings.Fields(cmdStr)
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stderr = os.Stderr

		switch args[0] {
		case "cd":
			if len(args) < 2 {
				home, _ := os.UserHomeDir()
				return os.Chdir(home)
			}
			return os.Chdir(args[1])
		case "exit":
			os.Exit(0)
		case "history":
			var sb strings.Builder
			for j, cmdHist := range history {
				fmt.Fprintf(&sb, "%d: %s\n", j+1, cmdHist)
			}
			if i < len(commands)-1 {
				pipeIn = io.NopCloser(strings.NewReader(sb.String()))
			} else {
				fmt.Print(sb.String())
			}
			continue
		}

		if i > 0 {
			cmd.Stdin = pipeIn
		}

		if i < len(commands)-1 {
			var err error
			pipeIn, err = cmd.StdoutPipe()
			if err != nil {
				return err
			}
		} else {
			cmd.Stdout = os.Stdout
		}

		if err := cmd.Start(); err != nil {
			return err
		}

		cmds = append(cmds, cmd)
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}
