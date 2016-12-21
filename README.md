# iothub

IoT Hub for things


this project is a multiroom echo server using websockets, you could have 
any sensor sending data to a room and have any other device join that room to read
the data in realtime

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/iothub
```

Running with -version will get you the current version and git commit hash for the binary

```console
$ ./bin/iothub -version
```

### Configuration

IoT Hub uses environment variables for configuration all of them are prefixed with IOTHUB.

- IOTHUB_JSON_LOGS defaults to false, if true it will output log in json format.
- IOTHUB_LOGLEVEL defaults to debug
- IOTHUB_MODE this sets the gin mode, defaults to debug, other options are release and test
- IOTHUB_LISTEN_ADDRESS defaults to ":5000"
- IOTHUB_SECRET defaults to "887yff9898yfhuiew3489fy3hewfuig239f8ghew32yfh" it is higly recomended to change this

a note on logging, debug is very verbose as it outputs all the messages the server receives

### Testing

``make test``

## Features

- base project created with [cookiecutter-golang](https://github.com/lacion/cookiecutter-golang)
- Uses [gin](https://github.com/gin-gonic/gin) for http/s
- Uses [melody](github.com/olahol/melody) for websockets
- Uses [logrus](https://github.com/Sirupsen/logrus) for logging
- Uses [viper](https://github.com/spf13/viper) for config