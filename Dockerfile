# Builder 
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go ./main.go
COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY db db

# RUN go build -o /out/app ./cmd/app

RUN CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" -o /out/school-museum .


# Runtime
FROM alpine:3.23.3 AS runtime

WORKDIR /src

COPY --from=builder /out/school-museum /src/school-museum

COPY config/ /src/config/

EXPOSE 8080

ENV PATH_CONFIG=/src/config/config.kdl

ENTRYPOINT ["/src/school-museum", "-t", "kdl"]
