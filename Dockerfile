FROM scratch
MAINTAINER Raphael Randschau<nicolai86@me.com>

EXPOSE 8080

WORKDIR /src
ADD reverse-proxy /src

CMD ["./reverse-proxy"]
