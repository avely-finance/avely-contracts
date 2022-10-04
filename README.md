# Avely Finance :: Avely Contracts

## Docker Compose Setup

* `docker-compose run --rm runner` to start bash with scilla and local blockchain
* `docker-compose up devex` to start the [zilliqa blockchain explorer](https://github.com/Zilliqa/devex)

If you see auth error, read [this](https://github.community/t/docker-pull-from-public-github-package-registry-fail-with-no-basic-auth-credentials-error/16358/90)

## Configuration

For local development you'll need .env.local file, copied from .env.example.

## Docs

* [Tools](docs/tools.md)
* [How to debug Scilla](docs/debug.md)
* [Isolated Server know how](docs/isolated_server.md)
