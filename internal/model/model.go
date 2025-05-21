package model

type ModelProvider interface {
    Generate(prompt string) (string, error)
    GetModelInfo() (ModelInfo, error)
    Stop() error
}

type ModelResponse interface {
    Text()      string
    Metadata()  map[string]interface{}
    Raw()       []byte
}

type ModelInfo struct {
    Name        string
    Endpoint    string
    Stream      bool
}

func ServeModel(m ModelProvider) (Model, error)
func StopModel(m ModelProvider) error
