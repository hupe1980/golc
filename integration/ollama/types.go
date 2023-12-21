package ollama

import (
	"time"

	"github.com/hupe1980/golc/integration/stream"
)

type ErrorResponse struct {
	Message string `json:"error"`
}

type Metrics struct {
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

type Runner struct {
	NumCtx             int     `json:"num_ctx,omitempty"`
	NumBatch           int     `json:"num_batch,omitempty"`
	NumGQA             int     `json:"num_gqa,omitempty"`
	NumGPU             int     `json:"num_gpu,omitempty"`
	MainGPU            int     `json:"main_gpu,omitempty"`
	NumThread          int     `json:"num_thread,omitempty"`
	RopeFrequencyBase  float32 `json:"rope_frequency_base,omitempty"`
	RopeFrequencyScale float32 `json:"rope_frequency_scale,omitempty"`
	LogitsAll          bool    `json:"logits_all,omitempty"`
	VocabOnly          bool    `json:"vocab_only,omitempty"`
	UseMMap            bool    `json:"use_mmap,omitempty"`
	UseMLock           bool    `json:"use_mlock,omitempty"`
	EmbeddingOnly      bool    `json:"embedding_only,omitempty"`
	UseNUMA            bool    `json:"numa,omitempty"`
	F16KV              bool    `json:"f16_kv,omitempty"`
	LowVRAM            bool    `json:"low_vram,omitempty"`
}

type Options struct {
	Stop []string `json:"stop,omitempty"`
	Runner
	RepeatLastN      int     `json:"repeat_last_n,omitempty"`
	Seed             int     `json:"seed,omitempty"`
	TopK             int     `json:"top_k,omitempty"`
	NumKeep          int     `json:"num_keep,omitempty"`
	Mirostat         int     `json:"mirostat,omitempty"`
	NumPredict       int     `json:"num_predict,omitempty"`
	Temperature      float32 `json:"temperature,omitempty"`
	TypicalP         float32 `json:"typical_p,omitempty"`
	RepeatPenalty    float32 `json:"repeat_penalty,omitempty"`
	PresencePenalty  float32 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
	TFSZ             float32 `json:"tfs_z,omitempty"`
	MirostatTau      float32 `json:"mirostat_tau,omitempty"`
	MirostatEta      float32 `json:"mirostat_eta,omitempty"`
	TopP             float32 `json:"top_p,omitempty"`
	PenalizeNewline  bool    `json:"penalize_newline,omitempty"`
}

type ImageData []byte

type GenerationRequest struct {
	Model    string      `json:"model"`
	Prompt   string      `json:"prompt"`
	System   string      `json:"system"`
	Template string      `json:"template"`
	Context  []int       `json:"context,omitempty"`
	Stream   *bool       `json:"stream,omitempty"`
	Raw      bool        `json:"raw,omitempty"`
	Format   string      `json:"format"`
	Images   []ImageData `json:"images,omitempty"`

	Options Options `json:"options"`
}

type GenerationResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Response  string    `json:"response"`

	Done    bool  `json:"done"`
	Context []int `json:"context,omitempty"`

	Metrics
}

type GenerationStream struct {
	*stream.Stream[GenerationResponse]
}

type Message struct {
	Role    string      `json:"role"` // one of ["system", "user", "assistant"]
	Content string      `json:"content"`
	Images  []ImageData `json:"images,omitempty"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   *bool     `json:"stream,omitempty"`
	Format   string    `json:"format"`

	Options Options `json:"options"`
}

type ChatResponse struct {
	Model     string    `json:"model"`
	CreatedAt time.Time `json:"created_at"`
	Message   *Message  `json:"message,omitempty"`

	Done bool `json:"done"`

	Metrics
}

type ChatStream struct {
	*stream.Stream[ChatResponse]
}

type EmbeddingRequest struct {
	Model   string  `json:"model"`
	Prompt  string  `json:"prompt"`
	Options Options `json:"options"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}
