package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var version string

type Job struct {
	ID *string
}

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	jobHclFile := flag.String("job-hcl-file", `nomad-job.hcl`, "Path to the job hch file")
	nomadUrl := flag.String("nomad-url", `http://127.0.0.1:4646`, "Parse url")
	showVersion := flag.Bool("version", false, "Prints the version of the application and exits")
	flag.Parse()
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

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

	response, err := http.Post(*nomadUrl+"/v1/jobs/parse", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == 500 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("Error parsing hcl file: " + string(body))
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	updateJob(*nomadUrl, contents)
}

//Updates a job in nomad
//Parameter nomadUrl	The url of nomad which to target
//Parameter hclContent	The content of a hcl formatted file
func updateJob(nomadUrl string, hclContent []byte) {
	var jobFull map[string]interface{}
	json.Unmarshal(hclContent, &jobFull)

	job := new(Job)
	json.Unmarshal(hclContent, &job)

	message := map[string]interface{}{
		"Job": jobFull,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	response, err := http.Post(nomadUrl+"/v1/job/"+*job.ID, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == 500 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("Error updating job: " + string(body))
	}

	if response.StatusCode == 200 {
		log.Println(message)
		os.Exit(0)
	} else {
		os.Exit(1)
	}
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
