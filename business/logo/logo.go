// Package logo prints a CLI logo.
package logo

import "fmt"

// Print prints a Toqns ASCII art logo.
func Print() {
	fmt.Println(`  █████                                         `)
	fmt.Println(` ░░███                                          `)
	fmt.Println(` ███████    ██████   ████████ ████████    █████ `)
	fmt.Println(`░░░███░    ███░░███ ███░░███ ░░███░░███  ███░░  `)
	fmt.Println(`  ░███    ░███ ░███░███ ░███  ░███ ░███ ░░█████ `)
	fmt.Println(`  ░███ ███░███ ░███░███ ░███  ░███ ░███  ░░░░███`)
	fmt.Println(`  ░░█████ ░░██████ ░░███████  ████ █████ ██████ `)
	fmt.Println(`   ░░░░░   ░░░░░░   ░░░░░███ ░░░░ ░░░░░ ░░░░░░  `)
	fmt.Println(`                        ░███                    `)
	fmt.Println(`                        █████                   `)
	fmt.Println(`                       ░░░░░                    `)
}
