FROM golang:1.9 as builder

RUN curl https://glide.sh/get | sh

WORKDIR  /go/src/github.com/jnummelin/graceful-stop-test

# Add dependency graph and vendor it in
ADD glide.yaml glide.lock /go/src/github.com/jnummelin/graceful-stop-test/
RUN glide install

# Add source and compile
ADD graceful-stop.go /go/src/github.com/jnummelin/graceful-stop-test/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -o graceful-stop .

FROM scratch

COPY --from=builder /go/src/github.com/jnummelin/graceful-stop-test/graceful-stop .

CMD ["./graceful-stop"]
