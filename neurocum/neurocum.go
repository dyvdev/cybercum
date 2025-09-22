package neurocum

import (
	"log/slog"
	"strings"

	gigachat "github.com/saintbyte/gigachat_api"
)

func Respond(context []string, prompt string) string {
	if len(context) == 0 || prompt == "" {
		return ""
	}
	chat := &gigachat.Gigachat{
		ApiHost:           gigachat.GigaChatApiHost,
		RepetitionPenalty: 1,
		TopP:              1.0,
		Model:             "GigaChat-Preview",
		MaxTokens:         100,
		Temperature:       1,
		AuthData:          "",
	}
	answer, err := ask(chat, gigachat.GigaChatRoleSystem, prompt)
	if err != nil {
		slog.Error("Ask error:", err)
		return ""
	}
	in := strings.Join(context, "\n")
	answer, err = ask(chat, gigachat.GigaChatRoleUser, in)
	if err != nil {
		slog.Error("Ask error:", err)
		return ""
	}
	return answer
}

func ask(g *gigachat.Gigachat, role, input string) (string, error) {
	return g.ChatCompletions([]gigachat.MessageRequest{
		{
			Role:    role,
			Content: input,
		},
	})
}
func CheckConnect() bool {
	chat := &gigachat.Gigachat{
		ApiHost:           gigachat.GigaChatApiHost,
		RepetitionPenalty: 1,
		TopP:              1.0,
		Model:             "GigaChat-Preview",
		MaxTokens:         100,
		Temperature:       1,
		AuthData:          "",
	}
	_, err := chat.GetModels()
	if err != nil {
		slog.Error("CheckConnect error:", err)
	}
	return err == nil
}
