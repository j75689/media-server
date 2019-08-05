FROM golang as build-env
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o media-server

FROM alpine
COPY --from=build-env /app/media-server /media-server
WORKDIR /
CMD ["./media-server"]