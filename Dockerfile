FROM docker:dind

RUN apk add bash git

VOLUME /var/k8s/creator-node-db
VOLUME /var/k8s/mediorum

COPY ./audius /usr/local/bin/
RUN chmod +x /usr/local/bin/audius

WORKDIR /root
RUN git clone --single-branch --branch dev https://github.com/AudiusProject/audius-docker-compose.git ./audius-docker-compose

WORKDIR /root/audius-docker-compose/creator-node
COPY ./.env ./.env
# COPY ./override.env ./override.env
