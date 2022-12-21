FROM golang:1.19-alpine3.17 as builder
RUN apk --no-cache upgrade && apk --no-cache add git make g++
WORKDIR /project/
COPY ./ /project/
RUN make build

FROM alpine
RUN apk --no-cache upgrade && apk --no-cache add ca-certificates
COPY --from=builder /project/bin/ /usr/local/bin/
WORKDIR /home/
CMD ["service"]
