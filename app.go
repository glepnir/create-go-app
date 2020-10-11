// Copyright 2020 glepnir. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"
)

var (
	layoutYaml string
)

func main() {
	flag.StringVar(&layoutYaml, "l", "", "a project layout yaml file path")
	flag.Parse()
	projectLayout := make(map[string]interface{})
	parseLayout(projectLayout)
	createProject(projectLayout)
}

func fullpath(dir string, approot ...string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if len(approot) > 0 {
		return currentDir + osSparator() + approot[0] + osSparator() + dir
	}
	return currentDir + osSparator() + dir
}

func parseLayout(projectLayout map[string]interface{}) {
	contents, err := ioutil.ReadFile(layoutYaml)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(contents, projectLayout)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(projectLayout)
}

func osSparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func isFile(file interface{}) bool {
	f, ok := file.(string)
	if !ok {
		log.Fatal("not a string value")
		return false
	}
	ext := filepath.Ext(f)
	if len(ext) == 0 {
		return false
	}
	return true
}

func createProject(projectLayout map[string]interface{}) {
	appName, ok := projectLayout["app"].(string)
	if !ok {
		log.Fatal("layout yaml file format wrong miss filed app")
		return
	}
	err := os.Mkdir(fullpath(appName), 0777)
	if err != nil {
		log.Fatal(err)
	}
	layout, ok := projectLayout["layout"].([]interface{})

	for _, v := range layout {
		if f, ok := v.(map[interface{}]interface{}); ok {
			for a, _ := range f {
				if a, ok := a.(string); ok {
					if !isFile(a) {
						os.Mkdir(fullpath(a, appName), 0777)
					}
				}
			}
		}
		if f, ok := v.(string); ok {
			os.Mkdir(fullpath(f, appName), 0777)
		}
	}
}
