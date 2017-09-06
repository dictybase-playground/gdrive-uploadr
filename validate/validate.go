package validate

import (
	"gopkg.in/urfave/cli.v1"
)

// ValidateGdriveOptions validate presence of *gdrive-secret* value
func ValidateGdriveOptions(c *cli.Context) error {
	if !c.GlobalIsSet("gdrive-secret") {
		return cli.NewExitError("missing command line argument gdrive-secret", 2)
	}
	return nil
}

// ValidateGdriveFolderOptions validate presence of folder name
func ValidateGdriveFolderOptions(c *cli.Context) error {
	if err := ValidateGdriveOptions(c); err != nil {
		return err
	}
	if !c.IsSet("folder") {
		return cli.NewExitError("missing folder argument", 2)
	}
	return nil
}

// ValidateGdriveServer validate presence of folder id
func ValidateGdriveServer(c *cli.Context) error {
	if err := ValidateGdriveOptions(c); err != nil {
		return err
	}
	if !c.IsSet("folder-id") {
		return cli.NewExitError("missing folder-id argument", 2)
	}
	return nil
}
