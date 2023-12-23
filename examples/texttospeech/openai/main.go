package main

import (
	"context"
	"log"
	"os"

	"github.com/hupe1980/golc/texttospeech"
	"github.com/sashabaranov/go-openai"
)

func main() {
	openai := texttospeech.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *texttospeech.OpenAIOptions) {
		o.Voice = openai.VoiceEcho // optional
	})

	stream, err := openai.SynthesizeSpeech(context.Background(), "Hello world! My name is Echo")
	if err != nil {
		log.Fatal(err)
	}

	defer stream.Close()

	if err := stream.Play(); err != nil {
		log.Fatal(err)
	}
}
