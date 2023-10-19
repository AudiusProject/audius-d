FROM docker:dind

ARG NETWORK=prod
ARG BRANCH=main

RUN apk add bash git curl libc-dev gcc python3 py3-pip python3-dev linux-headers

VOLUME /var/k8s/creator-node-db
VOLUME /var/k8s/mediorum
VOLUME /var/k8s/discovery-provider-db
VOLUME /var/k8s/discovery-provider-chain

WORKDIR /root
RUN git clone --single-branch --branch "$BRANCH" https://github.com/AudiusProject/audius-docker-compose.git ./audius-docker-compose

WORKDIR /root/audius-docker-compose
RUN echo "NETWORK='$NETWORK'" > ./creator-node/.env
RUN echo "NETWORK='$NETWORK'" > ./discovery-provider/.env

# docker volumes will initially create these as dirs if they don't exist
# create them here since this is a new audius-docker-compose clone
RUN touch discovery-provider/chain/spec.json
RUN touch discovery-provider/chain/static-nodes.json

RUN python3 -m pip install -r requirements.txt
RUN ln -sf $PWD/audius-cli /usr/local/bin/audius-cli
