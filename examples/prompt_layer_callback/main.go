package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/go-promptlayer"
	"github.com/hupe1980/golc"
	"github.com/hupe1980/golc/callback"
	"github.com/hupe1980/golc/chain"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/prompt"
	"github.com/hupe1980/golc/schema"
)

func main() {
	client := promptlayer.NewClient(os.Getenv("PROMPTLAYER_API_KEY"))

	output, err := client.GetPromptTemplate(context.Background(), &promptlayer.GetPromptTemplateInput{
		PromptName: "joke",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ID:", output.ID)
	fmt.Println("Template:", output.PromptTemplate.Template)
	fmt.Println("InputVariables:", output.PromptTemplate.InputVariables)

	openai, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	llmChain, err := chain.NewLLM(openai, prompt.NewTemplate(output.PromptTemplate.Template, func(o *prompt.TemplateOptions) {
		o.TransformPythonTemplate = true
	}))
	if err != nil {
		log.Fatal(err)
	}

	cb := callback.NewPromptLayerHandler(os.Getenv("PROMPTLAYER_API_KEY"), func(o *callback.PromptLayerHandlerOptions) {
		o.PromptID = output.ID
		o.OnPromptLayerOutputFunc = func(output *promptlayer.TrackRequestOutput) error {
			fmt.Printf("The request was tracked under id %s\n", output.RequestID)
			return nil
		}
	})

	result, err := golc.SimpleCall(context.Background(), llmChain, "Tell me a joke about animals. Max. 10 words.", func(sco *golc.SimpleCallOptions) {
		sco.Callbacks = []schema.Callback{cb}
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}
