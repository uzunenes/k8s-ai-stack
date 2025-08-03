# Kubernetes AI Stack –  Flexible and High-Performance LLM Inference

This repository presents a production-ready AI inference stack for Kubernetes, designed to run advanced large language models (LLMs) of your choice—including but not limited to Qwen3-32B—on NVIDIA L40S GPUs using vLLM. The stack integrates Ollama for embeddings and image analysis, and provides a unified interface through OpenWebUI with secure Microsoft OAuth2 authentication. All services are orchestrated within Kubernetes and securely exposed via Ingress under an internal DNS hostname (e.g., https://ai.com.tr), enabling scalable, high-performance AI capabilities on your own infrastructure.

![](docs/k8s-ai-stack.drawio.png?raw=true)

## Table of Contents
- [Kubernetes AI Stack –  Flexible and High-Performance LLM Inference](#kubernetes-ai-stack---flexible-and-high-performance-llm-inference)
  - [Table of Contents](#table-of-contents)
  - [1. Installation](#1-installation)
    - [1.1 Prerequisites](#11-prerequisites)
    - [1.2 vLLM Deployment on k8s](#12-vllm-deployment-on-k8s)
      - [1.2.1 vLLM Deployment](#121-vllm-deployment)
      - [1.2.2 vLLM Service](#122-vllm-service)
      - [1.2.3 vLLM Inference Test](#123-vllm-inference-test)
    - [1.3 Ollama Deployment on k8s](#13-ollama-deployment-on-k8s)
      - [1.3.1 Ollama Deployment](#131-ollama-deployment)
      - [1.3.2 Ollama Service](#132-ollama-service)
      - [1.3.3 Download and Prepare the Model](#133-download-and-prepare-the-model)
      - [1.3.4 Ollama Inference Tests](#134-ollama-inference-tests)
    - [1.4 OpenWebUI Deployment on k8s](#14-openwebui-deployment-on-k8s)
      - [1.4.1 OpenWebUI Deployment](#141-openwebui-deployment)
      - [1.4.2 OpenWebUI Ingress Service](#142-openwebui-ingress-service)
      - [1.4.3 OpenWebUI Microsoft OAuth2](#143-openwebui-microsoft-oauth2)
  - [2. Model Performance Test](#2-model-performance-test)
    - [2.1 vLLM Local Load Test \& Tokens Per Second Benchmark](#21-vllm-local-load-test--tokens-per-second-benchmark)
  - [vLLM Load Test Results](#vllm-load-test-results)
  - [Acknowledgements](#acknowledgements)
  - [References](#references)


## 1. Installation

### 1.1 Prerequisites

Make sure you meet the following prerequisites before getting started:

- **Kubernetes (k8s)** is installed and running.
- A node with an **NVIDIA GPU** (e.g., L40S, A100, H100) is available.
- **NVIDIA Container Toolkit** is installed.

> 💡 For detailed instructions on setting up **Kubernetes**, **NVIDIA GPU support**, and the **NVIDIA Container Toolkit**, please refer to:  
> [uzunenes/triton-server-hpa](https://github.com/uzunenes/triton-server-hpa)

### 1.2 vLLM Deployment on k8s

#### 1.2.1 vLLM Deployment

Apply the vLLM deployment manifest ([`k8s/vllm-dep.yaml`](k8s/vllm-dep.yaml)):

```bash
kubectl apply -f k8s/vllm-dep.yaml
```

#### 1.2.2 vLLM Service

Apply the vLLM service manifest ([`k8s/vllm-svc-np.yaml`](k8s/vllm-svc-np.yaml)):


``` bash
kubectl apply -f k8s/vllm-svc-np.yaml
```

#### 1.2.3 vLLM Inference Test

To verify your vLLM deployment, send a sample completion request using `curl`.  
First, create a prompt JSON file:

```bash
cat <<EOF > prompt.json
{
  "model": "qwen/qwen3-32B",
  "prompt": "Hello AI, tell me something interesting about deep learning.",
  "temperature": 0.7,
  "max_tokens": 200
}
EOF
```

Then, run the following command (replace <NODE_IP> with your node's IP):


```bash
curl http://<NODE_IP>:30081/generate \
     -H "Content-Type: application/json" \
     -d @prompt.json
```

You should receive a JSON response with the model’s completion, similar to:
```json
{
  "id": "cmpl-12345",
  "object": "text_completion",
  "created": 1716744000,
  "model": "qwen/qwen3-32B",
  "choices": [
    {
      "text": "Deep learning has revolutionized many fields by enabling machines to learn complex patterns from large datasets. For example, deep learning models have achieved state-of-the-art performance in image recognition, natural language processing, and game playing.",
      "index": 0,
      "logprobs": null,
      "finish_reason": "stop"
    }
  ]
}
```

>Note: The actual response format may vary depending on the deployed model and API configuration. 

### 1.3 Ollama Deployment on k8s

#### 1.3.1 Ollama Deployment

Apply the ollama deployment manifest ([`k8s/ollama-dep.yaml`](k8s/ollama-dep.yaml)):

```bash
kubectl apply -f k8s/ollama-dep.yaml
```

#### 1.3.2 Ollama Service

Apply the Ollama service manifest ([`k8s/ollama-svc-np.yaml`](k8s/ollama-svc-np.yaml)):
``` bash
kubectl apply -f k8s/ollama-svc-np.yaml
```

#### 1.3.3 Download and Prepare the Model

Access the Ollama pod and download the desired model.
Example for the nomic-embed-text model:

``` bash
kubectl exec -it <OLLAMA_POD_NAME> -- bash
ollama pull nomic-embed-text
```
Once the model is downloaded, Ollama will be ready to handle requests.

#### 1.3.4 Ollama Inference Tests

1. Embedding Test
Create a JSON file for the embedding request:

``` bash
cat <<EOF > embed.json
{
  "model": "nomic-embed-text",
  "input": "Artificial intelligence is transforming the world."
}
EOF
```

Send the request to the Ollama API:

``` bash
curl http://<NODE_IP>:30134/api/embeddings \
     -H "Content-Type: application/json" \
     -d @embed.json
```

2. Text Generation Test (Example with llama2 model)

Download the text generation model inside the pod:

``` bash
kubectl exec -it <OLLAMA_POD_NAME> -- ollama pull llama2
```

Create a JSON file for the generation request:

``` bash
cat <<EOF > prompt.json
{
  "model": "llama2",
  "prompt": "Explain the concept of reinforcement learning in simple terms."
}
EOF
```

Send the request to the API:

``` bash
curl http://<NODE_IP>:30134/api/generate \
     -H "Content-Type: application/json" \
     -d @prompt.json
```

You should receive the model’s response in JSON format.

💡 Note:

>Replace <NODE_IP> with your Kubernetes worker node’s IP address.
>Replace 30134 with the actual NodePort value from your service manifest.


### 1.4 OpenWebUI Deployment on k8s

#### 1.4.1 OpenWebUI Deployment

Apply the OpenWebUI deployment manifest ([`k8s/openwebui-dep.yaml`](k8s/openwebui-dep.yaml)):

```bash
kubectl apply -f k8s/openwebui-dep.yaml
```

#### 1.4.2 OpenWebUI Ingress Service


```bash
kubectl apply -f k8s/openwebui-ingress.yaml
```

#### 1.4.3 OpenWebUI Microsoft OAuth2

To enable Microsoft OAuth2 authentication in OpenWebUI, follow these steps:

1. **Create an Azure AD Application**
   - Go to [Azure Portal](https://portal.azure.com/).
   - Navigate to **Azure Active Directory** → **App registrations** → **New registration**.
   - Give your application a name (e.g., `OpenWebUI`).
   - In the **Redirect URI** field, add the following:  
     ```
     https://ai.com.tr/auth/callback
     ```
   - Click **Register**.

2. **Obtain Client ID and Tenant ID**
   - After registration, go to your application page.
   - Note down the **Application (client) ID** and **Directory (tenant) ID**.

3. **Create a Client Secret**
   - Go to the **Certificates & secrets** section.
   - Create a new secret using **New client secret**.
   - Save the generated value (it will only be shown once).

4. **Configure OpenWebUI**
   - In your Kubernetes environment, update the OpenWebUI deployment manifest with the following environment variables:

     ```yaml
     env:
       - name: OAUTH_PROVIDER
         value: "microsoft"
       - name: OAUTH_MICROSOFT_CLIENT_ID
         value: "<YOUR_CLIENT_ID>"
       - name: OAUTH_MICROSOFT_CLIENT_SECRET
         value: "<YOUR_CLIENT_SECRET>"
       - name: OAUTH_MICROSOFT_TENANT_ID
         value: "<YOUR_TENANT_ID>"
       - name: OAUTH_REDIRECT_URI
         value: "https://ai.com.tr/auth/callback"
     ```


## 2. Model Performance Test

### 2.1 vLLM Local Load Test & Tokens Per Second Benchmark

``` bash
go run load-test.go
```

## vLLM Load Test Results

```
| Metric | Value |
|--------|-------|
| Total Requests | 20 |
| Total Tokens | 10773 |
| Total Time (s) | 135.46 |
| Tokens per Second | 79.53 |
```

## Acknowledgements
I would like to thank my teammates for their valuable support during this work.

- Zekeriye Altunkaynak
- Alper Aldemir
- Ceyhun Erdil
- Güneş Balçık

## References
- [NVIDIA L40S Data Center GPU](https://www.nvidia.com/en-us/data-center/l40s/)
- [vLLM Documentation](https://docs.vllm.ai/en/stable/)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Ollama Documentation](https://github.com/ollama/ollama/tree/main/docs)
- [OpenWebUI Documentation](https://docs.openwebui.com/)
