package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hupe1980/golc/documenttransformer"
	"github.com/hupe1980/golc/model/chatmodel"
	"github.com/hupe1980/golc/schema"
)

type Tagging struct {
	MovieTitle string `json:"movieTitle" description:"The title of the movie this critic is for"`
	Critic     string `json:"critic"`
	Tone       string `json:"tone" enum:"'happy','neutral','sad'"`
	Rating     int    `json:"rating" description:"The number of stars the critic rated the movie"`
}

func main() {
	chatModel, err := chatmodel.NewOpenAI(os.Getenv("OPENAI_API_KEY"), func(o *chatmodel.OpenAIOptions) {
		o.Temperature = 0
	})
	if err != nil {
		log.Fatal(err)
	}

	tagger, err := documenttransformer.NewMetaDataTagger(chatModel, &Tagging{})
	if err != nil {
		log.Fatal(err)
	}

	docs := []schema.Document{
		{PageContent: "Review of The Bee Movie\nBy Roger Ebert\n\nThis is the greatest movie ever made. 4 out of 5 stars."},
		{PageContent: "Review of The Godfather\nBy Anonymous\n\nThis movie was super boring. 1 out of 5 stars.", Metadata: map[string]any{
			"reliable": false,
		}},
	}

	enrichedDocs, err := tagger.Transform(context.Background(), docs)
	if err != nil {
		log.Fatal(err)
	}

	for _, ed := range enrichedDocs {
		fmt.Println(ed.PageContent)
		fmt.Println(ed.Metadata)
		fmt.Println("---")
	}
}
