ARG GO_VERSION=PLEASE_PROVIDE_GO_VERSION

FROM golang:${GO_VERSION}

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /mononoke-go

ARG VCS_REF="latest"
LABEL org.opencontainers.image.source="https://github.com/mexfoo/mononoke-go" \
  org.opencontainers.image.revision=${VCS_REF}

CMD ["/mononoke-go"]