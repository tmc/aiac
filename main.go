package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

var (
	flagModel   = flag.String("model", "gpt-3.5-turbo", "model to use")
	flagVerbose = flag.Bool("verbose", false, "verbose output")
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: [description of the code you want to generate]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var prompts = map[bool]string{
	// non-verbose
	false: "Generate sample code for a %s",
	// verbose
	true: "Generate sample code for a %s. Include explanations.",
}

func run() error {
	llm, err := openai.NewChat(openai.WithModel(*flagModel))
	if err != nil {
		return err
	}
	ctx := context.Background()
	_, err = llm.Call(ctx, []schema.ChatMessage{
		schema.HumanChatMessage{Content: fmt.Sprintf(prompts[*flagVerbose], strings.Join(flag.Args(), " "))},
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}),
	)
	return err
}
