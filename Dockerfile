FROM golang:1.16.2 as build

ENV BIN_FILE /opt/service/service-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/service/*

# тонкий образ
FROM alpine:3.9

LABEL SERVICE="service"

ENV BIN_FILE "/opt/service/service-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

EXPOSE 8080

ENTRYPOINT ${BIN_FILE} -port 8080