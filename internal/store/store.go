package store

// Store агрегирует все хранилища данных приложения.
type Store struct {
	News        *NewsStore
	Exhibitions *ExhibitionStore
	Exhibits    *ExhibitStore
}

func New() *Store {
	return &Store{
		News:        NewNewsStore(),
		Exhibitions: NewExhibitionStore(),
		Exhibits:    NewExhibitStore(),
	}
}
