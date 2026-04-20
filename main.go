package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

var history []string

func execInput(input string) error {
	input = strings.TrimSuffix(input, "\n")
	if input == "" {
		return nil
	}

	history = append(history, input)

	args := strings.Split(input, " ")

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return os.Chdir("/home/manu/")
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "history":
		for i, cmd := range history {
			fmt.Printf("%d: %s\n", i+1, cmd)
		}
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		path, err := os.Getwd()
		checkErr(err)
		u, err := user.Current()
		checkErr(err)
		host, err := os.Hostname()
		checkErr(err)
		fmt.Printf("[%s@%s %s]> ", u.Username, host, path)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
