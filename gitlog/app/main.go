package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type Output struct {
	Added    []string `json:"added,omitempty"`
	Modified []string `json:"modified,omitempty"`
	Removed  []string `json:"removed,omitempty"`
}

func main() {

	out, err := exec.Command("git", "log", "--name-status", "--stat", "-1", "-m", "--oneline").Output()
	if err != nil {
		log.Fatal(err)
	}

	sp := strings.Split(string(out), "\n")

	// trim empty array cells
	sp = sp[1:]
	if len(sp) > 0 {
		sp = sp[:len(sp)-1]
	}

	output := Output{}

	for _, spl := range sp {
		tmp := strings.Split(spl, "\t")

		switch tmp[0] {
		case "A":
			output.Added = append(output.Added, tmp[1])
			break
		case "M":
			output.Modified = append(output.Modified, tmp[1])
			break
		case "D":
			output.Removed = append(output.Removed, tmp[1])
			break
		}

	}

	// Jsonify output
	b, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}
