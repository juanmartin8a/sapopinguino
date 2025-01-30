package aiutils;

type Output struct {
	Mappings         []Mapping      `json:"mappings"`
}

type Mapping struct {
	Input            string         `json:"input"`
	Transcription    string         `json:"transcription"`
	Output           string         `json:"output"`
}
