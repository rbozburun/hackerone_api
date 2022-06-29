package main

import (
	"io/ioutil"
	"os"
	"strings"
)

type H1Target struct {
	Id          string
	Handle_name string
	URL         string
}

// Return the content of given file
func readFile(fname string) string {
	databyte, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return string(databyte)
}

func main() {
	programs := getProgramsH1()
	var h1targets = []H1Target{}
	file, err := os.OpenFile("targets.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for i := 0; i < len(programs); i++ {
		program_handle := getProgramHandleH1(programs[i])

		for k := 0; k < len(program_handle.Relationships.StructedScopes.Data); k++ {
			scope := program_handle.Relationships.StructedScopes.Data[k]
			// If asset_type is URL
			if scope.ScopeAttrs.Asset_type == "URL" {
				h1target := H1Target{Id: scope.Id, Handle_name: program_handle.Attributes.Handle, URL: scope.ScopeAttrs.Asset_identifier}
				h1targets = append(h1targets, h1target)
			}
		}
	}

	for i := 0; i < len(h1targets); i++ {
		// Write URL and Handle_Name information if it is not already wrote
		if !strings.Contains(readFile("targets.txt"), "URL: "+h1targets[i].URL) {
			if _, err := file.WriteString("Handle name: " + h1targets[i].Handle_name + "\nURL: " + h1targets[i].URL + "\n\n"); err != nil {
				panic(err)
			}
		}
	}
}
