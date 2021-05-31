# repico

REST API based Pi controller

## Introduction

This application has been written as a simple REST based GPIO pin controller (on-off switcher) for Raspberry Pi Zero W. As this application uses standard Linux driver for GPIO operations it should works also on other Pi models as well as on other GPIO-enabled SBCs with Linux on board.

## Downloading and building

Download newest sources directly from GitHub repository:

```bash
go get github.com/markamdev/repico
```

Build using standard *go build* command.

```bash
go build ./
```

If you're not building on a target platform use necessary architecture and OS settings. Below example for building for RPi Zero:

```bash
GOARCH=arm GOOS=linux go build ./
```

## Running RePiCo

Repico is a standalone application and does not need any additional service to be running. Application options can be passed using command line flags and system environment variables.

Table below lists all supported application options:

| ------- | ------- | ------- | ------- |
| System variable | Command line option | Default value | Description |
| ------- | --------| ------- | --------|
| REPICO_PORT | --repico-port | 8080 | Application listening port |
| LOG_LEVEL | --log-level | ERROR | Logging level. Allowed leves are ERROR, DEBUG and VERBOSE |

## Usage

As **repico** is a REST based application it can be fully controlled by HTTP request. Use your HTTP client of choice ([Insomnia](https://insomnia.rest/), [Postman](https://www.postman.com/) or even a command line based [cURL](https://curl.se/) ) to send command to application.

### GPIO pin configuration

To **enable GPIO pin** send HTTP POST request to */v2/gpio* endpoint with pin description in JSON content.

*Request example for output pin*:

```bash
curl -X POST -d '{
"pin" : 1,
"direction" : "out"
}' http://locahost:8080/v2/gpio
```

*Request example for input pin*:

```bash
curl -X POST -d '{
"pin" : 1,
"direction" : "in"
}' http://locahost:8080/v2/gpio
```

If successfully processed HTTP OK (code 200) is returned.

To **disable GPIO pin** send HTTP DELETE request to */v2/gpio/{X}* endpoint where {X} is a PIN number.

*Request example*:

```bash
curl -X DELETE http://localhost:8080/v2/gpio/1
```

### Setting and getting GPIO pin's state

To **set GPIO pin value** send HTTP PATCH request to */v2/gpio/{X}* endpoint (where {X} is a PIN number) with proper JSON data in body

*Request example*:

```bash
curl -X PATCH -d '{ "value" : 1 }' http://localhost:8080/v2/gpio/1
```

Please note that in case of pin configured as *input* it is not possible to set value. In such case API will return HTTP error BadRequest (400).

To **get current GPIO pin value** send HTTP GET request to */v2/gpio/{X}* endpoint (where {X} is a PIN number.

*Request example*:

```bash
curl -X GET http://localhost:8080/v2/gpio/1
```

*Response example*:

```json
{
  "pin": 1,
  "value": 1
}
```

### Listing all exported GPIO pins

It is possible to **list all exported GPIO pins** with their current direction using GET request to main endpoint.

*Request example*:

```bash
curl -X GET http://localhost:8080/v2/gpio
```

*Response example*:

```json
[
  {
    "pin": 1,
    "direction": "out"
  },
  {
    "pin": 5,
    "direction": "in"
  }
]
```

## Testing

### Unit tests

Unit tests are prepared as a pure Go test (with *testing* package used). Additionally *testify* package is used.

To execute unit tests just call command shown below in main repository directory:

```bash
go test -v ./...
```

### Functional tests

## Licensing

Code is published under [MIT License](https://opensource.org/licenses/MIT) as it seems to be the most permissive license. If for some reason you need to have this code published with other license (ex. to reuse the code in your project) please contact [author](#author-/-contact) directly.

## Author / Contact

If you need to contact me feel free to write me an email:
[markamdev.84#dontwantSPAM#gmail.com](mailto:)
