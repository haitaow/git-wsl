package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type IgnoreError struct {
	GitCommand    string
	ExitCode      int
	EmptyStdErr   bool
	StdErrPhrases []string // ignored when EmptyStdErr is true
}

// UnmarshalJSON always sets default values of IgnoreError, even when un-marshaled as list element in Config.
// Without this the fields of Config.IgnoreErrors unspecified in json will take the values of the corresponding
// element in the default Config.IgnoreErrors value which can be unexpected.
func (t *IgnoreError) UnmarshalJSON(data []byte) error {
	type alias IgnoreError // prevent recursive calls to UnmarshalJSON
	v := alias{
		GitCommand:    "",
		ExitCode:      0,
		EmptyStdErr:   false,
		StdErrPhrases: []string{},
	}

	err := json.Unmarshal(data, &v)
	*t = IgnoreError(v)
	return err
}

type Config struct {
	LogAll       bool
	IgnoreErrors []IgnoreError
}

func getConfig(configFile string) Config {
	config := Config{
		LogAll: false,
		IgnoreErrors: []IgnoreError{
			{
				GitCommand:  "config",
				ExitCode:    1,
				EmptyStdErr: true,
			},
			{
				GitCommand:    "lfs",
				ExitCode:      1,
				StdErrPhrases: []string{"is not a git command", "lfs"},
			},
			{
				GitCommand:    "flow",
				ExitCode:      1,
				StdErrPhrases: []string{"is not a git command", "flow"},
			},
			{
				GitCommand:  "",
				ExitCode:    11,
				EmptyStdErr: true,
			},
		},
	}

	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			// try to write default config to file
			if b, err2 := json.Marshal(config); err2 == nil {
				_ = ioutil.WriteFile(configFile, b, 0600)
			} else {
				// print error and continue with default config
				fmt.Printf("Error writing default config file %v: ", configFile)
				fmt.Println(err)
			}
		} else {
			// print error and continue with default config
			fmt.Printf("Error reading config file %v: ", configFile)
			fmt.Println(err)
		}
		configData = nil
	}

	if configData != nil {
		if err = json.Unmarshal(configData, &config); err != nil {
			// print error and continue with default config
			fmt.Printf("Error parsing config file %v: ", configFile)
			fmt.Println(err)
		}
	}
	return config
}

func anyString(ss []string, f func(string) bool) bool {
	for _, s := range ss {
		if f(s) {
			return true
		}
	}
	return false
}

func allOf(ss []string, f func(string) bool) bool {
	for _, s := range ss {
		if !f(s) {
			return false
		}
	}
	return true
}

func main() {
	cmdInfo := "\n"
	cmdInfo += fmt.Sprintln("Args: ", os.Args)
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmdInfo += fmt.Sprintln("Working dir: ", workingDir)
	executable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	//cmdInfo += fmt.Sprintln("Executable: ", executable)

	logFile := filepath.Join(filepath.Dir(executable), "wsl-git.log")
	errFile := filepath.Join(filepath.Dir(executable), "wsl-git.err")
	configFile := filepath.Join(filepath.Dir(executable), "wsl-git.json")

	// gitArgs are translated args that will be sent to git in wsl
	gitArgs := append([]string{"git"}, os.Args[1:]...)
	for i, arg := range gitArgs {
		if len(arg) >= 3 && arg[1] == ':' && arg[2] == '\\' {
			gitArgs[i] = fmt.Sprintf("/mnt/%s/%s", strings.ToLower(arg[:1]), strings.ReplaceAll(arg[3:], "\\", "/"))
		}

		const wslPathPrefix = `\\wsl$\Ubuntu`
		if strings.HasPrefix(arg, wslPathPrefix) {
			gitArgs[i] = strings.ReplaceAll(arg[len(wslPathPrefix):], `\`, "/")
		}

		// TODO: translate other paths in args if needed.
	}
	//cmdInfo += fmt.Sprintf("gitArgs=%v\n", gitArgs)

	cmd := exec.Command("wsl", gitArgs...)
	cmdInfo += cmd.String() + "\n"

	var log *os.File
	var config Config
	defer func() {
		if err != nil {
			reportError(err, errFile, gitArgs, cmdInfo, log, config)
		}
	}()

	config = getConfig(configFile)
	if config.LogAll {
		if log, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600); err != nil {
			return
		}
		defer func() {
			if err := log.Close(); err != nil {
				panic(err)
			}
		}()
	}

	logf := func(format string, args ...interface{}) error {
		if log != nil {
			_, err := log.WriteString(fmt.Sprintf(format, args...))
			return err
		}
		return nil
	}

	if err = logf(cmdInfo); err != nil {
		return
	}

	gitOutput, err := cmd.Output()
	if _, err2 := os.Stdout.Write(gitOutput); err2 != nil {
		if err == nil {
			err = err2
		}
	}
	if err2 := logf(string(gitOutput)); err2 != nil {
		if err == nil {
			err = err2
		}
	}
}

func reportError(err error, errFile string, gitArgs []string, cmdInfo string, log *os.File, config Config) {
	if err == nil {
		return
	}

	knownErrors := []func(exitError *exec.ExitError, gitArgs []string) bool{
		func(exitError *exec.ExitError, gitArgs []string) bool {
			if exitError == nil {
				return false
			}
			for _, ignoreError := range config.IgnoreErrors {
				gitCommandMatch := true
				if ignoreError.GitCommand != "" {
					gitCommandMatch = anyString(gitArgs, func(s string) bool { return s == ignoreError.GitCommand })
				}

				exitCodeMatch := ignoreError.ExitCode == exitError.ExitCode()

				stdErr := string(exitError.Stderr)
				stdErrMatch :=
					allOf(ignoreError.StdErrPhrases, func(s string) bool { return strings.Contains(stdErr, s) })
				if ignoreError.EmptyStdErr {
					stdErrMatch = len(stdErr) == 0
				}

				if exitCodeMatch && stdErrMatch && gitCommandMatch {
					return true
				}
			}
			return false
		},
		// NOTE: hard coded known errors can be added here.
	}

	isKnownError := false
	exitCode := 101
	errorMessage := err.Error()
	if exitError, ok := err.(*exec.ExitError); ok {
		errorMessage += "\nstderr:\n"
		exitCode = exitError.ExitCode()
		errorMessage += string(exitError.Stderr)

		for _, f := range knownErrors {
			if f(exitError, gitArgs) {
				isKnownError = true
				break
			}
		}
	}
	errorMessage += "\n"

	// best effort reporting the error to caller.
	_, _ = os.Stderr.WriteString(errorMessage)
	if log != nil {
		_, _ = log.WriteString(errorMessage)
	}

	if isKnownError {
		fmt.Println("Is known error.  Skip logging.")
	} else {
		fmt.Printf("Logging error to %v.\n", errFile)
		f2, err2 := os.OpenFile(errFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err2 != nil {
			panic(err2)
		}
		defer func() {
			if err := f2.Close(); err != nil {
				panic(err)
			}
		}()
		_, _ = f2.WriteString(cmdInfo)
		_, _ = f2.WriteString(errorMessage)
	}

	fmt.Printf("Exit with code %v.\n", exitCode)
	os.Exit(exitCode)
}
