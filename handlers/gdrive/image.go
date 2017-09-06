package gdrive

import (
	"net/http"

	"github.com/dictybase-playground/gdrive-uploadr/apihelpers/apherror"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	drive "google.golang.org/api/drive/v3"
)

type ImageHandler struct {
	Logger   *logrus.Logger
	Client   *drive.Service
	Key      string
	FolderId string
}

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
	gfile, err := img.Client.Files.Create(
		&drive.File{
			Parents: []string{img.FolderId},
			Name:    header.Filename,
		},
	).Fields("id", "webContentLink", "webViewLink").Media(file).Do()
	if err != nil {
		img.Logger.Errorf("unable to upload file %s %s", header.Filename, err)
		render.Render(w, r, apherror.ErrGdrive(err))
		return
	}
	img.Logger.Infof(
		"uploaded file %s to gdrive with id %s and location %s",
		header.Filename,
		gfile.Id,
		gfile.WebContentLink,
	)
	w.Header().Set("Location", gfile.WebContentLink)
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, struct {
		Id  string
		Url string
	}{
		gfile.Id,
		gfile.WebViewLink,
	})
}
