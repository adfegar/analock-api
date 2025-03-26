FROM golang:alpine

ARG DB_URL
ARG DB_TOKEN

LABEL maintainer="adferdev"

RUN apk update && apk add --no-cache \
    git \
    bash \
    build-base \
    libc6-compat \
    pkgconfig

RUN mkdir /app

WORKDIR /app

COPY . .

# Create .env file
RUN echo "TURSO_DB_URL=$DB_URL" > .env && \
    echo "TURSO_DB_TOKEN=$DB_TOKEN" >> .env

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]
