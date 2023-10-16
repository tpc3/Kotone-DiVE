FROM golang:alpine AS build
ADD . /go/src/Kotone-DiVE/
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 0
WORKDIR /go/src/Kotone-DiVE
RUN go build .

FROM alpine:3.17
COPY --from=build /go/src/Kotone-DiVE/Kotone-DiVE /Kotone-DiVE/
COPY --from=build /go/src/Kotone-DiVE/config-template.yaml /Kotone-DiVE/config.yaml
RUN apk add --no-cache ca-certificates ffmpeg
WORKDIR /Kotone-DiVE
ENTRYPOINT [ "/Kotone-DiVE/Kotone-DiVE" ]
