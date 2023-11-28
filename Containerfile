FROM docker.io/golang:alpine as build

RUN mkdir -p /go/pcsmpages
COPY . /go/pcsmpages
RUN cd /go/pcsmpages && \
        go mod vendor && \
        CGO_ENABLED=0 go build -o /pcsmpages main.go

FROM scratch
COPY --from=build /pcsmpages /
COPY theme /theme
ENTRYPOINT ["/pcsmpages"]
EXPOSE 1323/tcp
