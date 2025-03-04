package aiassist

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/generative-ai-go/genai"
	aiassistinterface "github.com/ray31245/seo_cluster/pkg/ai_assist/ai_assist_interface"
	"github.com/ray31245/seo_cluster/pkg/ai_assist/model"
	"google.golang.org/api/option"
)

var _ aiassistinterface.AIAssistInterface = &AIAssist{}

const (
	maxTemperature = 2
)

type AIAssist struct {
	token          string
	client         *genai.Client
	rewriter       *genai.GenerativeModel
	extendRewriter *genai.GenerativeModel
	commenter      *genai.GenerativeModel
	evaluator      *genai.GenerativeModel
	keyWordFinder  *genai.GenerativeModel
	keyWordMatcher *genai.GenerativeModel
	categorySelect *genai.GenerativeModel
	lock           sync.Mutex
	isLimitedUsage bool
}

func NewAIAssist(ctx context.Context, token string, isLimitedUsuage bool) (*AIAssist, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, fmt.Errorf("failed to create new client: %w", err)
	}

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	rewriter := client.GenerativeModel("gemini-2.0-flash")
	rewriter.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	extendRewriter := client.GenerativeModel("gemini-2.0-flash")
	extendRewriter.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	commenter := client.GenerativeModel("gemini-2.0-flash")
	commenter.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		Temperature: func() *float32 {
			temp := float32(maxTemperature)

			return &temp
		}(),
	}

	evaluator := client.GenerativeModel("gemini-2.0-flash")
	evaluator.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
	}

	keyWordFinder := client.GenerativeModel("gemini-1.5-pro")
	keyWordFinder.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"KeyWords": {
					Type: genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeString,
					},
				},
			},
		},
	}

	keyWordMatcher := client.GenerativeModel("gemini-2.0-flash")
	keyWordMatcher.GenerationConfig = genai.GenerationConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeString,
			},
		},
	}

	categorySelect := client.GenerativeModel("gemini-2.0-flash")
	categorySelect.ResponseMIMEType = "application/json"
	categorySelect.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"id": {
				Type: genai.TypeString,
			},
			"isFind": {
				Type: genai.TypeBoolean,
			},
		},
	}

	return &AIAssist{
		token:          token,
		client:         client,
		rewriter:       rewriter,
		extendRewriter: extendRewriter,
		commenter:      commenter,
		evaluator:      evaluator,
		keyWordFinder:  keyWordFinder,
		keyWordMatcher: keyWordMatcher,
		categorySelect: categorySelect,
		isLimitedUsage: isLimitedUsuage,
	}, nil
}

func (a *AIAssist) CustomRewrite(ctx context.Context, systemPrompt string, prompt string, content []byte) (string, error) {
	customRewriter := a.client.GenerativeModel("gemini-2.0-flash")

	customRewriter.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemPrompt)}}

	resp, err := customRewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, content)))
	if err != nil {
		return "", fmt.Errorf("failed to rewrite content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", errors.New("no content generated")
	}

	res := fmt.Sprint(resp.Candidates[0].Content.Parts[0])

	return res, nil
}

func (a *AIAssist) Close() error {
	err := a.client.Close()
	if err != nil {
		return fmt.Errorf("failed to close client: %w", err)
	}

	return nil
}

