FROM --platform=linux/amd64 golang
WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /app

ENTRYPOINT ["/app"]
