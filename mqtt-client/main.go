package main

import (
	"client/cmd"
	"os"
)

func main() {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	cmd.Execute()
}
