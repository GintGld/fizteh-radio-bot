FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

RUN --mount=type=cache,target=/go/pkg/mod/ \
    go build -o bot ./cmd/bot

FROM alpine as final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

WORKDIR /bot

# copy executables
COPY --from=builder /build/bot /bot/bot

# TODO: enable for webhook
# EXPOSE 8082

ENTRYPOINT [ "/bot/bot" ]