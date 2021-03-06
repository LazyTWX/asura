FROM golang:alpine AS builder
ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine

WORKDIR /dist
ARG TOKEN
ARG FIREBASE_CONFIG
ARG FIREBASE_PROJECT_ID
ARG DATADOG_API_KEY
ENV PRODUCTION=TRUE
ENV TOKEN=$TOKEN
ENV FIREBASE_CONFIG=$FIREBASE_CONFIG
ENV FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID
ENV DATADOG_API_KEY=$DATADOG_API_KEY

COPY --from=builder /build/main /dist
COPY --from=builder /build/resources /dist/resources

ENTRYPOINT ./main 