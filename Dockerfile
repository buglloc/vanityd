FROM golang:1.23.3 as build

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /vanityd

FROM debian:bookworm-slim

COPY --from=build /vanityd /usr/sbin/vanityd

CMD ["/usr/sbin/vanityd", "--addr=:80"]
