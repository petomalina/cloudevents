package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	project   string
	gitOrigin string
)

func walkFunc(data map[string]string) func(filePath string, info os.FileInfo, err error) error {
	return func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		split := strings.Split(filePath, "/")
		hidden := false
		for _, s := range split {
			if s[0] == '.' {
				hidden = true
				break
			}
		}
		if hidden {
			return nil
		}
		fmt.Printf("%s", filePath)
		t := template.New(path.Base(filePath))
		t, err = t.ParseFiles(filePath)
		if err != nil {
			return err
		}
		file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		err = t.Execute(file, data)
		if err != nil {
			return err
		}
		fmt.Println(" \u2714") // Check mark
		return nil
	}

}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
func main() {
	flag.StringVar(&project, "project", "forky", "Name for project")
	flag.StringVar(&gitOrigin, "git-origin", "", "remote git origin")

	flag.Parse()
	t := time.Now()
	data := map[string]string{
		"project": project,
	}
	err := filepath.Walk(".", walkFunc(data))
	if err != nil {
		panic(err)
	}
	if err := os.RemoveAll(".git"); err != nil {
		panic(err)
	}
	if err := runCmd("git", "init"); err != nil {
		panic(err)
	}
	if err := runCmd("git", "add", "-A"); err != nil {
		panic(err)
	}
	if err := runCmd("git", "commit", "-am", "Init commit from Forky", "--author", "Forky <forky@flowup.cz>"); err != nil {
		panic(err)
	}

	if gitOrigin != "" {
		if err := runCmd("git", "remote", "add", "origin", gitOrigin); err != nil {
			panic(err)
		}
		if err := runCmd("git", "push", "-u", "origin", "master"); err != nil {
			panic(err)
		}
	}

	if err := os.Remove("gen.go"); err != nil {
		panic(err)
	}
	fmt.Println("Done in ", time.Since(t))
}
