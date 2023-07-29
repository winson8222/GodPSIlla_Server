package create

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"text/template"
)

// Constants struct contains information needed for building constant.go in gateway
type Constants struct {
	FilepathToService   string
	ServiceName         string
	Methods             []Method
	IDLName             string
	GatewayName         string
	Load_Balancing_Type string
}

// Services struct contains necessary information for creating constants.go files
type Services struct {
	GATEWAY_URL         string
	ETCD_URL            string
	LOAD_BALANCING_TYPE string
	Service_Constants   []Constants
}

// Method struct contains information on a method in a service
type Method struct {
	MethodName string
}

// CreateConstant create constant.go files needed for connection to all services
func CreateConstant(services Services) {
	// Define the values for the constants

	// Define the template string

	templateString :=
		`package constants

	 import (
		"strings"
	 )
 
	func ToConstant(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
	}


	const (
		ETCD_URL = "{{ .ETCD_URL }}" //connects to a single etcd instance in the cluster
		LOAD_BALANCING = "{{ .LOAD_BALANCING_TYPE }}"
	)`

	servicetemplate :=
		`
	const (
		FILEPATH_TO_{{ .ServiceName | ToConstant }}  = "{{ .FilepathToService }}"
		{{ .ServiceName | ToConstant }}_NAME         = "{{ .ServiceName }}" //name registered with svc registry as rpcendpoint
	)
	`

	// ToConstant function converts the input string to a constant format

	// Create a new template
	tmpl := template.Must(template.New("constants").Funcs(template.FuncMap{"ToConstant": ToConstant}).Parse(templateString))

	// Create the output file
	err := os.MkdirAll("constants", os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output folder: %s\n", err)
	}
	outputFile, err := os.Create("constants/constants.go")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	// Execute the template with the constants and write the output to the file
	err = tmpl.Execute(outputFile, services)
	if err != nil {
		log.Fatalf("Error executing constants template: %s\n", err)
	}

	serviceTmpl := template.Must(template.New("services_constants").Funcs(template.FuncMap{"ToConstant": ToConstant}).Parse(servicetemplate))

	// Generate code for each method
	for _, constants := range services.Service_Constants {
		// Execute the method template with the current method and write the output to the file
		err = serviceTmpl.Execute(outputFile, constants)
		if err != nil {
			log.Fatalf("Error executing service constants template: %s\n", err)
		}
	}

	log.Println("Template generation completed successfully.")
}

// Configure the Nginx conf file based on the services and methods routes
func NginxConfig(services Services) {
	configString := `
events {
	worker_connections  1024;
}

http {
	upstream gateway {
		server 127.0.0.1:8888;
		server 127.0.0.1:8889;
		server 127.0.0.1:8890;
	}

	server {
		listen {{ .PORT }};
{{- range .Service_Constants }}
{{- $serviceName := .ServiceName }}
{{- range .Methods }}
		location /{{ $serviceName }}/{{ .MethodName }} {
			proxy_pass http://gateway/{{ $serviceName }}/{{ .MethodName }};
		}
{{- end }}
{{- end }}
		location /ping {
			proxy_pass http://gateway/ping;
		}
	}
}
`

	// Extract the port from GATEWAY_URL
	_, port, err := net.SplitHostPort(services.GATEWAY_URL)
	if err != nil {
		log.Fatal("invalid Address")
	}

	// Define the data to be passed to the template
	data := struct {
		PORT              string
		Service_Constants []Constants
	}{
		PORT:              port,
		Service_Constants: services.Service_Constants,
	}

	// Create a new template
	tmpl := template.Must(template.New("nginxConfig").Parse(configString))

	// Execute the template and write to the output file
	outputFile, err := os.Create("../nginx/conf/nginx.conf")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatalf("Error executing template: %s\n", err)
	}

	fmt.Println("NGINX configuration file generated successfully.")
}

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}
