FROM golang:1.22-alpine as build

WORKDIR /trade-builder

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /trade-builder .


FROM scratch

COPY --from=build /trade-builder /trade-builder

CMD ["/trade-builder"]