# zillica-dapp-template

## Requirements

* nodejs, yarn https://classic.yarnpkg.com/en/docs/install
* `yarn set version berry`
* `yarn install`
* `docker pull zilliqa/scilla:latest` (3.5Gb)
* `docker pull zilliqa/zillica-isolated-server:latest` (1.5Gb)
* `docker pull zilliqa/devex:latest` (37Mb)

## Yarn scripts
* `yarn run start:isolatedserver` run/start docker container with [zilliqa isolated server](https://github.com/Zilliqa/zilliqa-isolated-server)
* `yarn run start:devex` run/start docker container with [zilliqa blockchain explorer](https://github.com/Zilliqa/devex)
* `yarn run typecheck:helloWorld` to typecheck contracts/helloWorld.scilla contract
* `yarn run typecheck:all` to typecheck all contracts (not implemented for now)

## Docker Compose Setup
* `docker-compose up devex` to start the [zilliqa blockchain explorer](https://github.com/Zilliqa/devex)
* `docker-compose run --rm runner` to start bash with NodeJS, yarn, scilla, zli inside

If you see auth error, read [this](https://github.community/t/docker-pull-from-public-github-package-registry-fail-with-no-basic-auth-credentials-error/16358/90)

## How to deploy and test contracts on local blockchain
1. [**Zilliqa-JS SDK**](https://github.com/Zilliqa/Zilliqa-JavaScript-Library):
    run test script from this repo `yarn node tests/basic.ts`
2. to be continued...

## Other ways to query local isolated server
* through **Devex** web-interface (block explorer) http://localhost:5555
* through **zli** command-line utility [zli is a command line tool based on the Zilliqa Golang SDK](https://github.com/Zilliqa/zli)
    ```
    $ zli rpc balance --api http://localhost:5555 -a d90f2e538ce0df89c8273cad3b63ec44a3c4ed82
    cannot load wallet =  open /home/mx/.zilliqa: no such file or directory
    {"balance":"90000000000000000000000","nonce":0}
    ```
* any **REST-client** like [Insomnia](https://insomnia.rest)
```
> POST / HTTP/1.1
> Host: localhost:5555
> User-Agent: insomnia/2021.5.3
> Content-Type: application/json
> Accept: */*
> Content-Length: 100

{"id":1,"jsonrpc":"2.0","method":"GetBalance","params":["d90f2e538ce0df89c8273cad3b63ec44a3c4ed82"]

< HTTP/1.1 200 OK
< Connection: Keep-Alive
< Content-Length: 81
< Access-Control-Allow-Origin: *
< Content-Type: application/json
< Date: Mon, 11 Oct 2021 14:18:58 GMT

{
  "id": 1,
  "jsonrpc": "2.0",
  "result": {
    "balance": "90000000000000000000000",
    "nonce": 0
  }
}
```

## Predefined zilliqa accounts
https://github.com/Zilliqa/zilliqa-isolated-server/blob/master/boot.json
```
{
    "d90f2e538ce0df89c8273cad3b63ec44a3c4ed82": {
        "privateKey": "e53d1c3edaffc7a7bab5418eb836cf75819a82872b4a1a0f1c7fcf5c3e020b89",
        "amount": "90000000000000000000000",
        "nonce": 0
    },
    "381f4008505e940ad7681ec3468a719060caf796": {
        "privateKey": "d96e9eb5b782a80ea153c937fa83e5948485fbfc8b7e7c069d7b914dbc350aba",
        "amount": "90000000000000000000000",
        "nonce": 0
    },
    "b028055ea3bc78d759d10663da40d171dec992aa": {
        "privateKey": "e7f59a4beb997a02a13e0d5e025b39a6f0adc64d37bb1e6a849a4863b4680411",
        "amount": "90000000000000000000000",
        "nonce": 0
    },
    "f6dad9e193fa2959a849b81caf9cb6ecde466771": {
        "privateKey": "589417286a3213dceb37f8f89bd164c3505a4cec9200c61f7c6db13a30a71b45",
        "amount": "90000000000000000000000",
        "nonce": 0
    },
    "10200e3da08ee88729469d6eabc055cb225821e7": {
        "privateKey": "5430365143ce0154b682301d0ab731897221906a7054bbf5bd83c7663a6cbc40",
        "amount": "1000000000000000000",
        "nonce": 0
    },
    "ac941274c3b6a50203cc5e7939b7dad9f32a0c12": {
        "privateKey": "1080d2cca18ace8225354ac021f9977404cee46f1d12e9981af8c36322eac1a4",
        "amount": "1000000000000000000",
        "nonce": 0
    },
    "ec902fe17d90203d0bddd943d97b29576ece3177": {
        "privateKey": "254d9924fc1dcdca44ce92d80255c6a0bb690f867abde80e626fbfef4d357004",
        "amount": "1000000000000000000",
        "nonce": 0
    },
    "c2035715831ab100ec42e562ce341b834bed1f4c": {
        "privateKey": "b8fc4e270594d87d3f728d0873a38fb0896ea83bd6f96b4f3c9ff0a29122efe4",
        "amount": "1000000000000000000",
        "nonce": 0
    },
    "6cd3667ba79310837e33f0aecbe13688a6cbca32": {
        "privateKey": "b87f4ba7dcd6e60f2cca8352c89904e3993c5b2b0b608d255002edcda6374de4",
        "amount": "1000000000000000000",
        "nonce": 0
    }
}
```
