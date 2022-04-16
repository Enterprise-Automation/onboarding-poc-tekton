package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type manifest struct {
	Owners []string `yaml:"owners"`
}

var filePath string

func init() {
	flag.StringVar(&filePath, "filepath", "", "path to manifest")
	flag.Parse()
}

func (c *manifest) getManifest() *manifest {

	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func main() {

	if filePath == "" {
		log.Fatalf("No filepath provided. please use --filepath <path_to_file>")
		os.Exit(1)
	}

	appID := substr(filePath, 0, strings.Index(filePath, "/"))

	var m manifest
	m.getManifest()

	file, err := os.OpenFile("CODEOWNERS", os.O_RDWR, os.ModeAppend)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string

	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	found := false
	for idx, eachline := range txtlines {
		if strings.Contains(eachline, appID) {
			found = true
			owners := ""
			for _, user := range m.Owners {
				owners += fmt.Sprintf("@%s ", user)
			}
			txtlines[idx] = fmt.Sprintf("%s/* %s", appID, owners)
			fmt.Printf("Updated entry for %s in CODEOWNERS.\n", appID)
		}
	}

	if !found {
		owners := ""
		for _, user := range m.Owners {
			owners += fmt.Sprintf("@%s ", user)
		}
		txtlines = append(txtlines, fmt.Sprintf("%s/* %s", appID, owners))
		fmt.Printf("Added new entry for %s in CODEOWNERS.\n", appID)
	}

	err = file.Truncate(0)
	_, err = file.Seek(0, 0)
	if _, err := file.WriteString(strings.Join(txtlines[:], "\n")); err != nil {
		fmt.Println(err)
		file.Close()
		return
	}

	file.Close()

}

func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}
