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
)

var version string

type Job struct {
	ID *string
}

func updateJob(nomadUrl string, contents []byte) {
	var jobFull map[string]interface{}
	json.Unmarshal(contents, &jobFull)

	job := new(Job)
	json.Unmarshal(contents, &job)

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

	if response.StatusCode == 200 {
		log.Println(message)
		os.Exit(0)
	} else {
		os.Exit(1)
	}
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
		"JobHCL":       string(hcl),
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

	defer response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	updateJob(*nomadUrl, contents)
}
