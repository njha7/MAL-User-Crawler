FROM golang:1.11
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./... && go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
WORKDIR /root/
COPY --from=0 /go/src/app .
CMD ["./app"]