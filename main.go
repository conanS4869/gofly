package main

import "gofly/cmd"

func main() {
	defer cmd.Clean()
	cmd.Start()
}
