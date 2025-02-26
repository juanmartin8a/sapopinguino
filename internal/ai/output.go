package aiutils;

type Token struct {
	Type          string `json:"type"`
	Input         string `json:"input,omitempty"`
	Transcription string `json:"transcription,omitempty"`
	Output        string `json:"output,omitempty"`
	Value         string `json:"value,omitempty"`
}

// TokensResponse represents the complete JSON structure returned by GPT-4.
type TokensResponse struct {
	Tokens []Token `json:"tokens"`
}
