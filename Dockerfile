FROM golang:1.19.0-alpine3.16 as build

COPY . /irsa-emu/

WORKDIR /irsa-emu
RUN apk update \
  && \
  apk add make=4.3-r0 \
  && \
  make build

FROM scratch

COPY --from=build /irsa-emu/bin/irsa-emu-webhook /irsa-emu-webhook
ENTRYPOINT ["/irsa-emu-webhook"]
