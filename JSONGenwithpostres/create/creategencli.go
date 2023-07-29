package create

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"unicode"
)

// Create generic Client file based on the constants of a service
func Creategencli(constants Constants) {
	// Generic client generation template

	funcs := template.FuncMap{
		"title": func(s string) string {
			for i, v := range s {
				return string(unicode.ToUpper(v)) + s[i+1:]
			}
			return ""
		}, "ToConstant": ToConstant,
	}

	genericClientTemplate := `
package {{ .IDLName }}

import (
	"context"
	"{{ .GatewayName }}/constants"
	"log"
	"strings"
	"fmt"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	

	"github.com/cloudwego/kitex/pkg/klog"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}


// Creates generic client "[ServiceName]GenericClient"
func {{ .ServiceName }}GenericClient() genericclient.Client {
	r, err := etcd.NewEtcdResolver([]string{constants.ETCD_URL})
	if err != nil {
		log.Fatal(err)
	}

	path := constants.FILEPATH_TO_{{ .ServiceName | ToConstant }}
	p, err := generic.NewThriftFileProvider(path)
	if err != nil {
		klog.Fatalf("new thrift file provider failed: %v", err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		klog.Fatalf("new JSON thrift generic failed: %v", err)
	}

	if (constants.LOAD_BALANCING == "ROUND_ROBIN") {
		cli, err := genericclient.NewClient(constants.{{ .ServiceName | ToConstant }}_NAME, g, client.WithResolver(r),
			client.WithLoadBalancer(loadbalance.NewWeightedBalancer()))
		if err != nil {
			klog.Fatalf("new JSON generic client failed: %v", err)
		}
		return cli
	} else {
		cli, err := genericclient.NewClient(constants.{{ .ServiceName | ToConstant }}_NAME, g, client.WithResolver(r))
		if err != nil {
			klog.Fatalf("new JSON generic client failed: %v", err)
		}
		return cli
	}
}

`

	// Method template
	methodTemplate := `
func Do{{ title .MethodName }}(ctx context.Context, cli genericclient.Client, req string) (interface{}, error) {
	fmt.Print(req)
	resp, err := cli.GenericCall(ctx, "{{ .MethodName }}", req)

	if err != nil {
		return nil, err
	}
	//OWN CODE ABOVE
	return resp, nil
}

`

	// Create the output file
	outputFile, err := os.Create("biz/handler/" + constants.IDLName + "/gen_client.go")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	// Create a new template for the generic client
	clientTmpl := template.Must(template.New("genericClient").Funcs(funcs).Parse(genericClientTemplate))

	// Execute the generic client template with the constants and write the output to the file
	err = clientTmpl.Execute(outputFile, constants)
	if err != nil {
		log.Fatalf("Error executing generic client template: %s\n", err)
	}

	// Create a new template for the method
	methodTmpl := template.Must(template.New("method").Funcs(funcs).Parse(methodTemplate))

	// Generate code for each method
	for _, method := range constants.Methods {
		// Execute the method template with the current method and write the output to the file
		err = methodTmpl.Execute(outputFile, method)
		if err != nil {
			log.Fatalf("Error executing method template: %s\n", err)
		}
	}

	fmt.Print("Generic Client code for " + constants.ServiceName + " Generated successfully.\n")
}
