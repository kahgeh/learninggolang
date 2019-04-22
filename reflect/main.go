package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
)

type Lang struct {
	Name string
	Year int
	Url  string
}

func main() {
	myPrint("Hello", 42, "\n")
	lang := Lang{"Go", 2009, "http://golang.org/"}
	fmt.Printf("%+v\n", lang)
	fmt.Printf("%#v\n", lang)
	data, err := json.Marshal(lang)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)

	data, err = xml.MarshalIndent(lang, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", data)
}

func myPrint(args ...interface{}) {
	for _, arg := range args {
		switch v := reflect.ValueOf(arg); v.Kind() {
		case reflect.String:
			os.Stdout.WriteString(v.String())
		case reflect.Int:
			os.Stdout.WriteString(strconv.FormatInt(v.Int(), 10))
		}
	}
}
