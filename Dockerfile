FROM alpine:latest
ARG BINARY_NAME
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY ./$BINARY_NAME .