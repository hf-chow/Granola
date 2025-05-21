package model

type Model interface {
    Generate(prompt string) (string, error)
    Start() error
    Stop() error
    GetModelInfo() string
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

func Prompt(prompt string, m Model) (string, error) {
    resp, err := m.Generate(prompt)
    if err != nil {
        return "", err
    }
    return resp, nil
}
