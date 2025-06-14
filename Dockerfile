FROM golang:bookworm

ARG DB_URL
ARG DB_TOKEN

LABEL maintainer="adferdev"

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    pkg-config

RUN mkdir /app

WORKDIR /app

COPY . .

# Create .env file
RUN echo "TURSO_DB_URL=$DB_URL" > .env && \
    echo "TURSO_DB_TOKEN=$DB_TOKEN" >> .env \
    echo "API_ENVIRONMENT=production" >> .env \
    echo "API_URL_HOST=analock" >> .env

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build

EXPOSE 8080

CMD [ "/build" ]
