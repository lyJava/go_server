package utils

import (
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/prompts"
	"log"
)

type Chat struct {
	Text string `json:"text"`
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
