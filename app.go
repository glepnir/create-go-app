// Copyright 2020 glepnir. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	layoutYaml string
)

func main() {
	flag.StringVar(&layoutYaml, "l", "", "a project layout yaml file path")
	flag.Parse()
	appLayout := make(map[string]interface{})
	parseLayout(appLayout)
	createApp(appLayout)
}

func checkErr(err error) {
	if err != nil {
		printError(err.Error())
		return
	}
}

func fullpath(dir string, approot ...string) string {
	var sb strings.Builder
	currentDir, err := os.Getwd()
	checkErr(err)
	sb.WriteString(currentDir)
	sb.WriteString(osSparator())
	if len(approot) > 0 {
		sb.WriteString(approot[0])
		sb.WriteString(osSparator())
	}
	sb.WriteString(dir)
	return sb.String()
}

func parseLayout(appLayout map[string]interface{}) {
	contents, err := ioutil.ReadFile(layoutYaml)
	checkErr(err)
	err = yaml.Unmarshal(contents, appLayout)
	if err != nil {
		log.Fatal(err)
	}
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
		printWarn("Something wron in isfile function")
		return false
	}
	ext := filepath.Ext(f)
	if len(ext) == 0 {
		return false
	}
	return true
}

func createApp(appLayout map[string]interface{}) {
	appName, ok := appLayout["app"].(string)
	if !ok {
		printError("Failed create app missed app keyword")
		return
	}
	err := os.Mkdir(fullpath(appName), 0777)
	checkErr(err)
	layout, ok := appLayout["layout"].([]interface{})

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
