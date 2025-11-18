package bot

import (
	"github.com/tuor4eg/vsratoved/internal/llm"
)

// StartMessage returns the welcome message text for /start command
func StartMessage() string {
	return "Привет, я Всратовед. Я выдаю всратые псевдо-мотивационные советы.\nИспользуй /vsrata для мягкого режима и /vsrata_spicy для более жёсткого."
}

// AdviceMessage returns a message text with advice formatted as a quote with author
func AdviceMessage(response *llm.AdviceResponse) string {
	if response.Author != "" && response.Advice != "" {
		return "💡 «" + response.Advice + "»\n\n— " + response.Author
	}
	// Fallback if parsing failed
	if response.Advice != "" {
		return "💡 " + response.Advice
	}
	return ""
}

// ErrorMessage returns a message text with error fallback
func ErrorMessage() string {
	msg := llm.ErrorFallbackMessage()
	author := "Всратовед"

	return AdviceMessage(&llm.AdviceResponse{
		Author: author,
		Advice: msg,
	})
}

// UnknownCommandMessage returns a message text for unknown command
func UnknownCommandMessage() string {
	return "❓ Неизвестная команда. Используй /start для списка доступных команд."
}
