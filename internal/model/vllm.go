package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func ServeVLLMModel(name, port string, stream bool) (Model, error){
    os.Setenv("OLLAMA_HOST", fmt.Sprintf("localhost:%s", port))
    cmd := exec.Command("bash", "-c", "ollama serve")
    err := cmd.Start()
    if err != nil {
        return Model{}, err
    }
    cmd = exec.Command("bash", fmt.Sprintf("ollama run %s", name))
    err = cmd.Run()

    m := Model{
        Name:           name,
        Endpoint:       fmt.Sprintf("http://localhost:%s/api/generate", port),
        Stream:         false,
    }
    return m, nil
}

func StopVLLMService() error {
    cmd := exec.Command("bash", "-c", "sudo systemctl stop ollama.service")
    err := cmd.Run()
    if err != nil {
        return errors.New(fmt.Sprintf("failed to stop ollama: %s", err))
    }
    log.Println("Ollama stopped successfully")
    return nil
}

func (m *Model) PromptVLLM(p []byte) (ModelResponse, error) {
    log.Print(string(p))
    dat, err := json.Marshal(
        ModelRequest{
            Model:      m.Name,
            Prompt:     string(p),
            Stream:     m.Stream,
        })
    if err != nil {
        return ModelResponse{}, err
    }
    buf := bytes.NewBuffer(dat)
    log.Printf(m.Endpoint)
    resp, err := http.Post(m.Endpoint, "application/json", buf)
    if err != nil {
        return ModelResponse{}, err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return ModelResponse{}, err
    }
    var modelResp ModelResponse
    err = json.Unmarshal(body, &modelResp)
    if err != nil {
        log.Printf("failed to unmarshal model response, %s", err)
        return ModelResponse{}, err
    }
    return modelResp, nil
}
