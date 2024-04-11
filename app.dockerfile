FROM golang:1.21
#WORKDIR /app
#COPY go.mod go.sum ./
#RUN go mod download
#COPY ./ ./
#RUN go build -o /main main.go
#CMD ["./main"]
LABEL autor=ias
LABEL task=6
WORKDIR /app
COPY . .
RUN go build -o ./cmd/server/xx ./cmd/server 
EXPOSE 8080
ENTRYPOINT ["/app/cmd/server/xx"]