package utils

import (
	"apiProject/api/response"
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
	"net/http"
)

type Chat struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`  // 模型名称
	Format string `json:"format,omitempty"` // 返回响应的格式。目前唯一接受的值是json
	Stream bool   `json:"stream,omitempty"` // false响应是否作为单个响应对象返回，而不是对象流
	Raw    bool   `json:"raw,omitempty"`    // true没有格式化，将应用于提示。raw如果您在 API 请求中指定完整模板化提示，则可以选择使用该参数
	Images string `json:"images,omitempty"` // base64编码的图片
}

// CreateModel 创建大模型
func CreateModel(moduleName string) *ollama.LLM {
	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Printf("创建模型出现错误%v", err)
	}
	return llm
}

// CreatePrompt 创建消息模板
func CreatePrompt() prompts.ChatPromptTemplate {
	return prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate("", nil),
		prompts.NewHumanMessagePromptTemplate("{{.text}}", []string{
			"text",
		}),
	})
}

// CreateSimplePrompt 创建请求简单模版
func CreateSimplePrompt(llm *ollama.LLM, w http.ResponseWriter, content Chat) {
	ctx := context.Background()

	flusher, ok := w.(http.Flusher)
	if !ok {
		response.WriteJson(w, response.FailMessageResp("Streaming not supported"))
		return
	}

	//w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	//w.WriteHeader(http.StatusOK)

	_, err := llms.GenerateFromSinglePrompt(
		ctx,
		llm,
		content.Prompt,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			_, writeErr := w.Write(chunk)
			if writeErr != nil {
				fmt.Printf("w.Write(chunk) error==%v", writeErr)
				return writeErr
			}
			flusher.Flush()
			return nil
		}),
	)

	fmt.Print("\r\n")
	if err != nil {
		log.Printf("生成内容时出现错误==%v", err)
		response.WriteJson(w, response.FailMessageResp("生成内容失败"))
		return
	}
}

// CreateSimplePromptChan generates text from a given prompt using the provided LLM.
// It returns a channel from which the generated text can be read line by line.
func CreateSimplePromptChan(llm *ollama.LLM, prompt string) (<-chan string, error) {
	ctx := context.Background()
	lineChan := make(chan string)

	go func() {
		defer close(lineChan) // Ensure the channel is closed when done

		_, err := llms.GenerateFromSinglePrompt(
			ctx,
			llm,
			prompt,
			llms.WithTemperature(0.8),
			llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
				line := string(chunk)
				fmt.Print(line)
				lineChan <- line
				return nil
			}),
		)

		if err != nil {
			log.Printf("出现错误==%v", err)
			lineChan <- fmt.Sprintf("error: %v", err)
		}

		fmt.Print("\r\n")
	}()

	return lineChan, nil
}

func PromptChan(ll *ollama.LLM, prompt string) {
	lineChan, err := CreateSimplePromptChan(ll, prompt)
	if err != nil {
		log.Fatalf("Failed to create prompt: %v", err)
	}

	// Reading from the channel line by line
	for line := range lineChan {
		fmt.Println("Received:", line)
		//time.Sleep(1 * time.Second) // Simulating processing delay
	}
}
