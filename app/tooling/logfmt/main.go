// This program takes the structured log output and makes it readable.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

var service string

func init() {
	flag.StringVar(&service, "service", "", "filter which service to see")
}

func main() {
	flag.Parse()
	var b strings.Builder

	// Scan standard input for log data per line.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()

		// Convert the JSON to a map for processing.
		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// If a service filter was provided, check.
		if service != "" && m["service"] != service {
			continue
		}

		// Build out the know portions of the log in the order
		// I want them in.
		b.Reset()
		b.WriteString(fmt.Sprintf("[%s] %s %s %s: %s ",
			color.BlueString(m["service"].(string)),
			color.CyanString(m["ts"].(string)),
			coloredLevel(m["level"].(string)),
			color.BlueString(m["caller"].(string)),
			color.YellowString(m["msg"].(string)),
		))

		// Add the rest of the keys ignoring the ones we already
		// added for the log.
		for k, v := range m {
			switch k {
			case "service", "ts", "level", "caller", "msg":
				continue
			}

			// TODO: Check the additional trailing space.
			// Without it the last character gets truncated.
			b.WriteString(fmt.Sprintf("%s[%v]  ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}

func coloredLevel(l string) string {
	switch strings.ToLower(l) {
	case "error":
		return color.RedString("%s", l)
	case "warning":
		return color.YellowString("%s", l)
	case "debug":
		return color.MagentaString("%s", l)
	default:
		return color.GreenString("%s", l)
	}
}
