package rewritemanager

import (
	"context"
	"fmt"
	"log"
	"time"

	aiAssistInterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	aiAssistModel "github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	dbInterface "github.com/ray31245/seo_cluster/pkg/db/db_interface"
)

const (
	retryLimit = 50
	retryDelay = 100 * time.Millisecond

	defaultSystemPromptKey = "system_prompt"
	defaultPromptKey       = "prompt"

	defaultExtendSystemPromptKey = "extend_system_prompt"
	defaultExtendPromptKey       = "extend_prompt"
)

type RewriteManager struct {
	aiAssist  aiAssistInterface.AIAssistInterface
	configDAO dbInterface.KVConfigDAOInterface
}

func NewRewriteManager(aiAssist aiAssistInterface.AIAssistInterface, configDAO dbInterface.KVConfigDAOInterface) *RewriteManager {
	return &RewriteManager{
		aiAssist:  aiAssist,
		configDAO: configDAO,
	}
}

func (r *RewriteManager) Rewrite(ctx context.Context, text []byte) (aiAssistModel.RewriteResponse, error) {
	res, err := r.aiAssist.Rewrite(ctx, text)
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.Rewrite: %w", err)
	}

	return res, nil
}

func (r *RewriteManager) ExtendRewrite(ctx context.Context, text []byte) (aiAssistModel.ExtendRewriteResponse, error) {
	res, err := r.aiAssist.ExtendRewrite(ctx, text)
	if err != nil {
		return aiAssistModel.ExtendRewriteResponse{}, fmt.Errorf("RewriteManager.ExtendRewrite: %w", err)
	}

	return res, nil
}

func (r *RewriteManager) CustomRewrite(ctx context.Context, systemPrompt string, prompt string, content []byte) (aiAssistModel.RewriteResponse, error) {
	res, err := r.aiAssist.CustomRewrite(ctx, systemPrompt, prompt, content)
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.CustomRewrite: %w", err)
	}

	return res, nil
}

func (r *RewriteManager) DefaultRewrite(ctx context.Context, text []byte) (aiAssistModel.RewriteResponse, error) {
	systemPrompt, err := r.GetDefaultSystemPrompt()
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultRewrite: %w", err)
	}

	prompt, err := r.GetDefaultPrompt()
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultRewrite: %w", err)
	}

	return r.CustomRewrite(ctx, systemPrompt, prompt, text)
}

func (r *RewriteManager) DefaultExtendRewrite(ctx context.Context, text []byte) (aiAssistModel.RewriteResponse, error) {
	systemPrompt, err := r.GetDefaultExtendSystemPrompt()
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultExtendRewrite: %w", err)
	}

	prompt, err := r.GetDefaultPrompt()
	if err != nil {
		return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultExtendRewrite: %w", err)
	}

	return r.aiAssist.CustomRewrite(ctx, systemPrompt, prompt, text)
}

func (r *RewriteManager) RewriteUntil(ctx context.Context, text []byte) (res aiAssistModel.RewriteResponse, err error) {
	log.Println("rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		res, err = r.aiAssist.Rewrite(ctx, text)
		if err == nil {
			return res, nil
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.RewriteUntil: %w", err)
}

func (r *RewriteManager) ExtendRewriteUntil(ctx context.Context, text []byte) (res aiAssistModel.ExtendRewriteResponse, err error) {
	log.Println("extending rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		res, err = r.aiAssist.ExtendRewrite(ctx, text)
		if err == nil {
			return res, nil
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	return aiAssistModel.ExtendRewriteResponse{}, fmt.Errorf("RewriteManager.ExtendRewriteUntil: %w", err)
}

func (r *RewriteManager) CustomRewriteUntil(ctx context.Context, systemPrompt string, prompt string, content []byte) (res aiAssistModel.RewriteResponse, err error) {
	log.Println("custom rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		res, err = r.aiAssist.CustomRewrite(ctx, systemPrompt, prompt, content)
		if err == nil {
			return res, nil
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.CustomRewriteUntil: %w", err)
}

func (r *RewriteManager) DefaultRewriteUntil(ctx context.Context, text []byte) (res aiAssistModel.RewriteResponse, err error) {
	log.Println("default rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		res, err = r.DefaultRewrite(ctx, text)
		if err == nil {
			return res, nil
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultRewriteUntil: %w", err)
}

func (r *RewriteManager) DefaultExtendRewriteUntil(ctx context.Context, text []byte) (res aiAssistModel.RewriteResponse, err error) {
	log.Println("default extending rewriting...")

	r.aiAssist.Lock()
	defer r.aiAssist.Unlock()

	for range retryLimit {
		res, err = r.DefaultExtendRewrite(ctx, text)
		if err == nil {
			return res, nil
		}

		log.Println("retrying...")
		<-time.After(retryDelay)
	}

	return aiAssistModel.RewriteResponse{}, fmt.Errorf("RewriteManager.DefaultExtendRewriteUntil: %w", err)
}

func (r *RewriteManager) GetDefaultSystemPrompt() (string, error) {
	systemPrompt, err := r.configDAO.GetByKey(defaultSystemPromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultSystemPrompt: %w", err)
	}

	return systemPrompt.Value, nil
}

func (r *RewriteManager) SetDefaultSystemPrompt(systemPrompt string) error {
	err := r.configDAO.UpsertByKey(defaultSystemPromptKey, systemPrompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultSystemPrompt: %w", err)
	}

	return nil
}

func (r *RewriteManager) GetDefaultPrompt() (string, error) {
	prompt, err := r.configDAO.GetByKey(defaultPromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultPrompt: %w", err)
	}

	return prompt.Value, nil
}

func (r *RewriteManager) SetDefaultPrompt(prompt string) error {
	err := r.configDAO.UpsertByKey(defaultPromptKey, prompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultPrompt: %w", err)
	}

	return nil
}

func (r *RewriteManager) GetDefaultExtendSystemPrompt() (string, error) {
	systemPrompt, err := r.configDAO.GetByKey(defaultExtendSystemPromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultExtendSystemPrompt: %w", err)
	}

	return systemPrompt.Value, nil
}

func (r *RewriteManager) SetDefaultExtendSystemPrompt(systemPrompt string) error {
	err := r.configDAO.UpsertByKey(defaultExtendSystemPromptKey, systemPrompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultExtendSystemPrompt: %w", err)
	}

	return nil
}

func (r *RewriteManager) GetDefaultExtendPrompt() (string, error) {
	prompt, err := r.configDAO.GetByKey(defaultExtendPromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultExtendPrompt: %w", err)
	}

	return prompt.Value, nil
}

func (r *RewriteManager) SetDefaultExtendPrompt(prompt string) error {
	err := r.configDAO.UpsertByKey(defaultExtendPromptKey, prompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultExtendPrompt: %w", err)
	}

	return nil
}
