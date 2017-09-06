package commands

import (
	"fmt"
	"io/ioutil"

	drive "google.golang.org/api/drive/v3"

	"github.com/dictybase-playground/gdrive-uploadr/auth"
	"golang.org/x/oauth2/google"
	"gopkg.in/urfave/cli.v1"
)

// CreateSharedGdriveFolderAction creates a new public readonly gdrive folder
func CreateSharedGdriveFolderAction(c *cli.Context) error {
	client, err := auth.GetGDriveClient(c)
	if err != nil {
		return err
	}
	file, err := client.Files.Create(
		&drive.File{
			Name:     c.String("folder"),
			MimeType: "application/vnd.google-apps.folder",
		},
	).Fields("id").Do()
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to create folder %s %s", c.String("folder"), err.Error()),
			2,
		)
	}
	fmt.Printf("created folder %s with id %s\n", c.String("folder"), file.Id)
	_, err = client.Permissions.Create(
		file.Id,
		&drive.Permission{
			Role: "reader",
			Type: "anyone",
		},
	).Do()
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unbale to create readonly public permission for folder-id %s %s", file.Id, err),
			2,
		)
	}
	fmt.Printf("created readonly public permission for folder-id %s\n", file.Id)
	return nil
}

// CreateGdriveFolderAction creates a new gdrive folder
func CreateGdriveFolderAction(c *cli.Context) error {
	client, err := auth.GetGDriveClient(c)
	if err != nil {
		return err
	}
	file, err := client.Files.Create(
		&drive.File{
			Name:     c.String("folder"),
			MimeType: "application/vnd.google-apps.folder",
		},
	).Fields("id").Do()
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("unable to create folder %s %s", c.String("folder"), err.Error()),
			2,
		)
	}
	fmt.Printf("created folder %s with id %s", c.String("folder"), file.Id)
	return nil
}

// AuthGDriveAction retrieves and saves a temporary token for further use
func AuthGDriveAction(c *cli.Context) error {
	tokenFile, err := auth.TokenCacheFile(c)
	if err != nil {
		cli.NewExitError(fmt.Sprintf("error unable to set the token file path %s\n", err), 2)
	}
	cont, err := ioutil.ReadFile(c.GlobalString("gdrive-secret"))
	if err != nil {
		cli.NewExitError(fmt.Sprintf("error unable to read the secret json file %s\n", err), 2)
	}
	config, err := google.ConfigFromJSON(
		cont,
		drive.DriveScope,
		drive.DriveMetadataScope,
	)
	if err != nil {
		cli.NewExitError(fmt.Sprintf("error unable to create oauth config from secret file %s\n", err), 2)
	}
	tok := auth.GetTokenFromWeb(config)
	auth.SaveToken(tokenFile, tok)
	return nil
}
