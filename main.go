package main

import (
	"log"
	"os"

	"gopkg.in/urfave/cli.v1"

	"github.com/dictybase-playground/gdrive-uploadr/commands"
	"github.com/dictybase-playground/gdrive-uploadr/validate"
)

func main() {
	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Name = "gdrive-uploadr"
	app.Usage = "Manage google drive upload"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "gdrive-secret, gs",
			Usage:  "gdrive client secret json file",
			EnvVar: "GDRIVE_CLIENT_SECRET",
		},
		cli.StringFlag{
			Name:   "cache-file, cf",
			Usage:  "location of cached gdrive token file, defaults to ~/.credentials/gdrive.json",
			EnvVar: "CACHE_TOKEN_FILE",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "authorize",
			Usage:  "authorize gdrive client",
			Action: commands.AuthGDriveAction,
			Before: validate.ValidateGdriveOptions,
		},
		{
			Name:   "gdrive-folder",
			Usage:  "create new gdrive folder",
			Action: commands.CreateGdriveFolderAction,
			Before: validate.ValidateGdriveFolderOptions,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "folder,f",
					Usage:  "Name of the folder[required]",
					EnvVar: "FOLDER",
				},
			},
		},
		{
			Name:   "shared-gdrive-folder",
			Usage:  "create read only public gdrive folder",
			Action: commands.CreateSharedGdriveFolderAction,
			Before: validate.ValidateGdriveFolderOptions,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "folder,f",
					Usage:  "Name of the folder[required]",
					EnvVar: "FOLDER",
				},
			},
		},
		{
			Name:   "run-gdrive",
			Usage:  "starts the server for uploading image to google drive",
			Action: commands.RunGdriveServer,
			Before: validate.ValidateGdriveServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "app-log",
					Usage:  "Name of the application log file(optional), default goes to stderr",
					EnvVar: "APP_LOG",
				},
				cli.StringFlag{
					Name:   "app-log-level",
					Usage:  "log level of the application log(optional), default is json",
					Value:  "error",
					EnvVar: "APP_LOG_LEVEL",
				},
				cli.StringFlag{
					Name:   "app-log-fmt",
					Usage:  "Format of the application log(optional), default is json",
					Value:  "json",
					EnvVar: "APP_LOG_FMT",
				},
				cli.StringSliceFlag{
					Name:  "hooks",
					Usage: "hook names for sending log in addition to stderr",
					Value: &cli.StringSlice{},
				},
				cli.StringFlag{
					Name:   "slack-channel",
					EnvVar: "SLACK_CHANNEL",
					Usage:  "Slack channel where the log will be posted",
				},
				cli.StringFlag{
					Name:   "slack-url",
					EnvVar: "SLACK_URL",
					Usage:  "Slack webhook url[required if slack channel is provided]",
				},
				cli.StringFlag{
					Name:   "web-log",
					Usage:  "Name of the web request log file(optional), default goes to stderr",
					EnvVar: "WEB-LOG",
				},
				cli.StringFlag{
					Name:   "web-log-fmt",
					Usage:  "Format of the web log(optional), default is json",
					Value:  "json",
					EnvVar: "WEB_LOG_FMT",
				},
				cli.IntFlag{
					Name:  "port",
					Usage: "port on which the server listen",
					Value: 9998,
				},
				cli.StringFlag{
					Name:  "image-key",
					Usage: "The name of form key from where the image file will be retrieved from the request body",
					Value: "image",
				},
				cli.StringFlag{
					Name:   "folder-id",
					Usage:  "The folder id of gdrive[required]",
					EnvVar: "FOLDER_ID",
				},
			},
		},
		{
			Name:   "run-local",
			Usage:  "starts the server for uploading image to local storage",
			Action: commands.RunLocalServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "app-log",
					Usage:  "Name of the application log file(optional), default goes to stderr",
					EnvVar: "APP_LOG",
				},
				cli.StringFlag{
					Name:   "app-log-level",
					Usage:  "log level of the application log(optional), default is json",
					Value:  "error",
					EnvVar: "APP_LOG_LEVEL",
				},
				cli.StringFlag{
					Name:   "app-log-fmt",
					Usage:  "Format of the application log(optional), default is json",
					Value:  "json",
					EnvVar: "APP_LOG_FMT",
				},
				cli.StringSliceFlag{
					Name:  "hooks",
					Usage: "hook names for sending log in addition to stderr",
					Value: &cli.StringSlice{},
				},
				cli.StringFlag{
					Name:   "slack-channel",
					EnvVar: "SLACK_CHANNEL",
					Usage:  "Slack channel where the log will be posted",
				},
				cli.StringFlag{
					Name:   "slack-url",
					EnvVar: "SLACK_URL",
					Usage:  "Slack webhook url[required if slack channel is provided]",
				},
				cli.StringFlag{
					Name:   "web-log",
					Usage:  "Name of the web request log file(optional), default goes to stderr",
					EnvVar: "WEB-LOG",
				},
				cli.StringFlag{
					Name:   "web-log-fmt",
					Usage:  "Format of the web log(optional), default is json",
					Value:  "json",
					EnvVar: "WEB_LOG_FMT",
				},
				cli.IntFlag{
					Name:  "port",
					Usage: "port on which the server listen",
					Value: 9998,
				},
				cli.StringFlag{
					Name:  "image-key",
					Usage: "The name of form key from where the image file will be retrieved from the request body",
					Value: "image",
				},
				cli.StringFlag{
					Name:  "save-path",
					Usage: "The local path where the image file will be saved, default is the current folder",
				},
			},
		},
		{
			Name:   "run-s3",
			Usage:  "starts the server for uploading image to s3 storage",
			Action: commands.RunS3Server,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "app-log",
					Usage:  "Name of the application log file(optional), default goes to stderr",
					EnvVar: "APP_LOG",
				},
				cli.StringFlag{
					Name:   "app-log-level",
					Usage:  "log level of the application log(optional), default is json",
					Value:  "error",
					EnvVar: "APP_LOG_LEVEL",
				},
				cli.StringFlag{
					Name:   "app-log-fmt",
					Usage:  "Format of the application log(optional), default is json",
					Value:  "json",
					EnvVar: "APP_LOG_FMT",
				},
				cli.StringSliceFlag{
					Name:  "hooks",
					Usage: "hook names for sending log in addition to stderr",
					Value: &cli.StringSlice{},
				},
				cli.StringFlag{
					Name:   "slack-channel",
					EnvVar: "SLACK_CHANNEL",
					Usage:  "Slack channel where the log will be posted",
				},
				cli.StringFlag{
					Name:   "slack-url",
					EnvVar: "SLACK_URL",
					Usage:  "Slack webhook url[required if slack channel is provided]",
				},
				cli.IntFlag{
					Name:  "port",
					Usage: "port on which the server listen",
					Value: 9998,
				},
				cli.StringFlag{
					Name:  "image-key",
					Usage: "The name of form key from where the image file will be retrieved from the request body",
					Value: "image",
				},
				cli.StringFlag{
					Name:   "s3-host",
					Usage:  "S3 server host",
					EnvVar: "MINIO_SERVICE_HOST",
					Value:  "minio",
				},
				cli.StringFlag{
					Name:   "s3-port",
					Usage:  "S3 server port",
					EnvVar: "MINIO_SERVICE_PORT",
				},
				cli.StringFlag{
					Name:  "s3-bucket",
					Usage: "S3 bucket where the image will be saved",
					Value: "content",
				},
				cli.StringFlag{
					Name:   "access-key, akey",
					EnvVar: "S3_ACCESS_KEY",
					Usage:  "access key for S3 server, required based on command run",
				},
				cli.StringFlag{
					Name:   "secret-key, skey",
					EnvVar: "S3_SECRET_KEY",
					Usage:  "secret key for S3 server, required based on command run",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
