package api_user

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/avatarimg"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

// maxAvatarSize caps the upload. A photo straight from a phone is several
// megabytes, and since the server downscales it anyway, rejecting it for its
// weight alone only annoys the user.
const maxAvatarSize = 8 * 1024 * 1024

var allowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

func (s *Service) avatar(w http.ResponseWriter, r *http.Request, id domain.UserID) {
	avatar, err := s.avatarStore.Get(r.Context(), id)
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to get avatar")
		return
	}

	w.Header().Set("Content-Type", avatar.MimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar.Data)))
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Write(avatar.Data)
}

func (s *Service) MyAvatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	s.avatar(w, r, id)
}

func (s *Service) UserAvatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(chi.URLParam(r, "id"))

	s.avatar(w, r, id)
}

func (s *Service) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	err := s.avatarStore.Delete(r.Context(), id)
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to delete avatar")
		return
	}

	httpx.Render(w, nil)
}

func (s *Service) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize)

	err := r.ParseMultipartForm(maxAvatarSize)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			httpx.RenderError(w, s.logger, http.StatusRequestEntityTooLarge, err, "file is too large")
			return
		}

		httpx.RenderError(w, s.logger, http.StatusBadRequest, err, "failed parse multipart")
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusBadRequest, err, "missing file")
		return
	}
	defer file.Close()

	mimeType := header.Header.Get("Content-Type")
	if !slices.Contains(allowedMimeTypes, mimeType) {
		httpx.RenderError(w, s.logger, http.StatusBadRequest,
			fmt.Errorf("file type %s not allowed", mimeType),
			"unsupported file type",
		)
		return
	}

	data, err := io.ReadAll(file)
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusInternalServerError, err, "failed to read file")
		return
	}

	mimeType = http.DetectContentType(data)
	if !slices.Contains(allowedMimeTypes, mimeType) {
		httpx.RenderError(w, s.logger, http.StatusBadRequest,
			fmt.Errorf("file type %s not allowed", mimeType),
			"file content doesn't match allowed image types",
		)
		return
	}

	normalized, err := avatarimg.Normalize(data)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, avatarimg.ErrUnsupportedFormat) || errors.Is(err, avatarimg.ErrTooLarge) {
			status = http.StatusBadRequest
		}

		httpx.RenderError(w, s.logger, status, err, "failed to process image")
		return
	}

	avatar := domain.Avatar{
		MimeType: avatarimg.MimeType,
		Data:     normalized,
	}

	err = s.avatarStore.Save(r.Context(), id, avatar)
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to save avatar")
		return
	}

	httpx.Render(w, nil)
}
