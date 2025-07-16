package authproviders

import (
	"fmt"

	"github.com/go-pkgz/auth/v2"
	"github.com/go-pkgz/auth/v2/avatar"
	"github.com/go-pkgz/auth/v2/logger"
	"github.com/go-pkgz/auth/v2/provider"
	"github.com/go-pkgz/auth/v2/token"
)

// Service extends auth.Service for adding custom providers
type Service struct {
	opts auth.Opts

	logger      logger.L
	jwtService  *token.Service
	avatarProxy *avatar.Proxy
	issuer      string
	useGravatar bool
}

// NewService initializes Service with defaults values
func NewService(opts auth.Opts) *Service {
	s := &Service{
		opts:        opts,
		logger:      opts.Logger,
		issuer:      opts.Issuer,
		useGravatar: opts.UseGravatar,
	}

	if opts.Issuer == "" {
		s.issuer = "Neskuchka"
	}

	if opts.Logger == nil {
		s.logger = logger.NoOp
	}

	jwtService := token.NewService(token.Opts{
		SecretReader:      opts.SecretReader,
		ClaimsUpd:         opts.ClaimsUpd,
		SecureCookies:     opts.SecureCookies,
		TokenDuration:     opts.TokenDuration,
		CookieDuration:    opts.CookieDuration,
		DisableXSRF:       opts.DisableXSRF,
		DisableIAT:        opts.DisableIAT,
		JWTCookieName:     opts.JWTCookieName,
		JWTCookieDomain:   opts.JWTCookieDomain,
		JWTHeaderKey:      opts.JWTHeaderKey,
		XSRFCookieName:    opts.XSRFCookieName,
		XSRFHeaderKey:     opts.XSRFHeaderKey,
		XSRFIgnoreMethods: opts.XSRFIgnoreMethods,
		SendJWTHeader:     opts.SendJWTHeader,
		JWTQuery:          opts.JWTQuery,
		Issuer:            s.issuer,
		AudienceReader:    opts.AudienceReader,
		AudSecrets:        opts.AudSecrets,
		SameSite:          opts.SameSiteCookie,
	})

	if opts.SecretReader == nil {
		jwtService.SecretReader = token.SecretFunc(func(string) (string, error) {
			return "", fmt.Errorf("secrets reader not available")
		})
		s.logger.Logf("[WARN] no secret reader defined")
	}

	s.jwtService = jwtService

	if opts.AvatarStore != nil {
		s.avatarProxy = &avatar.Proxy{
			Store:       opts.AvatarStore,
			URL:         opts.URL,
			RoutePath:   opts.AvatarRoutePath,
			ResizeLimit: opts.AvatarResizeLimit,
			L:           s.logger,
		}
		if s.avatarProxy.RoutePath == "" {
			s.avatarProxy.RoutePath = "/avatar"
		}
	}

	return s
}

func (s *Service) NewVerifyProvider(name string, msgTemplate string, sender provider.Sender) *VerifyHandler {
	return &VerifyHandler{
		L:            s.logger,
		ProviderName: name,
		Issuer:       s.issuer,
		TokenService: s.jwtService,
		AvatarSaver:  s.avatarProxy,
		Sender:       sender,
		Template:     msgTemplate,
		UseGravatar:  s.useGravatar,
	}
}
