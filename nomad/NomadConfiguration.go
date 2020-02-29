package nomad

import (
	"flag"
)

type Configuration struct {
	BasicAuthEnabled          *bool
	Url                       *string
	BasicAuthUsername         *string
	BasicAuthPassword         *string
	AllowInsecureCertificates *bool
}

func NewNomadConfiguration() *Configuration {
	config := Configuration{}

	config.Url = flag.String("nomad-url", `http://127.0.0.1:4646`, "Nomad url")
	config.BasicAuthEnabled = flag.Bool("nomad-basic-auth-enabled", false, "Add a basic authentication header to all nomad requests")
	config.BasicAuthUsername = flag.String("nomad-basic-auth-username", `user`, "Basic authentication username")
	config.BasicAuthPassword = flag.String("nomad-basic-auth-password", `password`, "Basic authentication password")
	config.AllowInsecureCertificates = flag.Bool("nomad-allow-insecure-certificates", false, "Allows insecure certificates / self signed certificates")

	return &config
}

func (this *Configuration) Parse() {

}
