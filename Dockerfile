FROM golang:1.12 as builder
WORKDIR /azure-resources-exporter
COPY . .
RUN make getpromu build

FROM quay.io/prometheus/busybox:glibc AS app
LABEL maintainer="FXinnovation CloudToolDevelopment <CloudToolDevelopment@fxinnovation.com>"
COPY --from=builder /azure-resources-exporter/azure-resources-exporter /bin/azure-resources-exporter

EXPOSE      9259
WORKDIR /
ENTRYPOINT  [ "/bin/azure-resources-exporter" ]
