FROM golang:1.23

# Install dependencies
RUN apt-get update && \
    apt-get install -y \
        sqlite3 \
        libsqlite3-dev \
        firejail

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend .

EXPOSE 6969

CMD ["./backend"]
