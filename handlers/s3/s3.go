package s3

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/dictybase-playground/gdrive-uploadr/apihelpers/apherror"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var imgRgx = regexp.MustCompile(`(jp(e)*g|png)$`)

// ImageHandler is an net/http handler for managing images
type ImageHandler struct {
	Logger *logrus.Logger
	Key    string
	Bucket string
	Client *minio.Client
	Folder string
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
	thumb, err := generateThumb(file)
	if err != nil {
		img.Logger.Error(err)
		render.Render(
			w, r,
			apherror.ErrServer(errors.Wrap(err, "error in generating thumbnail")),
		)
		return
	}
	_, err = img.Client.PutObject(
		img.Bucket,
		fmt.Sprintf("images/%s/%s/%s", chi.URLParam(r, "year"), img.Folder, strings.ToLower(header.Filename)),
		file,
		header.Size,
		minio.PutObjectOptions{ContentType: detectContentType(header.Filename)},
	)
	if err != nil {
		img.Logger.Errorf("unable to upload file %s to s3 object storeage %s", header.Filename, err)
		render.Render(
			// Create is a http.Handler method for handling POST requests
			w, r,
			apherror.ErrServer(
				errors.Wrapf(err, "unable to upload file %s to s3 object storage", header.Filename),
			),
		)
		return
	}
	img.Logger.Infof("uploaded file %s to object storage", header.Filename)
	// Open a local file for saving the image content
	render.NoContent(w, r)
}

func detectContentType(h string) string {
	switch imgRgx.FindString(strings.ToLower(h)) {
	case "jpg":
	case "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	default:
		return "application/octet-stream"
	}
	return ""
}

func generateThumb(r io.Reader) (image.Image, error) {
	var rimg image.Image
	buff := bytes.NewBuffer(make([]byte, 0))
	tr := io.TeeReader(r, buff)
	img, _, err := image.Decode(tr)
	if err != nil {
		return rimg, fmt.Errorf("error in decoding image %s", err)
	}
	nr := ioutil.NopCloser(buff)
	defer nr.Close()
	config, _, err := image.DecodeConfig(nr)
	if err != nil {
		return rimg, fmt.Errorf("error in getting config %s", err)
	}
	width := 200.0
	if config.Width <= int(width) {
		return img, nil
	}
	ar := float64(config.Height) / float64(config.Width)
	height := ar * width
	return transform.Resize(
		img,
		int(width),
		int(height),
		transform.Lanczos,
	), nil
}
