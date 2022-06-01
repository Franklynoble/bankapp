
#Build stage
#also in the build  stage Download and  extract the golang Migrate binaries to the docker image,
#also install the curl command into the image using RUN apk add curl 
FROM golang:1.16-alpine3.13 AS builder 
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl  
RUN  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz     
 
 #in the Run Stage,  copy the golang migrate to  file /app/migrate.linux-amd64 the path to the final image , then put it in the final run stage image ./migrate 
#Run stage
FROM alpine:3.13
WORKDIR /app

COPY --from=builder /app/main . 
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY start.sh . 
COPY wait-for.sh . 
COPY db/migration ./migration

EXPOSE 8080

CMD ["/app/main"]

ENTRYPOINT [ "/app/start.sh" ]
