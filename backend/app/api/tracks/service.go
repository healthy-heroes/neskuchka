package tracks

import (
	"github.com/healthy-heroes/neskuchka/backend/app/store/datastore"
)

type Service struct {
	store *datastore.DataStore
}

func NewService(store *datastore.DataStore) *Service {
	return &Service{store}
}
