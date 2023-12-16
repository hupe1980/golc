package main

import (
	"context"
	"fmt"
	"log"
	"os"

	cohereclient "github.com/cohere-ai/cohere-go/v2/client"
	"github.com/hupe1980/golc/documentcompressor"
	"github.com/hupe1980/golc/schema"
)

func main() {
	docs := []schema.Document{
		{PageContent: "Apples are pomaceous fruits that come in a variety of colors, such as red, green, and yellow. They are known for their sweet or tart taste, crisp texture, and are a popular choice for snacks, desserts, and beverages."},
		{PageContent: "Bananas are elongated, curved fruits with a distinctive yellow skin when ripe. They belong to the genus Musa and are cultivated worldwide for their sweet taste, creamy texture, and high nutritional value."},
		{PageContent: "Apples are versatile in the kitchen, used in a range of culinary creations from classic apple pies and sauces to modern salads and juices. Their natural sweetness and refreshing quality make them a staple in both sweet and savory dishes."},
		{PageContent: "Bananas are a versatile ingredient in the culinary world, enjoyed in a multitude of ways. Whether eaten fresh, blended into smoothies, or baked into desserts, their mild flavor and natural sweetness make them a popular choice for a variety of dishes."},
		{PageContent: "Bananas are a tropical fruit, thriving in warm climates and originating in Southeast Asia. They are now a global commodity, with major banana-producing countries including Ecuador, the Philippines, and Costa Rica, contributing to the widespread availability of this beloved fruit."},
		{PageContent: "Apples have a rich history, dating back thousands of years and playing a significant role in various cultures and mythologies. From the biblical story of Adam and Eve to the legend of Isaac Newton's gravity revelation, apples symbolize knowledge, temptation, and discovery."},
	}

	client := cohereclient.NewClient(cohereclient.WithToken(os.Getenv("COHERE_API_KEY")))

	cohereRerank := documentcompressor.NewCohereRank(client, func(o *documentcompressor.CohereRerankOptions) {
		o.TopN = 2
	})

	compressedDocs, err := cohereRerank.Compress(context.Background(), docs, "All about apples")
	if err != nil {
		log.Fatal(err)
	}

	for _, cd := range compressedDocs {
		fmt.Println(cd.PageContent)
		fmt.Printf("Score: %f\n\n", cd.Metadata["relevanceScore"].(float64))
	}
}
