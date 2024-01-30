# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY . .

ARG COMPONENT=dataprovider
ENV COMPONENT=${COMPONENT}
ENV CGO_ENABLED=1
RUN go build -o component ./$COMPONENT


CMD ["sh", "-c", "/app/component"]
