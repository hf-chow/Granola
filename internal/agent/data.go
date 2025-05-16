package agent

import "github.com/hf-chow/tofu/internal/model"

type Agent struct{
    Name    string
    Topic   string 
    Model   model.Model
}
