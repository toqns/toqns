package main

import (
	"fmt"

	"github.com/toqns/toqns/app/tooling/wallet/toqns-cli/cmd"
	"github.com/toqns/toqns/business/logo"
)

func main() {
	logo.Print()
	fmt.Println()
	fmt.Println("Copyright (c) 2022, Toqns Inc.")

	cmd.Execute()
}
