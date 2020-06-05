FROM golang:1.14.3 as build

WORKDIR /go/src/gitlab.com/emanueljoivo/workflows

COPY ./ ./

RUN go install -v

FROM debian:stretch-slim as production

COPY --from=build /go/bin/workflows /usr/local/bin/workflows

RUN chmod +x /usr/local/bin/workflows

EXPOSE 5000

CMD ["/usr/local/bin/workflows"]
