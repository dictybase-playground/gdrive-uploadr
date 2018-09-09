package validate

import (
	"fmt"

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

// ValidateS3Server validate cmdline arguments of s3-server
func ValidateS3Server(c *cli.Context) error {
	for _, p := range []string{"access-key", "secret-key", "bucket", "redis-master", "redis-slave"} {
		if !c.IsSet(p) {
			return cli.NewExitError(
				fmt.Sprintf("missing argument %s", p),
				2,
			)
		}
	}
	return nil
}
