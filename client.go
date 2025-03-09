package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func cmpBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func Quest(model Model, messages []Message) string {
	// 创建http客户端
	client := &http.Client{}
	// 构建请求体
	requestBody := RequestBody{
		Model:    model.Name,
		Messages: messages,
		Stream:   true,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		Fatal(err.Error())
	}
	// 创建POST请求
	req, err := http.NewRequest(
		"POST",
		model.Url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		Fatal(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+model.APIKey)
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		Fatal(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}
	// 读取响应体
	var fullContent strings.Builder
	lastContent := ""
	reader := bufio.NewReader(resp.Body)
	for {
		chunk, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			Fatal(err.Error())
		}
		if len(chunk) <= 7 {
			continue // avoid slice bounds out of range
		}
		if cmpBytes(chunk[:6], []byte("data: ")) {
			if chunk[6] == '[' {
				break // meet "data: [Done]"
			}
			chunk = chunk[6:] // remove "data: " to make it a valid json
		} else {
			fmt.Printf("\n 未知数据: %s\n", chunk)
			continue
		}
		var jsonData ResponseBody
		err = json.Unmarshal(chunk, &jsonData)
		if err != nil {
			Fatal(err.Error())
		}
		if jsonData.getReasoningContent() != "" {
			if lastContent == "" {
				fmt.Println("<think>")
				lastContent = "<think>"
			}
			fmt.Print(jsonData.getReasoningContent())
			fmt.Fprint(&fullContent, jsonData.getReasoningContent())
		} else {
			if lastContent == "<think>" {
				fmt.Print("\n</think>\n\n")
				lastContent = ""
			}
			fmt.Print(jsonData.getContent())
			fmt.Fprint(&fullContent, jsonData.getContent())
		}
	}
	return fullContent.String()
}

func chatSession(model Model, systemMsg string) {
	fmt.Print("q to quit")
	var messages []Message
	reader := bufio.NewReader(os.Stdin)
	messages = append(messages, Message{"system", systemMsg})
	for {
		fmt.Print("\n\n>>> ")
		msg, err := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if err != nil {
			Fatal(err.Error())
		}
		if msg == "q" {
			break
		}
		messages = append(messages, Message{"user", msg})
		response := Quest(model, messages)
		messages = append(messages, Message{"assistant", response})
	}
	fmt.Print("\n")
}

// !!NOTE!!: Qwen-MT模型暂时不支持增量式流式输出 (2025-3-9)
// 暂时自行实现增量式输出, 可能会有一些问题
func Translate(model Model, targetLang, domain, text string) {
	client := &http.Client{}
	requestBody := TranslationReq{
		Model: model.Name,
		Messages: []Message{
			{"user", text},
		},
		TranslationOptions: TransOpt{
			SourceLang: "auto",
			TargetLang: targetLang,
			Domains:    domain,
		},
		Stream: true,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		Fatal(err.Error())
	}
	req, err := http.NewRequest(
		"POST",
		model.Url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		Fatal(err.Error())
	}
	req.Header.Set("Authorization", "Bearer "+model.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		Fatal(err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}
	reader := bufio.NewReader(resp.Body)
	lastContentLen := 0 // 记录上一次输出的内容长度，用于增量输出
	for {
		chunk, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			Fatal(err.Error())
		}
		if len(chunk) <= 7 {
			continue // avoid slice bounds out of range
		}
		if cmpBytes(chunk[:6], []byte("data: ")) {
			if chunk[6] == '[' {
				break // meet "data: [Done]"
			}
			chunk = chunk[6:] // remove "data: " to make it a valid json
		} else {
			fmt.Printf("\n 未知数据: %s\n", chunk)
			continue
		}
		var jsonData ResponseBody
		err = json.Unmarshal(chunk, &jsonData)
		if err != nil {
			Fatal(err.Error())
		}
		fmt.Print(jsonData.getContent()[lastContentLen:])
		lastContentLen = len(jsonData.getContent())
	}
}
