#build stage
FROM golang:alpine AS builder

WORKDIR /app
# COPY go.mod .
# COPY go.sum .
COPY . .
RUN go mod download

RUN go build -o demyst main.go

#run stage
FROM alpine 
WORKDIR /app
COPY --from=builder /app/demyst .
COPY .env .

EXPOSE 5000
CMD ["/app/demyst"]

