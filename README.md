# docker-compose-lb

A simple reverse proxy to distribute requests between multiple container instances,
build to work with docker-compose.

## Usage

```
app:
  image: corp.com/devision/app

proxy:
  image: nicolai86/docker-compose-reverse-proxy
  ports:
    - 80:8080
  volumes:
    - "${DOCKER_CERT_PATH}:${DOCKER_CERT_PATH}"
  environment:
    DOCKER_HOST: "${DOCKER_HOST}"
    DOCKER_TLS_VERIFY: "${DOCKER_TLS_VERIFY}"
    DOCKER_CERT_PATH: "${DOCKER_CERT_PATH}"
    DOCKER_COMPOSE_SERVICE_NAME: app
```

## Limitations

- only supports HTTP (you need to terminate HTTPS elsewhere)
- assumes your load-balanced service only exposes a single TCP port
- distributes traffic randomly
