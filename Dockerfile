FROM golang:alpine AS build
ADD . /go/src/Kotone-DiVE/
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 0
WORKDIR /go/src/Kotone-DiVE
RUN go build .

FROM alpine
COPY --from=build /go/src/Kotone-DiVE/Kotone-DiVE /go/src/Kotone-DiVE/config.yaml /Kotone-DiVE/
RUN apk add --no-cache ca-certificates ffmpeg
WORKDIR /Kotone-DiVE
ENTRYPOINT [ "/Kotone-DiVE/Kotone-DiVE" ]