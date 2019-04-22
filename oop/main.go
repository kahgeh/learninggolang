package main

import "fmt"

type LangType string

const (
	imperative  LangType = "Imperative"
	declarative LangType = "Declarative"
)

type Lang struct {
	Name string
	Year int
	Url  string
	Type LangType
}

func (lang *Lang) generateLanguagesReport() {
	fmt.Printf("%s  %d %s\n", lang.Name, lang.Year, lang.Url)
}

func main() {
	languages := []Lang{
		Lang{"Go", 2009, "http://golang.org/", imperative},
		Lang{"Python", 1994, "https://www.python.org/", imperative},
	}

	languages[0].generateLanguagesReport()
}
