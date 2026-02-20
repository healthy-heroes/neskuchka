package email

import (
	emailpkg "github.com/go-pkgz/email"
	"github.com/rs/zerolog"

	lgr "github.com/healthy-heroes/neskuchka/backend/app/internal/logger"
)

type Service struct {
	sender *emailpkg.Sender

	from   string
	logger zerolog.Logger
}

type Opts struct {
	Host string
	Port int
	From string

	Logger zerolog.Logger
}

func NewService(opts Opts) *Service {
	logger := opts.Logger.With().Str("pkg", "email").Logger()

	from := opts.From
	if from == "" {
		from = "noreply@neskuchka"
	}

	return &Service{
		from: from,
		sender: emailpkg.NewSender(
			opts.Host,
			emailpkg.ContentType("text/plain"),
			emailpkg.Port(opts.Port),
			emailpkg.Log(lgr.NewSimple(logger)),
		),

		logger: opts.Logger.With().Str("pkg", "email").Logger(),
	}
}

func (s *Service) Send(to, subject, text string) error {
	return s.sender.Send(text, emailpkg.Params{
		From:    s.from,
		To:      []string{to},
		Subject: subject,
	})
}
