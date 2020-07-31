// Copyright 2020 glepnir. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

var (
	appname  string
	platform string
	repo     string
)

func main() {
	flag.StringVar(&appname, "a", "app", "app name")
	flag.StringVar(&platform, "p", "github", "platform name")
	flag.StringVar(&repo, "r", "glepnir", "repo name")
	flag.Parse()
	task()
}

func currentpath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func task() {
	options := platform + ".com/" + repo
	app := currentpath() + string(os.PathSeparator) + appname
	if err := os.MkdirAll(app, 0755); err != nil {
		log.Fatal("create folder failed")
	}

	err := os.Chdir(app)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("go", "mod", "init", options)
	cmd.Dir = app
	if err := cmd.Run(); err != nil {
		log.Fatal("failed go mod init")
	}
	fmt.Println("Success")
}