func (a *AIAssist) Rewrite(ctx context.Context, text []byte) (model.RewriteResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位收悉区块链的简体中文语系专栏作家，请你将以下内容用你的话重新阐述文章中的内容，并且以你认为没有AI痕迹的方式表达，并订一个标题。请使用json格式输出：{Title: string,Content: string}"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("failed to rewrite content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.RewriteResponse{}, errors.New("no content generated")
	}

	res := model.RewriteResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.RewriteResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) ExtendRewrite(ctx context.Context, text []byte) (model.ExtendRewriteResponse, error) {
	promt := "你是一位收悉区块链的简体中文语系专栏作家，你需要对这篇文章进行扩展，并且以你认为没有AI痕迹的方式表达。请使用json格式输出：{Title: string,Content: string}"

	resp, err := a.extendRewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", promt, text)))
	if err != nil {
		return model.ExtendRewriteResponse{}, fmt.Errorf("failed to extend rewrite content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.ExtendRewriteResponse{}, errors.New("no content generated")
	}

	res := model.ExtendRewriteResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.ExtendRewriteResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) Comment(ctx context.Context, text []byte) (model.CommentResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位简体中文语系读者,你在网路上看到以下文章，请随性且简洁地在这篇文章下留言。并且以一位看新闻的人的角度记录这个文章能够为你提供的价值。最低0分滿分100分。请使用json格式输出：{Comment: string, Score: int}\n。"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.CommentResponse{}, fmt.Errorf("failed to comment content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.CommentResponse{}, errors.New("no content generated")
	}

	res := model.CommentResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.CommentResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) Evaluate(ctx context.Context, text []byte) (model.EvaluateResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	prompt := "你是一位区块链专栏作家，你的文章被编辑修改过，你需要评价这篇文章的质量。最低0分滿分100分。请使用json格式输出：{ Score: int}。"

	resp, err := a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.EvaluateResponse{}, fmt.Errorf("failed to evaluate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.EvaluateResponse{}, errors.New("no content generated")
	}

	res := model.EvaluateResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.EvaluateResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) SelectCategory(ctx context.Context, req model.SelectCategoryRequest) (model.SelectCategoryResponse, error) {
	optionStr, err := json.Marshal(req.CategoriesOption)
	if err != nil {
		return model.SelectCategoryResponse{}, fmt.Errorf("failed to marshal categories option: %w", err)
	}
	prompts := []genai.Part{
		genai.Text("請幫我對以下文章根據提供的選項選擇最適合的分類，如果有多個分類選擇順序最優先的分類，例如選項（币安交易所、加密货币交易所、火币交易所、加密货币再质押）那順序最優先的選項就是\"币安交易所\"。"),
		genai.Text("完整的例子如下："),
		genai.Text("文章：<p>本文分析了埃隆·马斯克在2024年美国大选中对唐纳德·特朗普胜选所起到的关键作用，以及这种合作关系对美国政治和社会可能产生的深远影响。马斯克利用其在科技、商业和社交媒体上的巨大影响力，为特朗普的竞选提供了资金、地面组织和宣传支持，吸引了大量对特朗普的政策感兴趣但对其性格感到厌倦的年轻选民。</p>\n\n<p>文章指出，马斯克与特朗普的合作关系并非完全一致，两人的个性都强势，未来可能出现冲突。马斯克的动机也值得深思，其最终目标可能是为了实现其在太空探索方面的宏伟计划，而将美国政府作为实现这一目标的工具。</p>\n\n<p>马斯克被特朗普任命领导一个名为“政府效率部 (DOGE)”的新部门，其目标是精简联邦政府机构，削减开支。文章质疑了这一目标的可行性以及对社会福利计划可能造成的负面影响，特别是对弱势群体的冲击。马斯克的“效率”承诺可能导致医疗、教育等公共服务的削减，对依赖政府支持的民众造成严重后果。</p>\n\n<p>文章还探讨了马斯克与特朗普合作关系中潜在的利益冲突，以及马斯克在影响政府监管机构方面可能扮演的角色。马斯克旗下公司特斯拉和SpaceX此前都与政府监管机构发生过冲突，而“政府效率部”的成立可能使其更容易影响甚至规避这些监管。</p>\n\n<p>文章最后指出，马斯克的政治立场并不明确，其与特朗普的关系也曾经历过动荡。马斯克所宣扬的“言论自由”理念与其在Twitter上的行为存在矛盾，其对“工作思维病毒”的批评也反映出其某种意识形态立场。文章警告说，马斯克和特朗普的合作可能导致一种“寡头政治”的局面，对美国的民主和社会公平造成威胁，普通民众可能无法从这种权力结合中获益，反而可能遭受损失。</p>\n<img src=\"https://img.jinse.cn/7324717_watermarknone.png\"/><img src=\"https://img.jinse.cn/7324724_watermarknone.png\"/><img src=\"https://img.jinse.cn/7324738_watermarknone.png\"/><img src=\"https://img.jinse.cn/7324746_watermarknone.png\"/>"),
		genai.Text("選項：[{name:DOGE,id:ba9f1ec9-bfd7-405f-967b-45449011fbe5},{name:马斯克,id:22f585dd-c50d-4d8f-93c6-534eb89682c7},{name:時事,id:c32084ff-b335-4c60-8804-8fa9d54cf216},{name:美食,id:dcdc2fe2-6a1b-4231-9179-48ee3520de2a},{name:美国政府,id:48fb04ce-1834-474b-94e6-addaaedf6dc2}]，這些選項中的id為ba9f1ec9-bfd7-405f-967b-45449011fbe5或22f585dd-c50d-4d8f-93c6-534eb89682c7或c32084ff-b335-4c60-8804-8fa9d54cf216或48fb04ce-1834-474b-94e6-addaaedf6dc2都符合這個文章的分類，但是順序最優先的分類為ba9f1ec9-bfd7-405f-967b-45449011fbe5。那回傳結果會是：{id:\"ba9f1ec9-bfd7-405f-967b-45449011fbe5\",isFind:true}"),
		genai.Text("另外如果選項中找不到符合的分類則回傳空字串，以上面同一例子來說，選項為[{name:沙丁魚,id:ba9f1ec9-bfd7-405f-967b-45449011fbe5},{name:園藝景觀,id:22f585dd-c50d-4d8f-93c6-534eb89682c7},{name:咖啡,id:c32084ff-b335-4c60-8804-8fa9d54cf216},{name:美食,id:dcdc2fe2-6a1b-4231-9179-48ee3520de2a},{name:宮崎英高,id:48fb04ce-1834-474b-94e6-addaaedf6dc2}]，這些選項都不符合此文章的分類則回傳空字串。那回傳結果會是：{id:\"\",isFind:false}"),
		genai.Text("請注意，以上的文章和選項只是範例，請勿使用述的選項及文章做爲回答。"),
		genai.Text(fmt.Sprintf("文章：%s\n", req.Text)),
		genai.Text("選項："),
		genai.Text(fmt.Sprintf("%s\n", optionStr)),
	}

	resp, err := a.categorySelect.GenerateContent(ctx, prompts...)
	if err != nil {
		return model.SelectCategoryResponse{}, fmt.Errorf("failed to select category: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.SelectCategoryResponse{}, errors.New("no content generated")
	}

	res := model.SelectCategoryResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.SelectCategoryResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) Lock() {
	if a.isLimitedUsage {
		a.lock.Lock()
	}
}

func (a *AIAssist) TryLock() bool {
	if a.isLimitedUsage {
		return a.lock.TryLock()
	}
	return true
}

func (a *AIAssist) Unlock() {
	if a.isLimitedUsage {
		a.lock.Unlock()
	}
}

func (a *AIAssist) FindKeyWords(ctx context.Context, text []byte) (model.FindKeyWordsResponse, error) {
	//nolint:gosmopolitan // prompt is a string
	// prompt := "你是一位区块链专栏作家并且擅长seo，你需要列出这篇文章中与区块链、数位货币、投资相关的中文长尾关键字。请使用json格式输出：{KeyWords: []string}。"
	prompt := "你是一位区块链专栏作家并且擅长seo，你需要列出这篇文章中与区块链、数位货币、投资相关的中文核心关键字。请使用json格式输出：{KeyWords: []string}。"

	resp, err := a.keyWordFinder.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	if err != nil {
		return model.FindKeyWordsResponse{}, fmt.Errorf("failed to find keywords: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.FindKeyWordsResponse{}, errors.New("no content generated")
	}

	//nolint:gosmopolitan // prompt is a string
	// prompt = "请你将以下的词语中长度少于4个字的词语改写成长度4到5个字的同义词。请使用json格式输出：{KeyWords: []string}。"

	// resp, err = a.rewriter.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, text)))
	// if err != nil {
	// 	return model.FindKeyWordsResponse{}, fmt.Errorf("failed to filter keywords: %w", err)
	// }

	// if len(resp.Candidates) == 0 {
	// 	return model.FindKeyWordsResponse{}, errors.New("no content generated")
	// }

	res := model.FindKeyWordsResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.FindKeyWordsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) MatchKeyWords(ctx context.Context, text []byte, keywords []string) (model.MatchKeyWordsResponse, error) {
	optionStr, err := json.Marshal(keywords)
	if err != nil {
		return model.MatchKeyWordsResponse{}, fmt.Errorf("failed to marshal keywords: %w", err)
	}

	prompt := fmt.Sprintf("以下将提供文章及标签，请从标签的阵列中找出最多五个适合此文章的，如果找不到就回传空阵列\n标签：\n%s\n文章：%s", optionStr, text)

	resp, err := a.keyWordMatcher.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return model.MatchKeyWordsResponse{}, fmt.Errorf("failed to match keywords: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return model.MatchKeyWordsResponse{}, errors.New("no content generated")
	}

	res := model.MatchKeyWordsResponse{}

	err = json.Unmarshal([]byte(fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0])), &res)
	if err != nil {
		return model.MatchKeyWordsResponse{}, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return res, nil
}

func (a *AIAssist) MakeTitle(ctx context.Context, systemPrompt string, prompt string, content []byte) (string, error) {
	model := a.client.GenerativeModel("gemini-2.0-flash")

	model.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemPrompt)}}

	resp, err := model.GenerateContent(ctx, genai.Text(fmt.Sprintf("%s\n%s", prompt, content)))
	if err != nil {
		return "", fmt.Errorf("failed to make title: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", errors.New("no content generated")
	}

	return fmt.Sprintf("%s", resp.Candidates[0].Content.Parts[0]), nil
}
