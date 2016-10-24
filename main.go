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
	string_names     []string
	def_string_file  string
	out_res_path     string
	config_file      string
	translation_file string
)

type StringInfo struct {
	IsAddResource   bool
	DefaultValue    string
	TranslatedValue string
}

// trans_android_str -i <default string file path> -o <output res path> -c <config file> -t <translation file>
func main() {
	var err error
	var buf []byte
	pattern_string := `<string name="(?P<name>.*)">(?P<value>.*)</string>`
	pattern_add_string := `<add-resource type="string" name="(?P<name>.*)"`
	sep := "\n\n"
	var stringInfoMap = make(map[string]*StringInfo) // key: string name
	var string_xml_file_name = ""

	flag.StringVar(&def_string_file, "i", "", "default string file(english). Ex: ./res/values/strings.xml")
	flag.StringVar(&out_res_path, "o", "", "output resource path. Ex: ./res")
	flag.StringVar(&config_file, "c", "config.json", "config JSON file contains names of strings need to be translated.")
	flag.StringVar(&translation_file, "t", "translation.txt", "translation file contains translated strings, iso 639-1 language name and  iso 3166-1 locale name.")

	flag.Parse()

	fmt.Println("def_string_file: " + def_string_file)
	fmt.Println("out_res_path: " + out_res_path)
	fmt.Println("config_file: " + config_file)
	fmt.Println("translation_file: " + translation_file)

	if buf, err = ioutil.ReadFile(config_file); err != nil {
		fmt.Println("Read config file err:")
		fmt.Println(err)
		return
	}

	if err = json.Unmarshal(buf, &string_names); err != nil {
		fmt.Println("Parse string names err:")
		fmt.Println(err)
		return
	}

	fmt.Println("String Names:")
	for _, v := range string_names {
		fmt.Println(v)
	}

	string_xml_file_name = path.Base(def_string_file)
	fmt.Printf("xml file name: %s\n", string_xml_file_name)

	if buf, err = ioutil.ReadFile(def_string_file); err != nil {
		fmt.Println("Read string file err:")
		fmt.Println(err)
		return
	}

	re := regexp.MustCompile(pattern_string)
	all_strings := re.FindAllStringSubmatch(string(buf), -1)
	for _, v := range all_strings {
		info := StringInfo{false, v[2], ""}
		stringInfoMap[v[1]] = &info
	}

	re = regexp.MustCompile(pattern_add_string)
	add_strings := re.FindAllStringSubmatch(string(buf), -1)
	for _, v := range add_strings {
		if _, ok := stringInfoMap[v[1]]; ok {
			stringInfoMap[v[1]].IsAddResource = true
		}
	}

	if buf, err = ioutil.ReadFile(translation_file); err != nil {
		fmt.Println("Read translation file err:")
		fmt.Println(err)
		return
	}

	trans_string_collection := strings.Split(string(buf), sep)
	for _, v := range trans_string_collection {
		//fmt.Println(v)
		//fmt.Println()

		trans_strings := strings.Split(v, "\n")
		if len(trans_strings) != len(string_names)+1 {
			fmt.Println("Translated string count does not match the count of string name in config.json.")
			return
		}

		language := trans_strings[0]
		fmt.Printf("Language: %s\n", language) // 1st string is iso 639-1 language name and iso 3166-1 locale name
		// mkdir for res/values/xx-xx/ and write string file
		lang_path := out_res_path + "/values-" + language
		if err = os.MkdirAll(lang_path, os.ModePerm); err != nil {
			fmt.Printf("Create out language dir err:%s\n", err)
			return
		}

		for i := 1; i < len(trans_strings); i++ {
			//fmt.Println(trans_strings[i])
			string_name := string_names[i-1]
			if _, ok := stringInfoMap[string_name]; !ok {
				fmt.Printf("string name: %s can not found in default string xml", string_name)
				return
			}
			stringInfoMap[string_name].TranslatedValue = trans_strings[i]
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

		file := lang_path + "/" + string_xml_file_name
		if err = ioutil.WriteFile(file, []byte(body), os.ModePerm); err != nil {
			fmt.Printf("Write string file err:%s\n", err)
			return
		}
	}
}
