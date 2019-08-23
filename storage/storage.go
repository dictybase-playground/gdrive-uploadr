package storage

import (
	"github.com/dictybase-playground/gdrive-uploadr/schema"
)

type Storage interface {
	Save(key, *schema.ImageAttributes) error
	SaveAll(key, ...*schema.ImageAttributes) error
	GetAll(key, int, int, int) ([]*schema.ImageAttributes, error)
}
