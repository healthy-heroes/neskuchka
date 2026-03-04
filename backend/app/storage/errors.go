package storage

import (
	"database/sql"
	"errors"

	"github.com/healthy-heroes/neskuchka/backend/app/domain"
)

func HandleSqlError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}

	return err
}
