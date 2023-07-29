// Package gen contains code generating files needed to create customised Godzilla Gateway
package main

import (
	"create"
	"hz_gen"
	"idl_gen"
	"log"
	"os"
)

func main() {

	//Clear IDL files
	idl_gen.ClearFolder("idl")

	// Download IDL files and get info
	info, serviceinfolist := idl_gen.GetIDL()

	//Reconfigure IDL files

	// Create Required stuct for service related info
	gatewayexample := idl_gen.MakeServices(info, serviceinfolist)

	//create the constant folder and files
	desiredDir := "gateway"
	err := os.Chdir(desiredDir)
	if err != nil {
		log.Fatalf("move to folder failed with %s\n", err)
	}

	//Setup Nginx config
	create.NginxConfig(gatewayexample)
	// //create gencli for all services
	for _, constant := range gatewayexample.Service_Constants {
		create.CreateIDL(constant)
	}
	hz_gen.Hzgen()
	for _, constant := range gatewayexample.Service_Constants {
		create.Creategencli(constant)
	}
	create.CreateConstant(gatewayexample)
	create.CreateMain()

	allhandlers := []create.HandlerInfo{}

	err = os.Chdir("../")
	if err != nil {
		log.Fatalf("move to directory folder failed with %s\n", err)
	}
	//Make Handler info for all every IDLs
	for _, idls := range serviceinfolist {
		allhandlers = append(allhandlers, idl_gen.MakeHandlerInfo(idls.IDLName, info.GatewayName))

	}

	//create handler for all methods
	for _, handler := range allhandlers {
		create.Createhandler(handler)
		err := os.Chdir("../")
		if err != nil {
			log.Fatalf("move to directory folder failed with %s\n", err)
		}
	}

	err = os.Chdir("gateway")
	if err != nil {
		log.Fatalf("move to directory folder failed with %s\n", err)
	}

	//change kitex v0.6.1 to v0.5.2 and netpoll v0.4.0 to v0.3.2 to fix some weird bug
	Update()

	//go mod tidy
	hz_gen.Tidy()

	//build exe
	hz_gen.Build(info.GatewayName)

}
