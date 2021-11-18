# Docker

Official docker images are provided through docker hub here: https://hub.docker.com/r/superseriousbusiness/gotosocial

TODO: currently the docker images only support the amd64 architechture. In the future we will add support for more architechtures. For now if you wish to run GoToSocial on, for example, a Raspberry Pi, you may [build the docker image yourself](https://github.com/ForestJohnson/gotosocial/compare/main...forest-jank-multi-arch-docker-build)

GoToSocial can be configured using [Environment Variables](../index.md#environment-variables) if you wish, allowing your GoToSocial configuration to be embedded inside your docker container configuration.

Assuming you will be using something like [docker compose](https://docs.docker.com/compose/) to configure your GoToSocial docker container,
you might write a `docker-compose.yml stanza something like this:

```
version: "3.3"
services:
  gotosocial:
    image: superseriousbusiness/gotosocial:0.1.0
    restart: always
    volumes:
      - type: bind
        source: ./gotosocial/storage
        target: /gotosocial/storage
    environment:
      GTS_PORT: '8080'
      GTS_PROTOCOL: 'https'
      GTS_TRUSTED_PROXIES: '0.0.0.0/0'
      GTS_HOST: 'gotosocial.example.com'
      GTS_ACCOUNT_DOMAIN: 'gotosocial.example.com'
      GTS_DB_TYPE: 'sqlite'
      GTS_DB_ADDRESS: '/gotosocial/storage/database/sqlite.db'
      GTS_STORAGE_SERVE_PROTOCOL: 'https'
      GTS_STORAGE_SERVE_HOST: 'gotosocial.example.com'
      GTS_STORAGE_SERVE_BASE_PATH: '/media'
      GTS_LETS_ENCRYPT_ENABLED: 'true'
      GTS_LETS_ENCRYPT_PORT: '1213'
      GTS_LETS_ENCRYPT_CERT_DIR: '/gotosocial/storage/certs'
      GTS_LETS_ENCRYPT_EMAIL_ADDRESS: 'admin@example.com'
```

Note how the secret value SMTP_PASSWORD is templated into the docker-compose file. You can store the actual password in a separate file named `.env` like so 