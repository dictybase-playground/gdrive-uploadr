package commands

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dictybase-playground/gdrive-uploadr/auth"
	"github.com/dictybase-playground/gdrive-uploadr/handlers/gdrive"
	"github.com/dictybase-playground/gdrive-uploadr/handlers/local"
	"github.com/dictybase-playground/gdrive-uploadr/handlers/s3"
	"github.com/dictybase-playground/gdrive-uploadr/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/minio/minio-go"
	"gopkg.in/urfave/cli.v1"
)

// RunGdriveServer starts a http server for managing images with gdrive storage
func RunGdriveServer(c *cli.Context) error {
	// logging middleware
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	appLogger, err := logger.GetAppLogger(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	client, err := auth.GetGDriveClient(c)
	if err != nil {
		return err
	}
	// handler
	imgHandler := &gdrive.ImageHandler{
		Key:      c.String("image-key"),
		FolderId: c.String("folder-id"),
		Logger:   appLogger,
		Client:   client,
	}
	// cors middleware
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	})
	// router
	r := chi.NewRouter()
	r.Use(lmw.Middleware)
	r.Use(crs.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Route("/images", func(r chi.Router) {
		r.Post("/", imgHandler.Create)
	})
	fmt.Printf("starting gdrive server at port %d\n", c.Int("port"))
	http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), r)
	return nil
}

// RunLocalServer starts a http server for local management for images
func RunLocalServer(c *cli.Context) error {
	var savePath string
	if c.IsSet("save-path") {
		savePath = c.String("save-path")
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("unable to get current working directory %s", err),
				2,
			)
		}
		savePath = dir
	}
	// logging middleware
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	appLogger, err := logger.GetAppLogger(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	// handler
	imgHandler := &local.ImageHandler{
		Logger:    appLogger,
		Key:       c.String("image-key"),
		LocalPath: savePath,
	}

	// cors middleware
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	})
	// router
	r := chi.NewRouter()
	r.Use(lmw.Middleware)
	r.Use(crs.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Route("/images", func(r chi.Router) {
		r.Post("/", imgHandler.Create)
	})
	fmt.Printf("starting local server at port %d\n", c.Int("port"))
	http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), r)
	return nil
}

// RunS3Server starts a http server for s3 backend
func RunS3Server(c *cli.Context) error {
	// logging middleware
	lmw, err := logger.GetLoggerMiddleware(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	appLogger, err := logger.GetAppLogger(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	s3Client, err := getS3Client(c)
	if err != nil {
		return cli.NewExitError(err.Error(), 2)
	}
	// handler
	imgHandler := &s3.ImageHandler{
		Key:    c.String("image-key"),
		Bucket: c.String("bucket"),
		Folder: c.String("subfolder"),
		Client: s3Client,
		Logger: appLogger,
	}

	// cors middleware
	crs := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	})
	// router
	r := chi.NewRouter()
	r.Use(lmw.Middleware)
	r.Use(crs.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Route(`/images/{year:[0-9]{4}}`, func(r chi.Router) {
		r.Post("/", imgHandler.Create)
	})
	fmt.Printf("starting s3 backend server at port %d\n", c.Int("port"))
	http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), r)
	return nil
}

func getS3Client(c *cli.Context) (*minio.Client, error) {
	s3Client, err := minio.New(
		fmt.Sprintf("%s:%s", c.String("s3-host"), c.String("s3-port")),
		c.String("access-key"),
		c.String("secret-key"),
		false,
	)
	if err != nil {
		return s3Client, fmt.Errorf("unable create the client %s", err.Error())
	}
	return s3Client, nil
}
