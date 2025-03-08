FROM golang:1.23-alpine AS build
WORKDIR /go/src/jobgolang
COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/jobgolang

FROM alpine
COPY --from=build /go/bin/jobgolang /bin/jobgolang
COPY --from=build /go/src/jobgolang/config /config
ENV APP_ENV=aws_lambda
ENTRYPOINT ["/bin/jobgolang"]