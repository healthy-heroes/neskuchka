package email

import (
	"fmt"
	"strings"

	"github.com/matcornic/hermes/v2"
)

type Templater struct {
	baseUrl string
	hermes  *hermes.Hermes
}

func NewTemplate(baseUrl string) *Templater {
	return &Templater{
		baseUrl: strings.TrimRight(baseUrl, "/"),
		hermes: &hermes.Hermes{
			Product: hermes.Product{
				Name: "Нескучка",
			},
		},
	}
}

func (t Templater) AuthLink(token string) (string, error) {
	text, err := t.hermes.GeneratePlainText(hermes.Email{
		Body: hermes.Body{
			Title: "Подтверждение входа",
			Intros: []string{
				"Кто-то запросил вход в ваш аккаунт.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Нажмите кнопку для подтверждения:",
					Button: hermes.Button{
						Color: "#22BC66",
						Text:  "Подтвердить вход",
						Link:  fmt.Sprintf("%s/login/confirm?token=%s", t.baseUrl, token),
					},
				},
			},
			Outros: []string{
				"Если вы не запрашивали вход — просто проигнорируйте письмо.",
			},
		},
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate plain text: %w", err)
	}

	return text, nil
}
