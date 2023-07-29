package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Removes All the content in the gateway folders
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	fmt.Print("gateway Deleted")
	return nil
}

func main() {
	err := RemoveContents("../gateway")
	if err != nil {
		// Handle error
		fmt.Println(err)
		log.Fatal(err)
	}
}
