FROM ubuntu

RUN apt update -y \
    && apt upgrade -y

COPY ./experiment-app /experiment/

WORKDIR /experiment