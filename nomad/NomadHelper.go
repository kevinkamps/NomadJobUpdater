package nomad

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Job struct {
	ID *string
}

type NomadHelper struct {
	configuration *Configuration
	httpClient    *http.Client
}

func NewNomadHelper(configuration *Configuration) *NomadHelper {
	nomadHelper := NomadHelper{}
	nomadHelper.configuration = configuration

	nomadHelper.httpClient = &http.Client{}
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = *configuration.AllowInsecureCertificates

	if *configuration.TlsAuthEnabled {
		// Load client cert
		cert, err := tls.LoadX509KeyPair(*configuration.TlsCertFile, *configuration.TlsKeyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Load CA cert
		caCert, err := ioutil.ReadFile(*configuration.TlsCaFile)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig.Certificates = []tls.Certificate{cert}
		tlsConfig.RootCAs = caCertPool
		tlsConfig.BuildNameToCertificate()
	}
	nomadHelper.httpClient.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	return &nomadHelper
}

func (this *NomadHelper) ParseHclJob(bodyReader io.Reader) []byte {
	response := this.doPostRequest(*this.configuration.Url+"/v1/jobs/parse", bodyReader)
	if response.StatusCode == 500 {
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal("Error parsing hcl file: " + string(body))
	}

	jsonContent, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	return jsonContent
}

//Updates a job in nomad
//Parameter jsonContent	The content of a json formatted job file
func (this *NomadHelper) UpdateJob(jsonContent []byte) {
	var jobFull map[string]interface{}
	json.Unmarshal(jsonContent, &jobFull)

	job := new(Job)
	json.Unmarshal(jsonContent, &job)

	message := map[string]interface{}{
		"Job": jobFull,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}
	response := this.doPostRequest(*this.configuration.Url+"/v1/job/"+*job.ID, bytes.NewBuffer(bytesRepresentation))
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

func (this *NomadHelper) doPostRequest(url string, bodyReader io.Reader) *http.Response {
	req, err := http.NewRequest("POST", url, bodyReader)
	if *this.configuration.BasicAuthEnabled {
		req.SetBasicAuth(*this.configuration.BasicAuthUsername, *this.configuration.BasicAuthPassword)
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := this.httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
