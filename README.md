# EdgeX Keeper

EdgeX Keeper is a lightweight configuration and registry service that is aimed to replace Consul in the EdgeX architecture.
It uses Redis from the existing EdgeX service architecture as the data persistent store, and implements the configuration and registry abstractions by adopting the `go-mod-configuration`
and `go-mod-registry` modules from EdgeX.

The code base of `edgex-core-keeper` is a clone of the [edgexfoundry/edgex-go](https://github.com/edgexfoundry/edgex-go) repository, and with the new `core-keeper` source code added. This provides an easier way to integrate the new code back into `edgexfoundry/edgex-go` in the future. 

>See EdgeX [Configuration and Registry Providers](https://docs.edgexfoundry.org/2.2/microservices/configuration/ConfigurationAndRegistry/) documentation for more details about the roles they play in the EdgeX architecture.

## High level architecture design
EdgeX Keeper acts as the Configuration and Registry Providers by using the **Config Handler** and **Registry Handler** modules to communicate with other EdgeX services.
Each EdgeX service will bootstrap its configuration and registry information into Keeper when started, listens for any **Writable** configuration change events via the message bus, and runs health checks to other services.

In addition, users can manage the configuration and registry information by calling the REST APIs provided by EdgeX Keeper.

The following architecture diagram demonstrates how Keeper works to function as configuration/registry providers for other services and communicate with external users.

![keeper_architecture_diagram](images/keeper_architecture.png?raw=true "keeper_architecture")

## Installation and deployment options
### Generate a binary executable and run
#### Prerequisites
##### Go
- The current targeted version of the Go language runtime for release artifacts is v1.17.x
- The minimum supported version of the Go language runtime is v1.17.x

##### Clone the edgex-core-keeper source code
Clone the source code and switch to the working directory.
```sh
git clone git@github.com:IOTechSystems/edgex-core-keeper.git
cd edgex-core-keeper
````

##### Redis
EdgeX Keeper is used Redis DB as the persistence layer. Please start the Redis server before running the EdgeX Keeper application.
###### Update the config file for Redis connection
In the `cmd/core-keeper/res/configuration.toml` file, enter the Redis connection detail in the following section.
```sh
[Databases]
  [Databases.Primary]
  Host = "edgex-redis" # Update the Redis db host for edgex-core-keeper to store data
  Name = "corekeeper"
  Port = 6379
  Timeout = 5000
  Type = "redisdb"
```

#### Build and run the binary executable
Follow the "Clone the edgex-core-keeper source code" and "Redis" steps above to clone the source and set up the Redis server that is used for connection.
```sh
make run
```


### Build your own Docker Container

In addition to running the services directly, Docker and Docker Compose can be used.

#### Prerequisites

See [the install instructions](https://docs.docker.com/install/) to learn how to obtain and install Docker.

#### Build
Follow the "Build and run the binary executable" steps above for obtaining the code, then build the dockcer image:
```sh
make docker 
```

#### Run the core-keeper image as a container

##### Use the docker container file in this repository to start the `edgex-core-keeper` and `edgex-redis` container
```sh
docker-compose up -d
```

## EdgeX Keeper API examples
After the EdgeX Keeper service is started, the `Get key` and `Create/Update key` endpoints can be invoked with the following examples.
### Create/Update key
There's a query parameter `flatten` defined in the Create/Update key API, and it controls the way to store the json object in the `value` field. The default value of `flatten` is `false`.
#### `flatten` is `true`
Enter the following command to create a series of keys starting with the prefix `core/data`, and the json object defined in the `value` field should be flattened before storing into database.
- Request sample
```shell
curl --location --request PUT 'localhost:59883/api/v2/kvs/key/core/data?flatten=true' \
--header 'Content-Type: application/json' \
--data-raw '{
  "value":{
    "MaxEventSize":25000,
    "Writable":{
      "PersistData":true,
      "LogLevel":"INFO",
      "InsecureSecrets":{
        "DB":{
          "path":"redisdb"
        }
      }
    }
  }
}'
```
- Response sample

The following individual keys were created with prefix `core/data` from the url path, and property name as key(ex. `MaxEventSize`) from the `value` json object in the request payload.
```shell
{
    "apiVersion": "v2",
    "statusCode": 200,
    "response": [
        "core/data/MaxEventSize",
        "core/data/Writable/PersistData",
        "core/data/Writable/LogLevel",
        "core/data/Writable/InsecureSecrets/DB/path"
    ]
}
```
#### `flatten` is `false`
Enter the following command to create a single key `core/data`, and the json object defined in the `value` field should be stored as a string.
- Request sample
```shell
curl --location --request PUT 'localhost:59883/api/v2/kvs/key/core/data?flatten=false' \
--header 'Content-Type: application/json' \
--data-raw '{
  "value":{
    "MaxEventSize":25000,
    "Writable":{
      "PersistData":true,
      "LogLevel":"INFO",
      "InsecureSecrets":{
        "DB":{
          "path":"redisdb"
        }
      }
    }
  }
}'
```
- Response sample

Only one key `core/data` has been created from the url path definition, and the json object defined in the `value` field is stored as a string.
```shell
{
    "apiVersion": "v2",
    "statusCode": 200,
    "response": [
        "core/data"
    ]
}
```

### Get key
#### Request example
Run the following command to get all the stored values of keys starting with the prefix `core/data`.
```shell
curl --location --request GET 'localhost:59883/api/v2/kvs/key/core/data'
```
#### Response example
```shell
{
    "apiVersion": "v2",
    "statusCode": 200,
    "response": [
        {
            "key": "core/data/MaxEventSize",
            "created": 1655372656274362,
            "modified": 1655372656274362,
            "value": 25000
        },
        {
            "key": "core/data/Writable/LogLevel",
            "created": 1655372656268428,
            "modified": 1655372656268428,
            "value": "SU5GTw=="
        },
        {
            "key": "core/data/Writable/InsecureSecrets/DB/path",
            "created": 1655372656269610,
            "modified": 1655372656269610,
            "value": "cmVkaXNkYg=="
        },
       ...
    ]
}
```

> #### Change summary to the key-value APIs:
> 
> - Update url path name from `/kv/{key}` to `/kvs/key/{key}`
> - Rename the GET API query parameters:
>   - `key` to `keyOnly`
>   - `raw` to `plaintext`
> - Rename the DELETE API query parameter:
>  - `recurse` to `prefixMatch`

## Use EdgeX Keeper as Configuration Provider
EdgeX Keeper uses the message bus from EdgeX go-mod-messaging module to publish the configuration change event. 
In order to update the `Writable` configuration settings for a service without restarting, EdgeX Keeper and other services need to connect to the same message bus.

Every time when a configuration value gets created or updated, Keeper will publish a key change event via the message bus; the other service which subscribes to the particular topic will get the updated configuration value.

The following diagram shows how the **Watch** feature is implemented on EdgeX Keeper.

![watch_diagram](images/watch.png?raw=true "watch")

### Steps for setting Keeper as Configuration Provider
Please see the following steps for EdgeX Core Data service to use EdgeX Keeper as Configuration Provider.

#### Update the`MessageQueue` of both EdgeX Keeper (Publisher) and Core Data (Subscriber)
As previously mentioned, EdgeX Keeper will publish the key change event whenever the configuration value is updated. In order for Core Data to
get notified when the `Writable` configuration setting is changed, EdgeX Keeper and Core Data need to connect to the same message bus.
Therefore, both EdgeX Keeper and Core Data have to connect to the same message bus, and the `MessageQueue` section in the configuration.toml of each should match up.

Here we use the default `MessageQueue` type `redis` as example (MQTT message bus type is supported as well).
See the following example defined in the configuration.toml:
```shell
[MessageQueue]
Protocol = "redis"
Host = "localhost"
Port = 6379
Type = "redis"
```

> The `MessageQueue` configuration setting is only defined in the configuration.toml of Core Data in the Kamakura release.
>
> For other EdgeX services, the `MessageQueue` section needs to be added for adopting the [go-mod-messaging](https://github.com/edgexfoundry/go-mod-messaging) module to subscribe to the message bus.
> 
> In the future EdgeX release, the message bus implementation will be applied to every EdgeX service to publish service metrics.
> Therefore, no additional tweak needed for each EdgeX service to use EdgeX Keeper as the configuration provider.



#### Replace go-mod-configuration module in edgex-go
Skip this step if using the source code to build Core Data in this repository directly.

If you would like to use the edgex-go source code to build Core Data from other repository,
add the following line in the end of [edgex-go/go.mod](https://github.com/edgexfoundry/edgex-go/blob/main/go.mod) file to replace `go-mod-configuration` module in to the one under [IOTechSystems](https://github.com/IOTechSystems/go-mod-configuration) which implements Keeper as the configuration provider option.
```shell
replace github.com/edgexfoundry/go-mod-configuration/v2 => github.com/IOTechSystems/go-mod-configuration/v2 core-keeper
```

#### Start EdgeX Keeper and Core Data
- Start EdgeX Keeper as mentioned in [Build and run the binary executable](https://github.com/IOTechSystems/edgex-core-keeper#build-and-run-the-binary-executable). 
- Start Core Data with the following command
    ```shell
    ./core-data -cp=keeper.http://localhost:59883
    ```
    or use the [Makefile](https://github.com/IOTechSystems/edgex-core-keeper/blob/main/Makefile) in this repository
    ```shell
    make run_core_data
    ```
    to build and run the binary executable of Core Data to use Keeper as configuration provider.
#### Invoke the EdgeX Keeper Update Key API to change the Writable setting
```shell
curl --location --request PUT 'localhost:59883/api/v2/kv/edgex/core/2.0/core-data/Writable/LogLevel' \
--header 'Content-Type: application/json' \
--data-raw '{
  "value": "DEBUG"
}'
```

#### Invoke the Core Data config API and verify the configuration has been updated without a restart
```shell
curl --location --request GET 'localhost:59880/api/v2/config'
```

## Use EdgeX Keeper as Registry
Please see the following steps for EdgeX Core Data service to use EdgeX Keeper as Registry.
### Set the `Registry` section in the configuration.toml
#### Update the`Registry` of Core Data
```shell
[Registry]
Host = "localhost"
Port = 59883
Type = "keeper"
```

#### Replace go-mod-registry module in edgex-go
Skip this step if using the source code to build Core Data in this repository directly.

If you would like to use the edgex-go source code to build Core Data from other repository,
add the following line in the end of [edgex-go/go.mod](https://github.com/edgexfoundry/edgex-go/blob/main/go.mod) file to replace `go-mod-registry` module in to the one under [IOTechSystems](https://github.com/IOTechSystems/go-mod-registry) which implements Keeper as the registry option.
```shell
replace github.com/edgexfoundry/go-mod-registry/v2 => github.com/IOTechSystems/go-mod-registry/v2 core-keeper
```

#### Start EdgeX Keeper and Core Data
- Start EdgeX Keeper as mentioned in [Build and run the binary executable](https://github.com/IOTechSystems/edgex-core-keeper#build-and-run-the-binary-executable).
- Start Core Data with the following command
    ```shell
    ./core-data --registry
    ```
  or use the [Makefile](https://github.com/IOTechSystems/edgex-core-keeper/blob/main/Makefile) in this repository
    ```shell
    make run_core_data
    ```
  to build and run the binary executable of Core Data to use Keeper as registry.
#### Invoke the EdgeX Keeper Registry API to check Core Data is registered
```shell
curl --location --request GET 'localhost:59883/api/v2/registry/serviceId/core-data' 
```
