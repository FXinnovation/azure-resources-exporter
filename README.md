# azure-resources-exporter
[![Build Status](https://travis-ci.org/FXinnovation/azure-resources-exporter.svg?branch=master)](https://travis-ci.org/FXinnovation/azure-resources-exporter)

**Warning, this exporter is in early stages, and is not ready to be used**

Prometheus exporter exposing Azure resources API results as metrics

## Getting Started

### Prerequisites

To run this project, you will need a [working Go environment](https://golang.org/doc/install).

### Installing

```bash
go get -u github.com/FXinnovation/azure-resources-exporter
```

## Building

Build the sources with

```bash
make build
```

## Run the binary

```bash
./azure-resources-exporter
```

By default, the exporter config is expected in config.yml.

Use -h flag to list available options.

## Testing

### Running unit tests

```bash
make test
```

## Docker image

You can build a docker image using:
```bash
make docker
```
The resulting image is named `fxinnovation/azure-resources-exporter:<git-branch>`.
It exposes port 9259 and expects the config in /config.yml. To configure it, you can bind-mount a config from your host: 
```
$ docker run -p 9259:9259 -v /path/on/host/config.yml:/opt/azure-resources-exporter/config.yml fxinnovation/azure-resources-exporter:<git-branch>
```

## Contributing

Refer to [CONTRIBUTING.md](https://github.com/FXinnovation/azure-resources-exporter/blob/master/CONTRIBUTING.md).

## License

Apache License 2.0, see [LICENSE](https://github.com/FXinnovation/azure-resources-exporter/blob/master/LICENSE).
