# syntax=docker/dockerfile:1

##
## Build
##
# FROM golang:1.16-buster AS build
FROM registry.duvall.org.uk/docker1/containers/golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /exam-app-backend

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /exam-app-backend /exam-app-backend

EXPOSE 8081

USER nonroot:nonroot

ENTRYPOINT ["/exam-app-backend"]
