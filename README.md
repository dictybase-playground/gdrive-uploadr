# gdrive-uploadr
A server to upload images to google drive folder.   
There is also a `local` version, it is primarily for testing purpose.

# Usage
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
     help, h               Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --gdrive-secret value, --gs value  gdrive client secret json file [$GDRIVE_CLIENT_SECRET]
   --cache-file value, --cf value     location of cached gdrive token file, defaults to ~/.credentials/gdrive.json [$CACHE_TOKEN_FILE]
   --help, -h                         show help
   --version, -v                      print the version
```

## Subcommands

```
NAME:
   shared-gdrive-folder - create read only public gdrive folder

USAGE:
   gdrive-uploadr shared-gdrive-folder [command options] [arguments...]

OPTIONS:
   --folder value, -f value  Name of the folder[required] [$FOLDER]
```

```
NAME:
   run-gdrive - starts the server for uploading image to google drive

USAGE:
   gdrive-uploadr run-gdrive [command options] [arguments...]

OPTIONS:
   --app-log value        Name of the application log file(optional), default goes to stderr [$APP_LOG]
   --app-log-level value  log level of the application log(optional), default is json (default: "error") [$APP_LOG_LEVEL]
   --app-log-fmt value    Format of the application log(optional), default is json (default: "json") [$APP_LOG_FMT]
   --hooks value          hook names for sending log in addition to stderr
   --slack-channel value  Slack channel where the log will be posted [$SLACK_CHANNEL]
   --slack-url value      Slack webhook url[required if slack channel is provided] [$SLACK_URL]
   --web-log value        Name of the web request log file(optional), default goes to stderr [$WEB-LOG]
   --web-log-fmt value    Format of the web log(optional), default is json (default: "json") [$WEB_LOG_FMT]
   --port value           port on which the server listen (default: 9998)
   --image-key value      The name of form key from where the image file will be retrieved from the request body (default: "image")
   --folder-id value      The folder id of gdrive[required] [$FOLDER_ID]
```


```
NAME:
   authorize - authorize gdrive client

USAGE:
   gdrive-uploadr authorize [arguments...]
```

```
NAME:
   gdrive-folder - create new gdrive folder

USAGE:
   gdrive-uploadr gdrive-folder [command options] [arguments...]

OPTIONS:
   --folder value, -f value  Name of the folder[required] [$FOLDER]
```

