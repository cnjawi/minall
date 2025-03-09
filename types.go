package main

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"` // always true
}

type Choice struct {
	Delta struct {
		Content          string `json:"content"`
		ReasoningContent string `json:"reasoning_content"`
		FunctionCall     any    `json:"function_call"`
		Refusal          any    `json:"refusal"`
		Role             string `json:"role"`
		ToolCalls        []any  `json:"tool_calls"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason"`
	Index        int    `json:"index"`
	LogProbs     any    `json:"logprobs"` // even not in the official doc
}

// 文档: https://help.aliyun.com/zh/model-studio/developer-reference/use-qwen-by-calling-api
type ResponseBody struct {
	Id                string   `json:"id"`
	Choices           []Choice `json:"choices"`
	Created           int      `json:"created"`
	Model             string   `json:"model"`
	Object            string   `json:"object"`
	Usage             any      `json:"usage"`
	SystemFingerprint any      `json:"system_fingerprint"`
}

func (r *ResponseBody) getContent() string {
	return r.Choices[0].Delta.Content
}

func (r *ResponseBody) getReasoningContent() string {
	return r.Choices[0].Delta.ReasoningContent
}

type TransOpt struct {
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	Domains    string `json:"domains"`
}

type TranslationReq struct {
	Model              string    `json:"model"`
	Messages           []Message `json:"messages"`
	TranslationOptions TransOpt  `json:"translation_options"`
	Stream             bool      `json:"stream"`
}
