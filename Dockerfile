FROM golang:alpine

RUN apk add --no-cache g++ git

RUN mkdir /app
WORKDIR /app
COPY . ./

RUN go mod download
RUN go build -o /weather-web-service-go

EXPOSE 8080

CMD [ "/weather-web-service-go" ]
