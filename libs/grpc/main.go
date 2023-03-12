package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	os.RemoveAll("generated")

	files, err := filepath.Glob("./protos/*.proto")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.Contains(file, "google") {
			continue
		}

		parts := strings.Split(file, "/")
		if len(parts) != 2 {
			log.Fatalf("Folder structure %s doesn't supported", file)
		}
		fname := parts[1]

		parts = strings.Split(fname, ".")
		if len(parts) != 2 {
			log.Fatalf("File name %s is not expected", file)
		}
		name := parts[0]

		if _, err := os.Stat(fmt.Sprintf("generated/%s", file)); os.IsNotExist(err) {
			os.MkdirAll(fmt.Sprintf("generated/%s", name), os.ModePerm)
		}

		cmd := exec.Command("protoc",
			"--go_opt=paths=source_relative",
			fmt.Sprintf("--go-grpc_out=./generated/%s", name),
			"--go-grpc_opt=paths=source_relative",
			"--experimental_allow_proto3_optional",
			fmt.Sprintf("--go_out=./generated/%s", name),
			"--proto_path=./protos",
			file)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal("Error running protoc:", string(output))
		}

		log.Println("✅ Generated", file, "proto definitions", string(output))
	}
}
