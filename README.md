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

## Usage

As **repico** is a REST based application it can be fully controlled by HTTP request. Use your HTTP client of choice ([Insomnia](https://insomnia.rest/), [Postman](https://www.postman.com/) or even a command line based [cURL](https://curl.se/) ) to send command to application.

### GPIO pin configuration

### Setting and getting GPIO pin's state

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
