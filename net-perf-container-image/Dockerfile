FROM ubuntu

RUN apt update -y \
    && apt upgrade -y \
    && apt install inetutils-ping curl iproute2 mysql-client -y

COPY ./client.sh /net-perf-container-image/
COPY ./server.sh /net-perf-container-image/

WORKDIR /net-perf-container-image

ENTRYPOINT ["bash", "server.sh"]