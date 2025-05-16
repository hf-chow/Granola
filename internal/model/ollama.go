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

func ServeOllamaModel(name, port string, stream bool) (Model, error){
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

func StopOllamaService() error {
    cmd := exec.Command("bash", "-c", "sudo systemctl stop ollama.service")
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


func (m *Model) Prompt(p []byte) (ModelResponse, error) {
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
