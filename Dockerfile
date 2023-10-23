FROM docker:dind

ARG NETWORK=prod
ARG BRANCH=main

RUN apk add bash git

VOLUME /var/k8s/creator-node-db
VOLUME /var/k8s/mediorum
VOLUME /var/k8s/discovery-provider-db
VOLUME /var/k8s/discovery-provider-chain

WORKDIR /root
RUN git clone --single-branch --branch "$BRANCH" https://github.com/AudiusProject/audius-docker-compose.git ./audius-docker-compose

WORKDIR /root/audius-docker-compose
RUN echo "NETWORK='$NETWORK'" > ./creator-node/.env
RUN echo "NETWORK='$NETWORK'" > ./discovery-provider/.env
