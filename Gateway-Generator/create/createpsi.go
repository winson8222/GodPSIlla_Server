package create

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"unicode"
)

func ConfigMid() {

	configfile, err := os.Create("middleman/config.txt")
	if err != nil {
		log.Fatalf("Error creating middleman config file: %s\n", err)
	}

	_, err = configfile.WriteString(os.Args[1])
	if err != nil {
		log.Fatalf("Error Writing middleman config file: %s\n", err)
	}

}

func CreatePSI(services []HandlerInfo) {
	funcs := template.FuncMap{
		"title": func(s string) string {
			for i, v := range s {
				return string(unicode.ToUpper(v)) + s[i+1:]
			}
			return ""
		},
	}
	fmt.Print(services[0].ServiceInfo.IDLName)
	outputFile, err := os.Create("gateway/biz/handler/psi.go")
	if err != nil {
		log.Fatalf("Error creating psi file: %s\n", err)
	}
	defer outputFile.Close()
	psitemplate :=
		`package handler

	import (
		"context"
		"encoding/json"
		"fmt"
		`
	outputFile.WriteString(psitemplate)

	serivceimport := `"gateway/biz/handler/{{.ServiceInfo.IDLName}}"
	`
	importtemp := template.Must(template.New("method").Parse(serivceimport))
	for _, method := range services {
		err = importtemp.Execute(outputFile, method)
		if err != nil {
			log.Fatalf("Error executing method template: %s\n", err)
		}
	}

	importemp := `
		"github.com/cloudwego/hertz/pkg/app"
		"github.com/cloudwego/hertz/pkg/protocol/consts"
	)
	`

	_, err = outputFile.WriteString(importemp)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	maptemp := `var functionMap = make(map[string]map[string]func(ctx context.Context, c *app.RequestContext))
	`

	_, err = outputFile.WriteString(maptemp)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	temp := `
	func InitPSISvcInfo() {
		`
	_, err = outputFile.WriteString(temp)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	// Initialize the functionMap with predefined contents
	svcstring := `
	functionMap["{{(index .Handlers 0).ServiceName}}"] = map[string]func(ctx context.Context, c *app.RequestContext){
		`

	methodstring := `"{{.MethodName}}": {{.IDLName}}.{{ title .MethodName }},`

	quotetemp := `
	}`
	svctemp := template.Must(template.New("services2").Parse(svcstring))
	methodtemp := template.Must(template.New("methods").Funcs(funcs).Parse(methodstring))

	for _, service := range services {
		err = svctemp.Execute(outputFile, service)
		if err != nil {
			log.Fatalf("Error executing method template: %s\n", err)
		}
		for _, method := range service.Handlers {
			err = methodtemp.Execute(outputFile, method)
			if err != nil {
				log.Fatalf("Error executing method template: %s\n", err)
			}
		}

		_, err = outputFile.WriteString(quotetemp)
	}
	_, err = outputFile.WriteString(quotetemp)

	// PSI.
	reststring := `
	func PSI(ctx context.Context, c *app.RequestContext) {
		//First extract microservice and method names from c
		jsonBody, _ := c.Body()
		svcinfo, err := extractSvcInfoAsMap(jsonBody)
		if err != nil {
			return
		}
	
		//Then make generic calls to the microservices
		// Note that the Request struct of the methods being called is just "{}" so we can generalise it
		var resList [][]string
		for svcname, mtname := range svcinfo {
			if functionMap[svcname][mtname] != nil { //method exists
				ftn := functionMap[svcname][mtname]
				reqContext := app.NewContext(0) //What is maxParams supposed to be?
				reqContext.Request.SetBodyString("{}")
	
				ftn(context.Background(), reqContext) //Call the method
	
				jsonResp := reqContext.Response.Body()
	
				resList = append(resList, extractValueAsMap(jsonResp))
			} else {
				c.String(consts.StatusBadRequest, "Service or Method name not found for: "+svcname+", "+mtname)
				return
			}
		}
	
		//FINDING INTERSECTION OF RESPONSES
		//print out values received, set the min size array idx arbitrarily
		//listsChecked is the number of lists checked, as part of the algorithm
		minSizeListIdx := 0
		listsChecked := 0
		// then find and assign the min size list index
		for i, _ := range resList {
			if len(resList[i]) < len(resList[minSizeListIdx]) {
				minSizeListIdx = i
			}
		}
		//create a hashmap with the elements of the minSizeList
		intersectionMap := make(map[string]int)
		for _, elem := range resList[minSizeListIdx] {
			intersectionMap[elem] = 1
		}
		listsChecked = 1
		//The moment any req fails, immediately terminate and return the error
		//Once all data is collected, find the smaller set
		for i, list := range resList {
			if i == minSizeListIdx {
				continue
			} else {
				Intersection(&list, intersectionMap, &listsChecked)
			}
		}
	
		// Extract values
		// Create a slice to store the keys that map to the target value
		var matchingKeys []string
	
		// Iterate through the map to find keys that map to the target value
		for key, value := range intersectionMap {
			if value == listsChecked {
				matchingKeys = append(matchingKeys, key)
			}
		}
	
		//Return response
		fmt.Println(matchingKeys)
		c.JSON(consts.StatusOK, matchingKeys)
	}
	
	// a is the smaller list, and b is the larger list
	// use pointers to conserve memory
	func Intersection(list *[]string, intersectionMap map[string]int, listsChecked *int) {
		for _, elem := range *list {
			if intersectionMap[elem] == *listsChecked {
				intersectionMap[elem] += 1
			}
		}
		*listsChecked += 1
	}
	
	// Takes the json response in bytes and extracts the list of strings
	func extractValueAsMap(jsonResponse []byte) []string {
	
		// Declare a variable to store the parsed JSON data
		// Unmarshal the JSON string to remove Go string escaping
		var unescapedJSON string
		if err := json.Unmarshal(jsonResponse, &unescapedJSON); err != nil {
			fmt.Println("Error:", err)
			return nil
		}
	
		// Unmarshal the unescaped JSON string to parse the JSON data
		var data map[string][]string
		if err := json.Unmarshal([]byte(unescapedJSON), &data); err != nil {
			fmt.Println("Error:", err)
			return nil
		}
	
		for _, v := range data {
			return v
		}
		//if somehow there isn't data returned
		return nil
	}
	
	func extractSvcInfoAsMap(jsonBody []byte) (map[string]string, error) {
		// Parse the JSON body from the request
		resMap := make(map[string]string)
		var stringList [][]string
		if err := json.Unmarshal(jsonBody, &stringList); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return nil, err
		}
	
		for _, l := range stringList {
			resMap[l[0]] = l[1]
		}
		return resMap, nil
	}`
	_, err = outputFile.WriteString(reststring)

	routeString := `// Code generated by hertz generator.

	package main
	
	import (
		handler "gateway/biz/handler"
	
		"github.com/cloudwego/hertz/pkg/app/server"
	)
	
	// customizeRegister registers customize routers.
	func customizedRegister(r *server.Hertz) {
		r.GET("/ping", handler.Ping)
		r.POST("/PSI", handler.PSI)
	
		// your code ...
	}
	`
	routerfile, err := os.Create("gateway/router.go")
	if err != nil {
		log.Fatalf("Error Creating router file")
	}
	defer routerfile.Close()

	_, err = routerfile.WriteString(routeString)
	if err != nil {
		panic(err)
	}

	log.Println("psi.go has been created")
}
