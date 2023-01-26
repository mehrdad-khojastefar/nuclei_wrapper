FROM golang:1.19.5-alpine3.17

WORKDIR /app 
EXPOSE 8000
COPY . /app

RUN apk update && \
    apk add build-base

RUN go mod tidy && \
    sed -i 's/int64(rateLimit)/uint(rateLimit)/g' /go/pkg/mod/github.com/projectdiscovery/subfinder/v2@v2.5.5/pkg/subscraping/agent.go && \
    go build -o /nuclei-wrapper


RUN chmod +x /nuclei-wrapper

CMD [ "/nuclei-wrapper" ]

