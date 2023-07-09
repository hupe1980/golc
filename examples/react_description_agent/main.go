package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/agent"
	"github.com/hupe1980/golc/agent/toolkit"
	"github.com/hupe1980/golc/model/llm"
	"github.com/playwright-community/playwright-go"
)

func main() {
	golc.Verbose = true

	if err := playwright.Install(); err != nil {
		log.Fatal(err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatal(err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatal(err)
	}

	openai, err := llm.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	browserKit, err := toolkit.NewBrowser(browser)
	if err != nil {
		log.Fatal(err)
	}

	agent, err := agent.NewReactDescription(openai, browserKit.Tools())
	if err != nil {
		log.Fatal(err)
	}

	result, err := golc.SimpleCall(context.Background(), agent, "Navigate to https://news.ycombinator.com and summarize the text")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
