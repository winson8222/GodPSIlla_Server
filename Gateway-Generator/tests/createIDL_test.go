package test

import (
	"create"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

// Test if the IDL file can be update successfully
func TestCreateIDL1(t *testing.T) {
	testConstants1 := create.Constants{
		ServiceName:       "HelloService",
		FilepathToService: "../testfiles/idl/hello.thrift",
		Methods: []create.Method{
			{MethodName: "HelloMethod"},
		},
	}

	original, err := createBackup(testConstants1.FilepathToService)
	if err != nil {
		log.Fatal("back up cannot be created")
	}

	create.CreateIDL(testConstants1)

	match, err := CompareFileContent("../testfiles/idltest/hello.thrift", "../testfiles/idl/hello.thrift")
	if err != nil {
		log.Fatal("comparing failed")
	}

	if match {
		err = restoreFromBackup("../testfiles/idl/hello.thrift", original)
		if err != nil {
			fmt.Print("restore from backup failed")
		}
		t.Log("IDL file updated")
	} else {
		err = restoreFromBackup("..../testfiles/idl/hello.thrift", original)
		if err != nil {
			fmt.Print("restore from backup failed")
		}
		t.Error("IDL file updating failed")
	}

}

// Test if the IDL file can be update successfully with multiple methods
func TestCreateIDL2(t *testing.T) {
	testConstants1 := create.Constants{
		ServiceName:       "BizService",
		FilepathToService: "../testfiles/idl/bizrequests.thrift",
		Methods: []create.Method{
			{MethodName: "BizMethod1"},
			{MethodName: "BizMethod2"},
			{MethodName: "BizMethod3"},
		},
	}

	original, err := createBackup(testConstants1.FilepathToService)
	if err != nil {
		log.Fatal("back up cannot be created")
	}

	create.CreateIDL(testConstants1)

	match, err := CompareFileContent("../testfiles/idltest/bizrequests.thrift", "../testfiles/idl/bizrequests.thrift")
	if err != nil {
		log.Fatal("comparing failed")
	}

	if match {
		err = restoreFromBackup("../testfiles/idl/bizrequests.thrift", original)
		if err != nil {
			fmt.Print("restore from backup failed")
		}
		t.Log("IDL file updated")
	} else {
		err = restoreFromBackup("..../testfiles/idl/bizrequests.thrift", original)
		if err != nil {
			fmt.Print("restore from backup failed")
		}
		t.Error("IDL file updating failed")
	}

}

func createBackup(filePath string) ([]byte, error) {
	// Read the content of the file.
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Create a backup file with the content of the original file.
	backupPath := filePath + ".backup"
	err = ioutil.WriteFile(backupPath, content, 0644)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func restoreFromBackup(filePath string, content []byte) error {
	// Restore the content of the original file from the backup.
	err := ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// CompareFileContent compares the content of two files and returns true if their contents match.
func CompareFileContent(filePath1, filePath2 string) (bool, error) {
	// Read the content of the first file.
	content1, err := ioutil.ReadFile(filePath1)
	if err != nil {
		return false, err
	}

	// Read the content of the second file.
	content2, err := ioutil.ReadFile(filePath2)
	if err != nil {
		return false, err
	}

	fmt.Print(string(content1))
	fmt.Print(string(content2))

	// Compare the contents as strings (ignoring leading/trailing whitespaces).
	return strings.TrimSpace(string(content1)) == strings.TrimSpace(string(content2)), nil
}
