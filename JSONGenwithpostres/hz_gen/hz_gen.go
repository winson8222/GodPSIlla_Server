package hz_gen

import (
	"fmt"
	"idl_gen"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Install hz for the local system
func Hzinstall() {
	// execute install hz
	cmd1 := exec.Command("go", "install", "github.com/cloudwego/hertz/cmd/hz@latest")
	err := cmd1.Run()
	if err != nil {
		log.Fatalf("install hertz failed with %s\n", err)
	}

	// Set the environment variable for the second command
	os.Setenv("GO111MODULE", "on")

	// execute install thriftgo
	cmd2 := exec.Command("go", "install", "github.com/cloudwego/thriftgo@latest")
	err = cmd2.Run()
	if err != nil {
		log.Fatalf("install thriftgo failed with %s\n", err)
	}
}

// Generate the hz based gateway under gateway folder
func Hzgen() {
	//check idl files
	IDLs, err := idl_gen.GetIDLs()
	if err != nil {
		log.Fatalf("get IDL files failed with %s\n", err)
	}

	ClearGateway()

	number := 0

	for _, file := range IDLs {

		if number <= 0 {
			cmd1 := exec.Command("hz", "new", "-module", "gateway", "-idl", "../idl/"+file)
			err = cmd1.Run()
			if err != nil {
				log.Fatalf("hz gen failed with %s\n", err)
			}
			Tidy()
			number += 1
		} else {
			cmd1 := exec.Command("hz", "update", "-idl", "../idl/"+file)
			err = cmd1.Run()
			if err != nil {
				log.Fatalf("hz update failed with %s\n", err)
			}
			Tidy()
		}

	}

}

func ClearGateway() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to retrieve current working directory:", err)
	}

	// Delete directories recursively
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error occurred:", err)
			return err
		}

		if info.IsDir() && path != wd {
			err := os.RemoveAll(path)
			if err != nil {
				log.Println("Failed to delete directory:", path, "-", err)
			} else {
				log.Println("Deleted directory:", path)
			}
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// Delete files (excluding directories)
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatal("Failed to read directory:", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			if !strings.HasPrefix(filename, "gateway") {
				filePath := filepath.Join(wd, filename)
				err := os.Remove(filePath)
				if err != nil {
					log.Println("Failed to delete file:", filePath, "-", err)
				} else {
					log.Println("Deleted file:", filePath)
				}
			}
		}
	}

	log.Println("Deletion completed.")
}
