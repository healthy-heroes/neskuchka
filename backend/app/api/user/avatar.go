package api_user

import (
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/healthy-heroes/neskuchka/backend/app/api/httpx"
	"github.com/healthy-heroes/neskuchka/backend/app/domain"
	"github.com/healthy-heroes/neskuchka/backend/app/internal/session"
)

const maxAvatarSize = 512 * 1024 // 512 KB

var allowedMimeTypes = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

func (s *Service) Avatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	avatar, err := s.avatarStorage.Get(r.Context(), id)
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to get avatar")
		return
	}

	w.Header().Set("Content-Type", avatar.MimeType)
	w.Header().Set("Content-Length", strconv.Itoa(len(avatar.Data)))
	w.Header().Set("Cache-Control", "private, max-age=86400")
	w.Write(avatar.Data)
}

func (s *Service) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	id := domain.UserID(session.MustGetUserID(r))

	r.Body = http.MaxBytesReader(w, r.Body, maxAvatarSize)

	err := r.ParseMultipartForm(maxAvatarSize)
	if err != nil {
		httpx.RenderError(w, s.logger, http.StatusBadRequest, err, "file is too large")
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

	avatar := domain.Avatar{
		MimeType: mimeType,
		Data:     data,
	}

	err = s.avatarStorage.Save(r.Context(), id, avatar)
	if err != nil {
		httpx.RenderDomainError(w, s.logger, err, "failed to save avatar")
		return
	}

	httpx.Render(w, nil)
}
