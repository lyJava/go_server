package controller

import (
	"apiProject/api/response"
	"apiProject/api/utils"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
	"net/http"
)

// LlmsController 大模型控制器
type LlmsController struct {
	Llm *ollama.LLM
}

// LlmsControllerInit 大模型控制器初始化
func LlmsControllerInit(llm *ollama.LLM) *LlmsController {
	return &LlmsController{
		Llm: llm,
	}
}

func (td *LlmsController) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/llms/ollama", td.handleOllama3Ask).Methods("POST")
}

func (td *LlmsController) handleOllama3Ask(w http.ResponseWriter, r *http.Request) {
	var body utils.Chat
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("大模型问答参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("大模型问答请求失败"))
		return
	}

	// 使用了defer会在请求完成后处理关闭
	defer utils.CloseBodyError("大模型问答请求body", w, r)
	/*defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("大模型问答请求body关闭错误===%v", err)
			response.WriteJson(w, response.FailMessageResp("请求体关闭失败"))
			return
		}
		log.Println("大模型问答请求body关闭成功")
	}(r.Body)*/

	log.Println("大模型问答请求body内容===", body)

	prompt := utils.CreatePrompt()

	data := map[string]any{
		"text": body.Text,
	}

	msg, _ := prompt.FormatMessages(data)
	messageContent := []llms.MessageContent{
		llms.TextParts(msg[0].GetType(), msg[0].GetContent()),
		llms.TextParts(msg[1].GetType(), msg[1].GetContent()),
	}

	content, err := td.Llm.GenerateContent(context.Background(), messageContent)
	log.Println("返回内容", content.Choices[0].Content)

	response.WriteJson(w, response.OkDataResp(content))
}
