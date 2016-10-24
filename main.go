package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	stringNames     []string
	defStringFile   string
	outResPath      string
	configFile      string
	translationFile string
)

// StringInfo contains information for translation use.
type StringInfo struct {
	IsAddResource   bool
	DefaultValue    string
	TranslatedValue string
}

// trans_android_str -i <default string file path> -o <output res path> -c <config file> -t <translation file>
func main() {
	var err error
	var buf []byte
	patternString := `<string name="(?P<name>.*)">(?P<value>.*)</string>`
	patternAddString := `<add-resource type="string" name="(?P<name>.*)"`
	sep := "\n\n"
	var stringInfoMap = make(map[string]*StringInfo) // key: string name
	var stringXMLFileName = ""

	flag.StringVar(&defStringFile, "i", "", "default string file(english). Ex: ./res/values/strings.xml")
	flag.StringVar(&outResPath, "o", "", "output resource path. Ex: ./res")
	flag.StringVar(&configFile, "c", "config.json", "config JSON file contains names of strings need to be translated.")
	flag.StringVar(&translationFile, "t", "translation.txt", "translation file contains translated strings, iso 639-1 language name and  iso 3166-1 locale name.")

	flag.Parse()

	fmt.Println("defStringFile: " + defStringFile)
	fmt.Println("outResPath: " + outResPath)
	fmt.Println("configFile: " + configFile)
	fmt.Println("translationFile: " + translationFile)

	if buf, err = ioutil.ReadFile(configFile); err != nil {
		fmt.Println("Read config file err:")
		fmt.Println(err)
		return
	}

	if err = json.Unmarshal(buf, &stringNames); err != nil {
		fmt.Println("Parse string names err:")
		fmt.Println(err)
		return
	}

	fmt.Println("String Names:")
	for _, v := range stringNames {
		fmt.Println(v)
	}

	stringXMLFileName = path.Base(defStringFile)
	fmt.Printf("xml file name: %s\n", stringXMLFileName)

	if buf, err = ioutil.ReadFile(defStringFile); err != nil {
		fmt.Println("Read string file err:")
		fmt.Println(err)
		return
	}

	re := regexp.MustCompile(patternString)
	allStrings := re.FindAllStringSubmatch(string(buf), -1)
	for _, v := range allStrings {
		info := StringInfo{false, v[2], ""}
		stringInfoMap[v[1]] = &info
	}

	re = regexp.MustCompile(patternAddString)
	addStrings := re.FindAllStringSubmatch(string(buf), -1)
	for _, v := range addStrings {
		if _, ok := stringInfoMap[v[1]]; ok {
			stringInfoMap[v[1]].IsAddResource = true
		}
	}

	if buf, err = ioutil.ReadFile(translationFile); err != nil {
		fmt.Println("Read translation file err:")
		fmt.Println(err)
		return
	}

	transStringCollection := strings.Split(string(buf), sep)
	for _, v := range transStringCollection {
		//fmt.Println(v)
		//fmt.Println()

		transStrings := strings.Split(v, "\n")
		if len(transStrings) != len(stringNames)+1 {
			fmt.Println("Translated string count does not match the count of string name in config.json.")
			return
		}

		language := transStrings[0]
		fmt.Printf("Language: %s\n", language) // 1st string is iso 639-1 language name and iso 3166-1 locale name
		// mkdir for res/values/xx-xx/ and write string file
		langPath := outResPath + "/values-" + language
		if err = os.MkdirAll(langPath, os.ModePerm); err != nil {
			fmt.Printf("Create out language dir err:%s\n", err)
			return
		}

		for i := 1; i < len(transStrings); i++ {
			//fmt.Println(transStrings[i])
			stringName := stringNames[i-1]
			if _, ok := stringInfoMap[stringName]; !ok {
				fmt.Printf("string name: %s can not found in default string xml", stringName)
				return
			}
			stringInfoMap[stringName].TranslatedValue = transStrings[i]
		}

		body := `<?xml version="1.0" encoding="utf-8"?>` + "\n" + "<resources>\n"
		for k, v := range stringInfoMap {
			//fmt.Printf("%s:%v\n", k, v)
			s := ""
			if len(v.TranslatedValue) == 0 { // no need to translated, just use default string
				continue
			}

			if !v.IsAddResource {
				s = fmt.Sprintf("  <string name=\"%s\">%s</string>\n", k, v.TranslatedValue)
				body += s
			} else {
				s = fmt.Sprintf("  <add-resource type=\"string\" name=\"%s\" />\n  <string name=\"%s\">%s</string>\n", k, k, v.TranslatedValue)
				body += s
			}
		}
		//fmt.Println()
		body += "</resources>\n"
		fmt.Println(body)

		file := langPath + "/" + stringXMLFileName
		if err = ioutil.WriteFile(file, []byte(body), os.ModePerm); err != nil {
			fmt.Printf("Write string file err:%s\n", err)
			return
		}
	}
}
