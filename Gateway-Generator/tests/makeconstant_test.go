package test

import (
	"create"
	"idl_gen"
	"log"
	"os"
	"reflect"
	"testing"
)

// Testing MakeConstant Function with Hello.thrift content
func TestMakeConstant(t *testing.T) {
	testIDL := "hellotest.thrift"
	serviceinfo := idl_gen.ServiceInfo{IDLName: testIDL}

	err := os.Chdir("../testfiles")
	if err != nil {
		log.Fatal("move to test folder failed")
	}

	result := idl_gen.MakeConstants("gateway", serviceinfo)

	expected := &create.Constants{
		FilepathToService:   "../idl/hellotest.thrift",
		ServiceName:         "HelloService",
		Methods:             []create.Method{{MethodName: "HelloMethod"}},
		IDLName:             "hello",
		GatewayName:         "gateway",
		Load_Balancing_Type: "",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test failed. \nExpected: %+v\nGot: %+v", expected, result)
	} else {
		t.Log("Make Constant for hello.thrift ok")
	}

}

// Testing MakeConstant Function with bizrequests.thrift content (Multiple Methods)
func TestMakeConstant2(t *testing.T) {
	testIDL := "bizrequeststest.thrift"
	serviceinfo := idl_gen.ServiceInfo{IDLName: testIDL}

	err := os.Chdir("../testfiles")
	if err != nil {
		log.Fatal("move to test folder failed")
	}

	result := idl_gen.MakeConstants("gateway", serviceinfo)

	expected := &create.Constants{
		FilepathToService:   "../idl/bizrequeststest.thrift",
		ServiceName:         "BizService",
		Methods:             []create.Method{{MethodName: "BizMethod1"}, {MethodName: "BizMethod2"}, {MethodName: "BizMethod3"}},
		IDLName:             "http",
		GatewayName:         "gateway",
		Load_Balancing_Type: "",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Test failed. \nExpected: %+v\nGot: %+v", expected, result)
	}
	t.Log("Make Constant for bizrequests.thrift ok")
}
