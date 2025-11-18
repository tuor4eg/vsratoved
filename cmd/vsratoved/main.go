package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/tuor4eg/vsratoved/internal/config"
	"github.com/tuor4eg/vsratoved/internal/llm"
)

func main() {
	mode := flag.String("mode", "clean", "mode: clean or spicy")
	flag.Parse()

	ctx := context.Background()
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	advice, err := llm.GetWeirdAdvice(ctx, *mode)
	if err != nil {
		fmt.Println(llm.ErrorFallbackMessage())
		return
	}
	if advice.Author != "" && advice.Advice != "" {
		fmt.Printf("%s - %s\n", advice.Advice, advice.Author)
	} else if advice.Advice != "" {
		fmt.Println(advice.Advice)
	}
}
