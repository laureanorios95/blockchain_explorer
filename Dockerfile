FROM golang:1.19

WORKDIR /usr/src/blockchain-explorerCRD

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
ENV GOOGLE_APPLICATION_CREDENTIALS=./dogwood-theorem-312714-4a9e90bc3a8f.json
RUN go build -v -o /usr/local/bin/blockchain-explorerCRD ./...

EXPOSE 8080

CMD ["blockchain-explorerCRD"]