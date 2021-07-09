package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	cmdInfo := "\n"
	cmdInfo += fmt.Sprintln("Args: ", os.Args)
	workingDir, err := os.Getwd()
	cmdInfo += fmt.Sprintln("Working dir: ", workingDir, " err: ", err)
	executable, err := os.Executable()
	cmdInfo += fmt.Sprintln("Executable: ", executable, " err: ", err)
	logFile := filepath.Join(filepath.Dir(executable), "wsl-git.log")
	cmdInfo += fmt.Sprintln("Log file: ", logFile)

	// TODO: use config file to control logging all/errors/none.
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err = f.WriteString(cmdInfo); err != nil {
		panic(err)
	}

	gitArgs := os.Args[1:]
	gitArgs = append([]string{"git"}, gitArgs...)
	for i, arg := range gitArgs {
		if len(arg) >= 3 && arg[1] == ':' && arg[2] == '\\' {
			gitArgs[i] = fmt.Sprintf("/mnt/%s/%s", strings.ToLower(arg[:1]), strings.ReplaceAll(arg[3:], "\\", "/"))
		}
		// TODO: translate other paths in args if needed.
	}

	cmd := exec.Command("wsl", gitArgs...)
	if _, err = f.WriteString(cmd.String() + "\n"); err != nil {
		panic(err)
	}

	gitOutput, err := cmd.Output()
	if err == nil {
		_, err = os.Stdout.Write(gitOutput)
	}
	if err == nil {
		_, err = f.Write(gitOutput)
	}

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// best effort reporting the error to caller.
			if _, err2 := os.Stderr.Write(exitError.Stderr); err2 != nil {
				panic(err2)
			}
			if _, err2 := f.Write(exitError.Stderr); err2 != nil {
				panic(err2)
			}
		} else {
			panic(err)
		}
	}
}
