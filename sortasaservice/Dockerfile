FROM golang:1.24-alpine AS build

WORKDIR /build
COPY * .
RUN go build .

FROM alpine:3

COPY --from=build /build/sortasaservice /sortasaservice
CMD /sortasaservice 