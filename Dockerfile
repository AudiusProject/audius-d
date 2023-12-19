FROM docker:dind

ARG NETWORK=prod
ARG BRANCH=main

RUN apk add bash git curl libc-dev gcc python3 py3-pip python3-dev linux-headers postgresql-client redis

VOLUME /var/k8s/creator-node-db
VOLUME /var/k8s/mediorum
VOLUME /var/k8s/discovery-provider-db
VOLUME /var/k8s/discovery-provider-chain

WORKDIR /root
RUN git clone --single-branch --branch ${BRANCH} https://github.com/AudiusProject/audius-docker-compose.git ./audius-docker-compose

WORKDIR /root/audius-docker-compose
RUN echo "NETWORK='$NETWORK'" > ./creator-node/.env
RUN echo "NETWORK='$NETWORK'" > ./discovery-provider/.env
RUN echo "NETWORK='$NETWORK'" > ./identity-service/.env

RUN cp "./discovery-provider/chain/${NETWORK}_spec_template.json" "./discovery-provider/chain/spec.json"
RUN echo '[]' > ./discovery-provider/chain/static-nodes.json

RUN python3 -m venv .venv && source .venv/bin/activate && python3 -m pip install -r requirements.txt
RUN ln -sf $PWD/audius-cli /usr/local/bin/audius-cli

COPY daemon.json /etc/docker/daemon.json
