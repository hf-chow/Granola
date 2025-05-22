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

type OllamaModel struct {
    Name        string
    Endpoint    string
    Stream      bool
}

func NewOllamaModel(name, port string, stream bool) *OllamaModel {
    return &OllamaModel{
        Name:       name,
        Endpoint:   fmt.Sprintf("http://localhost:%s/api/generate", port),
        Stream:     stream,
    }
}

func (m *OllamaModel) Start() error {
    os.Setenv("OLLAMA_HOST", m.Endpoint)
    cmd := exec.Command("bash", "-c", "ollama serve")
    err := cmd.Start()
    if err != nil {
        return err
    }
    cmd = exec.Command("bash", fmt.Sprintf("ollama run %s", m.Name))

    err = cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func (m *OllamaModel) Stop() error {
    cmd := exec.Command("bash", "-c", "sudo systemctl stop ollama.service")
    err := cmd.Run()
    if err != nil {
        return errors.New(fmt.Sprintf("failed to stop ollama: %s", err))
    }
    log.Println("Ollama stopped successfully")
    return nil
}

func (m *OllamaModel) Generate(prompt string) (string, error) {
    dat, err := json.Marshal(
        OllamaModelRequest{
            Model:      m.Name,
            Prompt:     prompt,
            Stream:     m.Stream,
        })
    if err != nil {
        return "", err
    }
    buf := bytes.NewBuffer(dat)
    log.Printf(m.Endpoint)
    resp, err := http.Post(m.Endpoint, "application/json", buf)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    var modelResp OllamaModelResponse
    err = json.Unmarshal(body, &modelResp)
    if err != nil {
        log.Printf("failed to unmarshal model response, %s", err)
        return modelResp.Response, err
    }
    return modelResp.Response, nil
}

func (m *OllamaModel) GetModelInfo() string { 
    return m.Name
}

func pullOllamaModel(name string) error {
    cmd := exec.Command("bash", "-c", fmt.Sprintf("ollama pull %s", name))
    err := cmd.Run()
    if err != nil {
        return errors.New(fmt.Sprintf("failed to pull a %s: %s", name, err))
    }
    return nil
}

