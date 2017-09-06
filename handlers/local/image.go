package local

import "github.com/dictybase-playground/gdrive-uploadr/apihelpers/apherror"
import "github.com/go-chi/render"
import "github.com/pkg/errors"
import "github.com/sirupsen/logrus"
import "io"
import "net/http"
import "os"
import "path/filepath"

// ImageHandler is an net/http handler for managing images
type ImageHandler struct {
	Logger    *logrus.Logger
	Key       string
	LocalPath string
}

// Create is a http.Handler method for handling POST requests
func (img *ImageHandler) Create(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile(img.Key)
	if err != nil {
		img.Logger.Errorf("unable to find a file in request with key %s %s", img.Key, err)
		render.Render(
			// Create is a http.Handler method for handling POST requests
			w, r,
			apherror.ErrServer(
				errors.Wrapf(err, "unable to find a file in request with key %s", img.Key),
			),
		)
		return
	}
	defer file.Close()
	// Open a local file for saving the image content
	name := filepath.Join(img.LocalPath, "uploaded-"+header.Filename)
	writer, err := os.Create(name)
	if err != nil {
		img.Logger.Errorf("unable to open file %s for writing %s", name, err)
		render.Render(
			w, r,
			apherror.ErrServer(
				errors.Wrapf(err, "unable to open file %s for writing", name),
			),
		)
		return
	}
	defer writer.Close()
	_, err = io.Copy(writer, file)
	if err != nil {
		img.Logger.Errorf("unable to write file %s %s", name, err)
		render.Render(
			w, r,
			apherror.ErrServer(
				errors.Wrapf(err, "unable to write to file %s", name),
			),
		)
		return
	}
	img.Logger.Infof("wrote image data to local file %s", name)
	render.NoContent(w, r)
}
