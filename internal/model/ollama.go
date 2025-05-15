package model

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
)

func ServeOllamaModel() error{
    err := stopOl
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
