FROM golang:1.26.3

WORKDIR /app

ENV GOTOOLCHAIN=local GOPROXY=https://goproxy.cn,direct GOSUMDB=sum.golang.google.cn

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build ./...

CMD ["bash"]
