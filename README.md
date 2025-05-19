# Granola

An event-driven AI agent framework implemented with Go and RabbitMQ with no other dependencies besides RabbitMQ and Docker.
Designed for high-level abstraction to orchestrate multiple agents. 
Support locally hosted LLM with vLLM and Ollama.

# Prerequisite
- Docker
- RabbitMQ 

# Quick Start
- Run `./rabbit.sh start` to start the RabbitMQ in the background.

## Topics
- `request`
- `quest_ans`
- `prod_search`
- `prod_query`

## Workflow

### Agents
- Orchestration (OR) Agent
- Question Answering (QA) Agent
- Product Search (PS) Agent
- Product Query (PQ) Agent

### Scenario 0: task creation
- OR retrive message from the `request` topic
- OR interpret the request and send message to either `quest_ans`, `prod_search`, and `prod_query` topics
 
### Scenario 1: simple question answering
- QA retreive message from `quest_ans` topic 
- QA answer the question and Ack the message

### Scenario 2: query about stock avaialbility
- PQ retreive from `prod_query` topic 
- PQ query the database and Ack the message
- PQ send message to the `quest_ans` topic 
- The rest follows Scenario 1

### Scenario 3: query about similar product and its stock avaialbility
- PS retrieve from `prod_search` topic
- PS search for the most similar product and Ack the message
- PS queue the task to `prod_query` topic with the search result
- The rest follows Scenario 2
