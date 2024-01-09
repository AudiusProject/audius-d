FROM docker:dind

RUN apk add bash git curl libc-dev gcc python3 py3-pip python3-dev linux-headers postgresql-client redis

VOLUME /var/k8s/creator-node-db
VOLUME /var/k8s/mediorum
VOLUME /var/k8s/discovery-provider-db
VOLUME /var/k8s/discovery-provider-chain

WORKDIR /root
RUN git clone https://github.com/AudiusProject/audius-docker-compose.git ./audius-docker-compose

WORKDIR /root/audius-docker-compose

RUN python3 -m venv .venv && source .venv/bin/activate && python3 -m pip install -r requirements.txt
RUN ln -sf $PWD/audius-cli /usr/local/bin/audius-cli

COPY daemon.json /etc/docker/dahttps://github.com/AudiusProject/auhttps://github.com/AudiusProject/audius-d/blob/main/.circleci/config.yml#L73dius-d/blob/main/.circleci/config.yml#L73emon.json
