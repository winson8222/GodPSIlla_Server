package test

import (
	"create"
	"fmt"
	"log"
	"testing"
)

// Testing NignxConfig function with Single Serivce of Single Method
func TestNconf(t *testing.T) {
	testservices := create.Services{
		GATEWAY_URL:         "0.0.0.0:80",
		ETCD_URL:            "0.0.0.0:20000",
		LOAD_BALANCING_TYPE: "ROUND_ROBIN",
		Service_Constants: []create.Constants{{
			FilepathToService:   "../idl/hellotest.thrift",
			ServiceName:         "HelloService",
			Methods:             []create.Method{{MethodName: "HelloMethod"}},
			IDLName:             "hello",
			GatewayName:         "gateway",
			Load_Balancing_Type: ""}},
	}

	original, err := createBackup("../nginx/conf/nginx.conf")
	if err != nil {
		log.Fatal("create back up failed")
	}

	create.NginxConfig(testservices)

	result, err := CompareFileContent("../nginx/conf/nginx.conf", "../testfiles/nginx/nginxtest1.conf")
	if err != nil {
		t.Error(err)
	}

	err = restoreFromBackup("../nginx/conf/nginx.conf", original)
	if err != nil {
		fmt.Print("restore backup failed")
	}

	if result {
		t.Log("nconf is correct")
	} else {
		t.Error("nconf is not correct")
	}

}

// Testing NignxConfig function with Multiple Serivce of Multiple Methods
func TestNconf2(t *testing.T) {
	testservices := create.Services{
		GATEWAY_URL:         "0.0.0.0:80",
		ETCD_URL:            "0.0.0.0:20000",
		LOAD_BALANCING_TYPE: "ROUND_ROBIN",
		Service_Constants: []create.Constants{{
			FilepathToService:   "../idl/hellotest.thrift",
			ServiceName:         "HelloService",
			Methods:             []create.Method{{MethodName: "HelloMethod"}},
			IDLName:             "hello",
			GatewayName:         "gateway",
			Load_Balancing_Type: ""},
			{
				FilepathToService:   "../idl/bizrequeststest.thrift",
				ServiceName:         "BizService",
				Methods:             []create.Method{{MethodName: "BizMethod1"}, {MethodName: "BizMethod2"}, {MethodName: "BizMethod3"}},
				IDLName:             "http",
				GatewayName:         "gateway",
				Load_Balancing_Type: ""}},
	}

	original, err := createBackup("../nginx/conf/nginx.conf")
	if err != nil {
		log.Fatal("create back up failed")
	}

	create.NginxConfig(testservices)

	result, err := CompareFileContent("../nginx/conf/nginx.conf", "../testfiles/nginx/nginxtest2.conf")
	if err != nil {
		t.Error(err)
	}

	err = restoreFromBackup("../nginx/conf/nginx.conf", original)
	if err != nil {
		fmt.Print("restore backup failed")
	}

	if result {
		t.Log("nconf is correct")
	} else {
		t.Error("nconf is not correct")
	}

}
