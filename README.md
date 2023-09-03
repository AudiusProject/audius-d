# dot-slash-audius

because running audius should look like this
```
./audius
```

## run

on an x86 linux machine that has docker installed
```
curl -o audius https://raw.githubusercontent.com/AudiusProject/dot-slash-audius/main/audius && chmod +x audius
```

minimal required config, can mount at this location or pass via `-c` flag at runtime
```
# ~/.audius/audius.conf

# content node
creatorNodeEndpoint=
delegateOwnerWallet=
delegatePrivateKey=
spOwnerWallet=

# discovery node
audius_discprov_url=
audius_delegate_owner_wallet=
audius_delegate_private_key=

```

run
```
./audius [-c audius.conf]
```

## build

```
make
```

## todo

- docker buildx manifests for multi arch
- fix discovery chain
