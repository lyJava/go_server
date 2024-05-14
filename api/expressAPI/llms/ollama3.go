package llms

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func MessageLlama3() {
	ctx := context.Background()

	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Printf("使用模型出现错误%v", err)
	}

	completion, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		"golang的发展前景，使用中文回答",
		//llms.WithTemperature(0.8),
		//// 这里是返回字符串
		//llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		//	fmt.Print(string(chunk))
		//	return nil
		//}),
	)

	fmt.Print("\r\n")
	if err != nil {
		log.Printf("出现错误==%v", err)
	}
	fmt.Println(completion)

	//_ = completion
}
