# gdrive-uploadr
A server to upload images to 
+ google drive folder.   
+ s3 storage
+ local filesystem

# Usage
## google drive only
To use the client, generate the `client secret` as described
[here](https://developers.google.com/drive/v3/web/quickstart/go). Then run the
`authorize` subcommand to create a token file for authorizing the client. Then
use both `client secret` and `token` files to run the server and/or the other
subcommands.

# Available commands
```
NAME:
   gdrive-uploadr - Manage google drive upload

USAGE:
   gdrive-uploadr [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
     authorize             authorize gdrive client
     gdrive-folder         create new gdrive folder
     shared-gdrive-folder  create read only public gdrive folder
     run-gdrive            starts the server for uploading image to google drive
     run-local             starts the server for uploading image to local storage
     run-s3                starts the server for uploading image to s3 storage
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gdrive-secret value, --gs value  gdrive client secret json file [$GDRIVE_CLIENT_SECRET]
   --cache-file value, --cf value     location of cached gdrive token file, defaults to ~/.credentials/gdrive.json [$CACHE_TOKEN_FILE]
   --help, -h                         show help
   --version, -v                      print the version
```

```
NAME:
   gdrive-uploadr run-s3 - starts the server for uploading image to s3 storage

USAGE:
   gdrive-uploadr run-s3 [command options] [arguments...]

OPTIONS:
   --app-log value                   Name of the application log file(optional), default goes to stderr [$APP_LOG]
   --app-log-level value             log level of the application log(optional), default is json (default: "error") [$APP_LOG_LEVEL]
   --app-log-fmt value               Format of the application log(optional), default is json (default: "json") [$APP_LOG_FMT]
   --hooks value                     hook names for sending log in addition to stderr
   --slack-channel value             Slack channel where the log will be posted [$SLACK_CHANNEL]
   --slack-url value                 Slack webhook url[required if slack channel is provided] [$SLACK_URL]
   --port value                      port on which the server listen (default: 9998)
   --image-key value                 The name of form key from where the image file will be retrieved from the request body (default: "image")
   --s3-host value                   S3 server host (default: "minio") [$MINIO_SERVICE_HOST]
   --s3-port value                   S3 server port [$MINIO_SERVICE_PORT]
   --s3-bucket value                 S3 bucket where the image will be saved (default: "content")
   --access-key value, --akey value  access key for S3 server, required based on command run [$S3_ACCESS_KEY]
   --secret-key value, --skey value  secret key for S3 server, required based on command run [$S3_SECRET_KEY]
   --log-file value                  name of log file
   --log-fmt value                   Format of the web log(optional), default is json (default: "json")
```

```
NAME:
   gdrive-uploadr gdrive-folder - create new gdrive folder

USAGE:
   gdrive-uploadr gdrive-folder [command options] [arguments...]

OPTIONS:
   --folder value, -f value  Name of the folder[required] [$FOLDER]
```
