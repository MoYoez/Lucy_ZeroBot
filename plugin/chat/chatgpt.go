package chat // Package chat base on Zerobot-Plugin Playground.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FloatTech/zbputils/img/text"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"

	"github.com/FloatTech/ttl"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	modelGPT3Dot5Turbo = "gpt-3.5-turbo"
)

type sessionKey struct {
	group int64
	user  int64
}

var cache = ttl.NewCache[sessionKey, []chatMessage](time.Minute * 15)

var GPTPromptHandlerLimitedTimeManager = rate.NewManager[int64](time.Minute*2, 3)

// chatGPTResponseBody 响应体
type chatGPTResponseBody struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int          `json:"created"`
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
}

// chatGPTRequestBody 请求体
type chatGPTRequestBody struct {
	Model       string        `json:"model,omitempty"` // gpt3.5-turbo
	Messages    []chatMessage `json:"messages,omitempty"`
	Temperature float64       `json:"temperature,omitempty"`
	N           int           `json:"n,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

// chatMessage 消息
type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatChoice struct {
	Index        int `json:"index"`
	Message      chatMessage
	FinishReason string `json:"finish_reason"`
}

// 可以直接使用web封包 无需二次构建
// completions gtp3.5文本模型回复
// curl https://api.openai.com/v1/chat/completions
// -H "Content-Type: application/json"
// -H "Authorization: Bearer YOUR_API_KEY"
// -d '{ "model": "gpt-3.5-turbo",  "messages": [{"role": "user", "content": "Hello!"}]}'
func completions(messages []chatMessage, apiKey string) (*chatGPTResponseBody, error) {
	com := chatGPTRequestBody{
		Messages: messages,
	}
	// default model
	if com.Model == "" {
		com.Model = modelGPT3Dot5Turbo
	}
	body, err := json.Marshal(com)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://openai.lemonkoi.one/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	var client http.Client
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	v := new(chatGPTResponseBody)
	if err = json.NewDecoder(res.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}

func init() {
	// easy and work well with chatgpt? key handler.
	// trigger for chatgpt.
	engine.OnRegex(`^/chat\s*(.*)`, zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !GPTPromptHandlerLimitedTimeManager.Load(ctx.Event.GroupID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Too quick! 慢一点再请求哦！"))
			return
		}
		args := ctx.State["regex_matched"].([]string)[1]
		key := sessionKey{
			group: ctx.Event.GroupID,
			user:  ctx.Event.UserID,
		}
		if args == "reset" {
			cache.Delete(key)
			ctx.Send("Session cleaned.")
			return
		}
		messages := cache.Get(key)
		messages = append(messages, chatMessage{
			Role:    "user",
			Content: args,
		})
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("收到信息了哦~大概需要等一下w..."))
		resp, err := completions(messages, os.Getenv("gptkey"))
		if err != nil {
			ctx.SendChain(message.Text("出现了一些问题:( : ", err))
			return
		}
		reply := resp.Choices[0].Message
		reply.Content = strings.TrimSpace(reply.Content)
		messages = append(messages, reply)
		cache.Set(key, messages)
		base64Format := HandleTextTobase54UsingWryh(reply.Content, ctx)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+helper.BytesToString(base64Format)))
	})
}

func HandleTextTobase54UsingWryh(txt string, ctx *zero.Ctx) (base64Raw []byte) {
	base64Raw, err := text.RenderToBase64(txt, text.FontPath+"wryh.ttf", 1280, 45)
	if err != nil {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Handle text err.", err))
		return nil
	}
	return base64Raw
}
