package hz_gen

import (
	"fmt"
	"log"
	"os/exec"
)

// Run go mod tidy
func Tidy() {
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("go mod tidy failed with %s\n", err)
	}

}

// Build gateway into exe file
func Build(name string) {
	cmd := exec.Command("go", "build", name)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("making exe failed with %s\n", err)
	}
	cmd = exec.Command("go", "build", "-o", name)
	err = cmd.Run()
	if err != nil {
		log.Fatalf("making exe failed with %s\n", err)
	}
	fmt.Print(name + " created\n")
}
