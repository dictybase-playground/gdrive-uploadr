package s3

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"mime/multipart"
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
	Logger      *logrus.Logger
	Key         string
	Bucket      string
	Client      *minio.Client
	Folder      string
	ThumbFolder string
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
	treader, err := img.generateThumb(file)
	if err != nil {
		img.Logger.Error(err)
		render.Render(
			w, r,
			apherror.ErrServer(errors.Wrap(err, "")),
		)
		return
	}
	// rewind the image file for re-reading after thumbnail creation
	_, err = file.Seek(0, 0)
	if err != nil {
		img.Logger.Errorf("error in seeking original image file %s %s", header.Filename, err)
		render.Render(
			w, r,
			apherror.ErrServer(
				errors.Wrapf(err, "error in seeking original image file %s %s", header.Filename, err),
			),
		)
		return
	}
	year := chi.URLParam(r, "year")
	origPath, err := img.storeOriginal(file, header, year)
	if err != nil {
		render.Render(
			w, r,
			apherror.ErrServer(errors.Wrap(err, "")),
		)
		return
	}
	thumbPath, err := img.storeThumbnail(treader, header, year)
	if err != nil {
		render.Render(
			w, r,
			apherror.ErrServer(errors.Wrap(err, "")),
		)
		return
	}
	img.Logger.Infof("uploaded %s file to %s and %s in object storage", header.Filename, origPath, thumbPath)
	// Open a local file for saving the image content
	render.NoContent(w, r)
}

func (img *ImageHandler) storeOriginal(r io.Reader, header *multipart.FileHeader, year string) (string, error) {
	origPath := fmt.Sprintf(
		"images/%s/%s/%s",
		year,
		img.Folder,
		strings.ToLower(header.Filename),
	)
	_, err := img.Client.PutObject(
		img.Bucket,
		origPath,
		r,
		header.Size,
		minio.PutObjectOptions{ContentType: detectContentType(header.Filename)},
	)
	if err != nil {
		img.Logger.Errorf("unable to upload file %s to s3 object storage %s", header.Filename, err)
		return origPath,
			fmt.Errorf("unable to upload file %s to s3 object storage %s", header.Filename, err)
	}
	return origPath, nil
}

func (img *ImageHandler) storeThumbnail(r io.Reader, header *multipart.FileHeader, year string) (string, error) {
	thumbPath := fmt.Sprintf(
		"images/%s/%s/thumb_%s",
		year,
		img.ThumbFolder,
		strings.ToLower(header.Filename),
	)
	_, err := img.Client.PutObject(
		img.Bucket,
		thumbPath,
		r,
		-1,
		minio.PutObjectOptions{ContentType: detectContentType(header.Filename)},
	)
	if err != nil {
		img.Logger.Errorf("unable to upload thumbnail %s to s3 object storage %s", header.Filename, err)
		return thumbPath,
			errors.Errorf("unable to upload thumbnail file %s to s3 object storage %s", header.Filename, err)
	}
	return thumbPath, nil
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

func (img *ImageHandler) generateThumb(r io.Reader) (io.Reader, error) {
	buff := bytes.NewBuffer(make([]byte, 0))
	tr := io.TeeReader(r, buff)
	dimg, _, err := image.Decode(tr)
	if err != nil {
		return buff, errors.Errorf("error in decoding image %s", err)
	}
	nr := ioutil.NopCloser(buff)
	defer nr.Close()
	config, _, err := image.DecodeConfig(nr)
	if err != nil {
		return buff, errors.Errorf("error in getting config %s", err)
	}
	var thumbImg image.Image
	width := 200.0
	if config.Width <= int(width) {
		thumbImg = dimg
	} else {
		ar := float64(config.Height) / float64(config.Width)
		height := ar * width
		thumbImg = transform.Resize(
			dimg,
			int(width),
			int(height),
			transform.Lanczos,
		)
	}
	tbuff := bytes.NewBuffer(make([]byte, 0))
	if err := jpeg.Encode(tbuff, thumbImg, &jpeg.Options{Quality: 75}); err != nil {
		img.Logger.Errorf("error in copying thumbnail to buffer %s", err)
		return tbuff, errors.Errorf("error in copying thumbnail to buffer %s", err)
	}
	return tbuff, nil
}
