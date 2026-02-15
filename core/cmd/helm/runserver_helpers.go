package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Mindburn-Labs/helm/core/pkg/kernel"
	"github.com/Mindburn-Labs/helm/core/pkg/llm"
	"github.com/Mindburn-Labs/helm/core/pkg/store"
)

func initLimiterStoreFromEnv() kernel.LimiterStore {
	if redisAddr := os.Getenv("REDIS_ADDR"); redisAddr != "" {
		log.Printf("[kernel] rate limiter: redis at %s", redisAddr)
		return kernel.NewRedisLimiterStore(redisAddr, "", 0)
	}

	log.Printf("[kernel] rate limiter: in-memory")
	return kernel.NewInMemoryLimiterStore()
}

func initEmbedderAndModels(openAIKey string) (store.Embedder, llm.Client, llm.Client, error) {
	if openAIKey == "" {
		return nil, nil, nil, fmt.Errorf("fail-closed: OPENAI_API_KEY not set")
	}

	log.Printf("[kernel] embedder: openai")
	embedder := store.NewOpenAIEmbedder(openAIKey)

	// Fast = GPT-4o-mini (default)
	fastModelID := os.Getenv("HELM_LLM_FAST_MODEL")
	if fastModelID == "" {
		fastModelID = "gpt-4o-mini"
	}
	fastModel := llm.NewOpenAIClient(openAIKey, fastModelID)

	// Smart = GPT-4o (default)
	smartModelID := os.Getenv("HELM_LLM_SMART_MODEL")
	if smartModelID == "" {
		smartModelID = "gpt-4o"
	}
	smartModel := llm.NewOpenAIClient(openAIKey, smartModelID)

	return embedder, fastModel, smartModel, nil
}
