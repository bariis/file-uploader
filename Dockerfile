FROM golang:1.16-alpine
WORKDIR /src
COPY go.sum go.mod /
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/file-uploader .

FROM alpine:3.15.0
COPY --from=0 /bin/file-uploader /bin/file-uploader
COPY ./images /bin/images
ENTRYPOINT ["/bin/file-uploader"]
