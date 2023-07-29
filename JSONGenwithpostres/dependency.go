package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

// Update the dependency of the generated files and write it into go.mod file of the gateway
func Update() {
	// Read the contents of the go.mod file
	goModFile, err := ioutil.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	// Parse the go.mod file
	modFile, err := modfile.Parse("go.mod", goModFile, nil)
	if err != nil {
		panic(err)
	}

	// Update the required version for github.com/cloudwego/kitex
	updateDependency(modFile, "github.com/cloudwego/kitex", "v0.5.2")

	// Update the required version for github.com/cloudwego/netpoll
	updateDependency(modFile, "github.com/cloudwego/netpoll", "v0.3.2")

	// Generate the updated go.mod content
	updatedContent, err := modFile.Format()
	if err != nil {
		panic(err)
	}

	// Write the updated content back to the go.mod file
	err = ioutil.WriteFile("go.mod", updatedContent, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("go.mod file updated successfully.")

}

// Change the dependencies to compatible versions
func updateDependency(modFile *modfile.File, module, version string) {
	// Add a replace directive for the module
	err := modFile.AddReplace(module, "", module, version)
	if err != nil {
		panic(err)
	}
}
