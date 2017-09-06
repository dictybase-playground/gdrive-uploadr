package auth

import (
	"context"
	"fmt"
	"io/ioutil"

	drive "google.golang.org/api/drive/v3"

	"golang.org/x/oauth2/google"
	"gopkg.in/urfave/cli.v1"
)

func GetGDriveClient(c *cli.Context) (*drive.Service, error) {
	var srv *drive.Service
	cacheFile, err := TokenCacheFile(c)
	if err != nil {
		return srv, fmt.Errorf("error unable to set the token file path %s\n", err)
	}
	tok, err := TokenFromFile(cacheFile)
	if err != nil {
		return srv, fmt.Errorf("error unable to get token from cache file: possibly run the authorize-drive command %s", err)
	}
	cont, err := ioutil.ReadFile(c.GlobalString("gdrive-secret"))
	if err != nil {
		return srv, fmt.Errorf("error unable to read the secret json file %s\n", err)
	}
	config, err := google.ConfigFromJSON(
		cont,
		drive.DriveScope,
		drive.DriveMetadataScope,
	)
	if err != nil {
		return srv, fmt.Errorf("error unable to create oauth config from secret file %s\n", err)
	}
	client := config.Client(context.Background(), tok)
	srv, err = drive.New(client)
	if err != nil {
		return srv, fmt.Errorf("error unable to set gdrive client %s\n", err)
	}
	return srv, nil
}
