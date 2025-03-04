package rewritemanager

import (
	"fmt"
)

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

func (r *RewriteManager) GetDefaultMakeTitleSystemPrompt() (string, error) {
	systemPrompt, err := r.configDAO.GetByKey(defaultMakeTitleSystemPromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultMakeTitleSystemPrompt: %w", err)
	}

	return systemPrompt.Value, nil
}

func (r *RewriteManager) SetDefaultMakeTitleSystemPrompt(systemPrompt string) error {
	err := r.configDAO.UpsertByKey(defaultMakeTitleSystemPromptKey, systemPrompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultMakeTitleSystemPrompt: %w", err)
	}

	return nil
}

func (r *RewriteManager) GetDefaultMakeTitlePrompt() (string, error) {
	prompt, err := r.configDAO.GetByKey(defaultMakeTitlePromptKey)
	if err != nil {
		return "", fmt.Errorf("RewriteManager.GetDefaultMakeTitlePrompt: %w", err)
	}

	return prompt.Value, nil
}

func (r *RewriteManager) SetDefaultMakeTitlePrompt(prompt string) error {
	err := r.configDAO.UpsertByKey(defaultMakeTitlePromptKey, prompt)
	if err != nil {
		return fmt.Errorf("RewriteManager.SetDefaultMakeTitlePrompt: %w", err)
	}

	return nil
}
