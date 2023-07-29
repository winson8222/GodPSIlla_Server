package test

import (
	"hz_gen"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMyFunction(t *testing.T) {
	err := os.Chdir("../gateway")
	if err != nil {
		log.Fatal("Error occur")
	}
	hz_gen.ClearGateway()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal("Error occur")
	}

	for _, file := range files {
		filename := file.Name()
		if !strings.HasPrefix(filename, "gateway") {
			t.Error("left over files in gateway")
		}
	}
	t.Log("Gateway files are cleared")
}
