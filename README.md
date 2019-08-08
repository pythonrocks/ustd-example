Omnicored Test Setup
--------------------

This project represents development setup of the omnicored node with complementing rpc/rest API on top of it. Please note, omnicore runs in test network mode.

To start project run following commands:

```
$ docker-compose up -d
```

Since api starts in generic golang container, you need to wait until API binary is built.
Type following to check container logs:

```
$ docker-compose logs -f api
```

When you'll see following message:

```
Starting RPC server on port 5001
```
project is up and running.

You can use following endpoints to test it.

To get wallet info, type:

```
$ curl localhost:8000/api/v1/wallet
```

To see your wallet addresses:

```
$ curl localhost:8000/api/v1/address
```

To get some address balance:

```
$ curl localhost:8000/api/v1/address/<address>
```

To create new address:

```
$ curl -XPOST localhost:8000/api/v1/address
```

To get transaction info:

```
$ curl localhost:8000/api/v1/transaction/<txID>
```

To send currency to another address:

```
$ curl -XPOST localhost:8000/api/v1/transaction -d '{"fromaddress": "<your address>", "toaddress": "<address to send currency to>", "propertyid": "<property id>", "amount": "<amount>"}'
```

This example uses Envoy is used to provide REST-GRPC transcoding, service itsself is GRPC based. Please check `service.proto` to see API definition. In order to generate new descriptor use command:

```
make descriptor
```

To generate GRPC stubs use:

```
make interface
```
