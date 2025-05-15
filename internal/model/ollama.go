package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func ServeOllamaModel() error{
    err := stopOllama()
    if err != nil {
        return err
    }
    cmd := exec.Command("bash", "-c", "ollama serve")
    err = cmd.Start()
    if err != nil {
        return err
    }
    return nil
}

func stopOllama() error {
    cmd := exec.Command("bash", "-c", "pkill -f ollama")
    err := cmd.Run()
    if err != nil {
        return errors.New(fmt.Sprintf("failed to stop ollama: %s", err))
    }
    log.Println("Ollama stopped successfully")
    return nil
}

func pullOllamaModel(name string) error {
    cmd := exec.Command("bash", "-c", fmt.Sprintf("ollama pull %s", name))
    err := cmd.Run()
    if err != nil {
        return errors.New(fmt.Sprintf("failed to pull a %s: %s", name, err))
    }
    return nil
}

func (m *Model) prompt(msg ChatMessage) (ChatResponse, error) {
    msgs := []ChatMessage{}
    msgs = append(msgs, msg)

    dat, err := json.Marshal(
        ChatRequest{
            Model:      m.Name,
            Messages:   msgs,
            Stream:     m.Stream,
        })
    if err != nil {
        return ChatResponse{}, err
    }
    buf := bytes.NewBuffer(dat)
    resp, err := http.Post(m.Endpoint, "application/json", buf)
    if err != nil {
        return ChatResponse{}, err
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return ChatResponse{}, err
    }
    var chatResp ChatResponse
    err = json.Unmarshal(body, &chatResp)
    if err != nil {
        return ChatResponse{}, err
    }
    return chatResp, nil
}

