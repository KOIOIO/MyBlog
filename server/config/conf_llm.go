package config

// LLM 大模型配置（阿里云百炼 DashScope，OpenAI 兼容接口）。
type LLM struct {
	BaseURL         string `json:"base_url" yaml:"base_url"`                 // 兼容接口 Base URL
	APIKey          string `json:"api_key" yaml:"api_key"`                   // API Key（支持 ${ENV} 环境变量展开）
	Model           string `json:"model" yaml:"model"`                       // 模型名，如 qwen-max
	MaxHistory      int    `json:"max_history" yaml:"max_history"`           // 注入上下文的历史消息条数上限
	MaxArticleChars int    `json:"max_article_chars" yaml:"max_article_chars"` // 单篇文章注入正文的最大字符数
}
