FROM alpine:latest

RUN apk add --update ca-certificates

COPY dist/convoy /

ENTRYPOINT ["/convoy"]
