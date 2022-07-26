#Build binary in alpine container
FROM golang:alpine AS alpine_builder

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/main main.go 

#Copy and run binary
FROM scratch

COPY --from=alpine_builder /go/bin/main /go/bin/main

EXPOSE 5555

ENTRYPOINT ["/go/bin/main"]