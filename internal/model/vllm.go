package model

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

type VLLMModel struct {
    Name        string
    Port        string
    Device      DeviceType
}

type DeviceType string 

const (
    DeviceCPU DeviceType = "cpu"
    DeviceGPU DeviceType = "gpu"
)

func (d DeviceType) Validate() error {
    switch d {
    case DeviceGPU, DeviceCPU: 
        return nil
    default:
        return errors.New("invalid device type: %s, supported types are [cpu, gpu]")
    }
}

func NewVLLMModel(name, port string, device DeviceType) *VLLMModel {
    return &VLLMModel{
        Name:    name,
        Port:    port,
        Device:  device,
    }
}

func (m *VLLMModel) Start() error {
    hfToken := os.Getenv("HUGGING_FACE_HUB_TOKEN")
    if hfToken == "" {
        return errors.New("missing HuggingFace Hub token in environment variable " +
        "use `export HUGGING_FACE_HUB_TOKEN=<your HuggingFace Hub token>` to set " +
        "your token to run the model with vLLM")
    }
    switch m.Device {
    case DeviceCPU:
        cmd := exec.Command(
            "docker build -f docker/Dockerfile.cpu --tag vllm --target vlllm-openai .",
        )
        err := cmd.Run()
        if err != nil {
            return err
        }
        cmd = exec.Command(
            "docker", 
            fmt.Sprintf("run --name=vllm --rm --privileged=true -p %s:8000", m.Port),
            fmt.Sprintf("--env 'HUGGING_FACE_HUB_TOKEN=%s'", hfToken),
            fmt.Sprintf("vllm --model=%s", m.Name),
        )
        err = cmd.Run()
        if err != nil {
            return err
        }

    case DeviceGPU:
        cmd := exec.Command(
            "DOCKER_BUILDKIT=1 docker build . --target vllm-openai --tag vllm --file docker/Dockerfile",
        )
        err := cmd.Run()
        if err != nil {
            return err
        }
        cmd = exec.Command(
            "docker", 
            fmt.Sprintf("run --name=vllm --rm --privileged=true -p %s:8000", m.Port),
            fmt.Sprintf("--env 'HUGGING_FACE_HUB_TOKEN=%s'", hfToken),
            fmt.Sprintf("vllm --model=%s", m.Name),
        )
        err = cmd.Run()
        if err != nil {
            return err
        }
    default:
        return errors.New("invalid device type: %s, supported types are [cpu, gpu]")
    }
    return nil
}

func (m *VLLMModel) Stop() error {
    cmd := exec.Command(
        "docker stop vllm",
    )
    err := cmd.Run() 
    if err != nil {
        return err
    }
    return nil
}

func (m *VLLMModel) GetModelInfo() string {
    return m.Name
}
