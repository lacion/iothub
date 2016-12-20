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

### Testing

``make test``