package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"kevinkamps.nl/gitlab-ci/nomad/configuration"
	nomad "kevinkamps.nl/gitlab-ci/nomad/nomad"
	"log"
	"os"
	"strings"
)

var version string

func main() {
	var configurations []configuration.Configuration

	nomadConfiguration := nomad.NewNomadConfiguration()
	configurations = append(configurations, nomadConfiguration)

	jobHclFile := flag.String("job-hcl-file", `nomad-job.hcl`, "Path to the job hch file")
	showVersion := flag.Bool("version", false, "Prints the version of the application and exits")
	flag.Parse()
	for _, configuration := range configurations {
		configuration.Parse()
	}
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	nomadHelper := nomad.NewNomadHelper(nomadConfiguration)

	// parse hcl
	hcl, err := ioutil.ReadFile(*jobHclFile)
	if err != nil {
		log.Fatalln(err)
	}

	message := map[string]interface{}{
		"JobHCL":       replaceVars(hcl),
		"Canonicalize": true,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	jsonContent := nomadHelper.ParseHclJob(bytes.NewBuffer(bytesRepresentation))
	nomadHelper.UpdateJob(jsonContent)
}

//Replaces variables in the content. Currently it only replaces environment variables.
func replaceVars(content []byte) string {
	contentString := string(content)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		searchKey := "$" + pair[0]
		i := strings.Index(contentString, searchKey)
		for i >= 0 {
			contentString, i = replaceString(contentString, searchKey, pair[1], i)

			// check if there is a next match
			newIndex := strings.Index(contentString[i+1:], searchKey)
			if newIndex == -1 {
				i = newIndex
				continue
			}
			i = i + 1 + strings.Index(contentString[i+1:], searchKey)
		}
	}

	return contentString
}

//Replaces the old with the new in the content string. Only replaces the first find.
//Parameter content 	The contents of the string in which to replace the old with the new
//Parameter old			Will be replaced with the new
//Parameter new			Used to replace the old
//Parameter fromIndex	The index from which to start looking
//return string 		Content with the replaced value
//return int 			The value of the index of the new content replaced value | the original fromIndex
func replaceString(content string, old string, new string, fromIndex int) (string, int) {
	start := strings.Index(content[fromIndex:], old)
	if start == -1 {
		return content, fromIndex
	}

	var end int
	for end = start + 1; end < len(content[fromIndex:]) && isShellVariableNameCharacter(content[fromIndex:][end]); end++ {
	}
	if len(old) != end-start {
		return content, fromIndex
	}

	return content[0:fromIndex] + new + content[fromIndex+end:], fromIndex + end
}

//Checks if a certain character is valid for a shell variable name.
//Parameter hclContent	The content of a hcl formatted file
func isShellVariableNameCharacter(c uint8) bool {
	return c == '_' || '0' <= c && c <= '9' || 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z'
}
