package idl_gen

import (
	"create"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"unicode"
)

// Gateway struct contains necessary infromation of the gateway
type GatewayInfo struct {
	GatewayPort         string
	ETCD_URL            string
	GatewayName         string
	Load_Balancing_Type string
}

// ServiceInfo structs contains information needed for the gateway's connection to microservices
type ServiceInfo struct {
	IDLName string
}

// MakeServices built information needed for creation of constants files
func MakeServices(info GatewayInfo, list []ServiceInfo) create.Services {

	exampleconstants := []create.Constants{}

	for _, service := range list {
		constants := MakeConstants(info.GatewayName, service)
		exampleconstants = append(exampleconstants, *constants)
	}

	url := "0.0.0.0:" + info.GatewayPort

	gateway := create.Services{
		GATEWAY_URL:         url,
		ETCD_URL:            info.ETCD_URL,
		LOAD_BALANCING_TYPE: info.Load_Balancing_Type,
		Service_Constants:   exampleconstants,
	}
	fmt.Print("Gateway info configured\n")
	return gateway
}

// MakeHandlerInfo returns information needed for create handler functions for a service in the gateway
func MakeHandlerInfo(idl string, gatename string) create.HandlerInfo {
	serviceinfo := create.HandlerServiceInfo{
		IDLName:     GetNameSpace(idl),
		GatewayName: gatename,
		HandlerFile: separateCamelCase(GetServiceName(idl)),
	}

	methods := GetMethods(idl)

	handlers := []create.Handler{}

	for _, method := range methods {
		handler := create.Handler{
			MethodName:    method.MethodName,
			ServiceName:   GetServiceName(idl),
			IDLName:       GetNameSpace(idl),
			RequestStruct: GetReqStruct(idl),
		}
		handlers = append(handlers, handler)
	}

	handlerinfo := create.HandlerInfo{
		ServiceInfo: serviceinfo,
		Handlers:    handlers,
	}

	return handlerinfo

}

// Get names of all the idl files retrieved from database
func GetIDLs() ([]string, error) {

	files, err := ioutil.ReadDir("../idl")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	IDLs := []string{}
	for _, file := range files {
		IDLs = append(IDLs, file.Name())
	}

	return IDLs, nil
}

// Get name of the IDL file without extension
func GetIDLName(idl string) string {

	name := idl[:len(idl)-7]

	return name
}

// Get filepath of the IDL file
func GetFilePath(idl string) string {

	path := "./idl/" + idl

	return path
}

// Get the name of the request struct using IDL
func GetReqStruct(idl string) string {
	path := GetFilePath(idl)

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read file " + path + " fail while getting req Struct")
	}

	// Convert byte slice to string
	fileString := string(content)

	typeRegex := regexp.MustCompile(`service\s+\w+\s*{[\s\n]*\w+\s+\w+\s*\(\d+:\s*[a-z]*\s*([A-Z]\w*)`)
	matches := typeRegex.FindStringSubmatch(fileString)
	if len(matches) >= 2 {
		requestType := matches[1]
		return requestType
	}

	log.Fatal("cannot find request type")
	return "none"
}

// Get the ServiceName using IDL
func GetServiceName(idl string) string {
	path := GetFilePath(idl)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read file " + path + " fail while getting service name")
	}

	stringcontent := string(content)
	serviceRegex := regexp.MustCompile(`service\s+(\w+)\s+{`)
	matches := serviceRegex.FindStringSubmatch(stringcontent)

	if len(matches) >= 2 {
		serviceName := matches[1]
		return serviceName
	}

	return ""
}

// Get the names of the methods of a service from IDL file
func GetMethods(idl string) []create.Method {
	path := GetFilePath(idl)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read file " + path + "fail while getting methods")
	}

	stringContent := string(content)

	methodRegex := regexp.MustCompile(`(\w+)\s+(\w+)\(.*?\);?`)
	matches := methodRegex.FindAllStringSubmatch(stringContent, -1)

	methods := []create.Method{}
	for _, match := range matches {
		methodName := match[2]

		newMethod := create.Method{
			MethodName: methodName,
		}

		methods = append(methods, newMethod)
	}
	return methods
}

// game namespace from thrift idl file (Package in go)
func GetNameSpace(idl string) string {
	path := GetFilePath(idl)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("read file " + path + "fail while getting namespace")
	}

	stringcontent := string(content)
	// Extract namespace
	namespaceRegex := regexp.MustCompile(`namespace\s+go\s+(\w+)`)
	match := namespaceRegex.FindStringSubmatch(stringcontent)

	if len(match) >= 2 {
		namespace := match[1]
		return namespace
	}

	return ""

}

// Create constants object based on service information
func MakeConstants(gateway string, info ServiceInfo) *create.Constants {
	con := create.Constants{
		FilepathToService:   "." + GetFilePath(info.IDLName),
		ServiceName:         GetServiceName(info.IDLName),
		Methods:             GetMethods(info.IDLName),
		IDLName:             GetNameSpace(info.IDLName),
		GatewayName:         gateway,
		Load_Balancing_Type: "",
	}

	return &con

}

func separateCamelCase(input string) string {
	var builder strings.Builder
	for i, char := range input {
		if i > 0 && unicode.IsUpper(char) {
			builder.WriteRune('_')
		}
		builder.WriteRune(unicode.ToLower(char))
	}
	return builder.String()
}
