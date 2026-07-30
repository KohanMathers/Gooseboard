package store

import "gooseboard/internal/schema"

type Store interface {
	Migrate(entities map[string]schema.Entity) error
	Create(entity string, data map[string]any) (id string, err error)
	Get(entity string, id string) (map[string]any, error)
	List(entity string, filter map[string]any) ([]map[string]any, error)
	Update(entity string, id string, data map[string]any) error
	Delete(entity string, id string) error
}
