package controller

import (
	"apiProject/api/response"
	"apiProject/api/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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
	router.HandleFunc("/llms/ollama", td.handleOllama).Methods("POST")
	router.HandleFunc("/llms/simple", td.handleSimple).Methods("POST")
	router.HandleFunc("/llms/stream", td.handleStream).Methods("POST")
}

func (td *LlmsController) handleOllama(w http.ResponseWriter, r *http.Request) {
	var body utils.Chat
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("大模型问答参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("大模型问答请求失败"))
		return
	}

	// 使用了defer会在请求完成后处理关闭
	defer utils.CloseBodyError("大模型问答请求body", w, r)
	log.Println("大模型问答请求body内容===", body)

	prompt := utils.CreatePrompt()

	data := map[string]any{
		"text": body.Prompt,
	}

	msg, _ := prompt.FormatMessages(data)
	messageContent := []llms.MessageContent{
		llms.TextParts(msg[0].GetType(), msg[0].GetContent()),
		llms.TextParts(msg[1].GetType(), msg[1].GetContent()),
	}

	content, _ := td.Llm.GenerateContent(context.Background(), messageContent)
	log.Println("返回内容", content.Choices[0].Content)

	response.WriteJson(w, response.OkDataResp(content))
}

func (td *LlmsController) handleSimple(w http.ResponseWriter, r *http.Request) {
	var body utils.Chat
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("大模型问答参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("大模型问答请求失败"))
		return
	}

	// 使用了defer会在请求完成后处理关闭
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("大模型问答请求body关闭错误===%v", err)
			response.WriteJson(w, response.FailMessageResp("请求体关闭失败"))
			return
		}
		log.Println("大模型问答请求body关闭成功")
	}(r.Body)

	log.Println("大模型问答请求body内容===", body)
	utils.CreateSimplePrompt(td.Llm, w, body)
}

func (td *LlmsController) handleStream(w http.ResponseWriter, r *http.Request) {
	var body utils.Chat
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("大模型问答参数解析失败===%v", err)
		response.WriteJson(w, response.FailMessageResp("大模型问答请求失败"))
		return
	}

	// 使用了defer会在请求完成后处理关闭
	defer utils.CloseBodyError("大模型问答请求body", w, r)
	log.Println("大模型问答请求body内容===", body)

	prompt := utils.CreatePrompt()

	data := map[string]any{
		"text": body.Prompt,
	}

	msg, _ := prompt.FormatMessages(data)
	messageContent := []llms.MessageContent{
		llms.TextParts(msg[0].GetType(), msg[0].GetContent()),
		llms.TextParts(msg[1].GetType(), msg[1].GetContent()),
	}

	ctx := context.Background()
	flusher, ok := w.(http.Flusher)
	if !ok {
		response.WriteJson(w, response.FailMessageResp("Streaming not supported"))
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.WriteHeader(http.StatusOK)

	_, err = td.Llm.GenerateContent(ctx, messageContent, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		_, writeErr := w.Write(chunk)
		fmt.Print(string(chunk))
		if writeErr != nil {
			return writeErr
		}
		flusher.Flush()
		return nil
	}))

	if err != nil {
		log.Printf("生成内容时出现错误==%v", err)
		response.WriteJson(w, response.FailMessageResp("生成内容失败"))
	}
}

func Read(resp *http.Response, out chan<- interface{}) error {
	r := resp.Body
	defer r.Close()
	cr := httputil.NewChunkedReader(r)
	for {
		m := &http.Response{}
		buf := make([]byte, 1024)
		l, err := cr.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if pe := json.Unmarshal(buf[:l], m); pe != nil {
			return pe
		}
		out <- m
		if err == io.EOF {
			return nil
		}
	}
}
