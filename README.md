# Avely Finance :: Avely Contracts

## Docker Compose Setup

- setup Zilliqa's Isolated Server as [described](https://github.com/Zilliqa/zilliqa-developer?tab=readme-ov-file#building-and-running-docker-images)
- `docker-compose run --rm runner` to start bash with scilla and local blockchain
- `docker-compose up devex` to start the [zilliqa blockchain explorer](https://github.com/Zilliqa/devex)

If you see auth error, read [this](https://github.community/t/docker-pull-from-public-github-package-registry-fail-with-no-basic-auth-credentials-error/16358/90)

## Configuration

For local development you'll need .env.local file, copied from .env.example.

## Docs

- [Isolated Server know how](docs/isolated_server.md)

## License Information

Please note that the following files and/or directories are subject to a restricted license:

- [contracts/stzil.scilla]
- [contracts/holder.scilla]
- [contracts/holder.scilla]
- [sdk]

These files and directories are protected by a separate license, which only allows viewing of the code. You may not use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of this code without the express written permission of the author. For more details, please refer to the RESTRICTED_LICENSE file in the root of this repository.
