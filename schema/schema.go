package schema

import (
	"time"
)

// ImageAttributes contains the storage paths for original and thumbnail images
// along with additional meta information
type ImageAttibutes struct {
	CheckSum    string    `json:"checksum"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	CreatedOn   time.Time `json:"created_on"`
	Original    string    `json:"original"`
	Thumbnail   string    `json:"thumbnail"`
}
