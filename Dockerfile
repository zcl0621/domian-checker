FROM golang:1.19 AS build-dist

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags netcgo -installsuffix cgo -o /bin/app main.go

FROM redhat/ubi8-minimal:8.8 as prod

RUN export TZ='Asia/Shanghai' && \
    microdnf reinstall tzdata -y

COPY --from=build-dist /bin/app /bin/app

RUN chmod +x /bin/app

ENV RUN_MODE='release'

CMD ["/bin/app"]
EXPOSE 8080