# Builder 
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./main.go ./main.go
COPY cmd cmd
COPY pkg pkg
COPY internal internal
# COPY db db
# COPY config config

# RUN go build -o /out/app ./cmd/app

RUN CGO_ENABLED=0 \
    go build -trimpath -ldflags "-s -w" -o /out/{{.APP_NAME}} .

# ENV PATH_CONFIG=/src/config/config.yaml

# ENTRYPOINT ["/out/app"]


# Runtime
FROM alpine:3.23.3 AS runtime

WORKDIR /src

COPY --from=builder /out/{{.APP_NAME}} /src/{{.APP_NAME}}

COPY config/ /src/config/

EXPOSE 8080

ENV PATH_CONFIG=/src/config/config.kdl

ENTRYPOINT ["/src/{{.APP_NAME}}", "-t", "kdl"]
