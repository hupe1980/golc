package main

import (
	"fmt"
	"log"

	"github.com/hupe1980/golc/prompt"
)

func main() {
	pt := prompt.NewTemplate(`Tell me a {{.adjective}} joke about {{.content}}.`)

	s, err := pt.Format(map[string]any{
		"adjective": "funny",
		"content":   "chickens",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
