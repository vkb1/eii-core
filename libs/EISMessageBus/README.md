# EIS Message Bus

Message bus used between containers inside of EIS.

## Dependency Installation

To install the dependencies for the message bus execute the following command:

```sh
$ sudo -E ./install.sh
```

If you wish to compile the Python binding as well, then run the `install.sh`
script with the `--cython` flag (as shown below).

```sh
$ sudo -E ./install --cython
```

## Compilation

The EIS Message Bus utilizes CMake as the build tool for compiling the library.
The simplest sequence of commands for building the library are shown below.

```sh
$ mkdir build
$ cd build
$ cmake ..
$ make
```

This will compile only the C library for the EIS Message Bus. If you wish to
build with the Python binding, then specify the `WITH_PYTHON` flag when
executing the `cmake ` command (as shown below).

```sh
$ cmake -DWITH_PYTHON=ON ..
```

If you wish to include installation of the Go binding with the installation of
the EIS library, then specify the `WITH_GO` flag when executing the `cmake`
command (as shown below).

```sh
$ cmake -DWITH_GO=ON ..
```

Note that this only copies the Go binding library to your system's `$GOPATH`.
If you do not have your `$GOPATH` specified in your system's environmental
variables then an error will occur while executing the `cmake` command.

In addition to the `WITH_PYTHON` and `WITH_GO` flags, the EIS Message Bus
CMake files add flags for building the C examples and the unit tests associated
with the library. The table below specifies all of the available flags that can
be given to CMake for building the EIS Message Bus.

|       Flag      | Default |                                       Description                                    |
| :-------------: | :-----: | ------------------------------------------------------------------------------------ |
| `WITH_PYTHON`   | `OFF`   | If set to `ON`, specifies to build with the Python binding                           |
| `WITH_GO`       | `OFF`   | If set to `ON`, then CMake will install the Go binding during the installation phase |
| `WITH_TESTS`    | `OFF`   | If set to `ON`, builds the C unit tests with the EIS Message Bus compilation         |
| `WITH_EXAMPLES` | `OFF`   | If set to `ON`, then CMake will compile the C examples in addition to the library    |

> **NOTE:** These flags are in addition to any and all flags that are available
> for the `cmake` command. See the CMake documentation for additional flags.

If you wish to compile the EIS Message Bus in debug mode, then you can set the
the `CMAKE_BUILD_TYPE` to `Debug` when executing the `cmake` command (as shown
below).

```sh
$ cmake -DCMAKE_BUILD_TYPE=Debug ..
```

### Potential Compilation Issues

- **CMake Python Version Issue:** If CMake raises an error for finding the
    incorrect Python version, then add the following flag to the CMake
    command: `-DPYTHON_EXECUTABLE=/usr/bin/python3`
- **Python Binding Changes Not Compiling:** If a change is made to the Python
    binding and make has already been previously ran, then you must run
    `make clean` before running make again to compile the changes in the
    Python binding. This will need to be fixed later.

## Installation

If you wish to install the EIS Message Bus on your system, execute the
following command after building the library:

```sh
$ sudo make install
```

If you specified the `WITH_PYTHON=ON` option when executing CMake, then the
above command will also install the Python module for using the Python binding
for the EIS Message Bus.

## Running Unit Tests

> **NOTE:** The unit tests will only be compiled if the `WITH_TESTS=ON` option
> is specified when running CMake.

Execute one of the following commands from the `build/tests` folder to execute
the message bus unit tests.

```sh
$ ./msgbus-tests
$ ./msg-envelope-tests
$ ./crc32-tests
```

It is important to note that the `msgbus-tests` executable has an extra CLI
option which the other unit test binaries do not have. This option allows for
running the message bus tests to run over TCP rather than the IPC.

To run the message bus tests over TCP, execute the following command:

```sh
$ ./msgbus-tests --tcp
```

## Example Usage

> **IMPORTANT NOTE:** Some of the example configurations contain public/private
> keys for the purpose of show how to use the message bus with security enabled.
> THESE KEYS SHOULD **NEVER** BE USED IN PRODUCTION.

> **NOTE:** The examples will only be compiled if the `WITH_EXAMPLES=ON` option
> is set when CMake is executed during compilation.

All of the examples provided for the EIS Message Bus use a JSON configuration
file to configure the EIS Message Bus. There are several example configurations
provided with the message bus for running in IPC and TCP mode accross the
various different messaging patterns (i.e. Publish/Subscribe and Request/Response).
All of these example configurations are in the `examples/configs/` directory.
However, all of them are copied into the `build/examples/configs/` directory
as well when you build the message bus.

The table below specifies all of the provided example configurations.

|                 Configuration                   |                                     Description                                     |
| :---------------------------------------------: | ----------------------------------------------------------------------------------- |
| ipc_example_config.json                         | Configuration for IPC based communication. Works with all examples.                 |
| tcp_publisher_no_security.json                  | TCP configuration for publishing with no security.                                  |
| tcp_publisher_with_security_no_auth.json        | TCP configuration for publishing with key based auth without ZAP auth.              |
| tcp_publisher_with_security_with_auth.json      | TCP configuration for publishing with key based auth and ZAP auth.                  |
| tcp_subscriber_no_security.json                 | TCP configuration for subscribing to a topic with no security.                      |
| tcp_subscriber_with_security.json               | TCP configuration for subscribing to a topic with security enabled.                 |
| tcp_service_server_no_security.json             | TCP configuration for a service server side (i.e. `echo-service`) without security. |
| tcp_service_server_with_security_no_auth.json   | TCP configuration for a service server side with key based auth without ZAP auth.   |
| tcp_service_server_with_security_with_auth.json | TCP configuration for a service server side with key based auth and ZAP auth.       |
| tcp_service_client_no_security.json             | TCP configuration for a service client side (i.e. `echo-client`) with no security.  |
| tcp_service_client_with_security.json           | TCP configuration for a service client side with security enabled.                  |

You will notice that for the publisher configurations and service server side
configurations there are 3 configurations each, where as subscribers and service
client side configurations only have 2. This is because for publishers and service
server side applications there are two forms of security to enable: with ZAP authentication,
and no ZAP authentication. In the configurations with ZAP authentication, an
additional configuration value is provided which specifies the list of clients
(i.e. subscribers or service client side connections) which are allowed to connect
to the specified port. This list oporates as a whitelist of allowed client public
keys. If a connection is attempted with a key not in that list, then the connection
is denied.

### C Examples

There are currently 5 C examples:

1. `examples/publisher.c`
2. `examples/subscriber.c`
3. `examples/echo_server.c`
4. `examples/echo_client.c`
5. `examples/publisher_many.c`

All of the C example executables are in the `build/examples/` directory. To run
them, execute the following command:

```sh
$ ./publisher ./configs/ipc_example_config.json
```

> **NOTE:** The `tcp_example_config.json` can also be used in lieu of the IPC
> configuration file.

All of the examples follow the command structure above, i.e.
`<command> <json-config-file>.json`, except for the `publisher_many.c`
example. This example is explained more in-depth in the next section.

#### Publisher Many Example

The `examples/publisher_many.c` example serves as a reference for implementing
an application which contains many publishers. This also serves as a way of
testing this functionality in the EIS Message Bus.

The example can be run with the following command (from the `build/examples/`
directory):

```sh
$ ./publisher-many ./configs/ipc_example_config.json 5
```

In the case above, the example will create 5 publishers where the topic
strings follow the pattern `pub-{0,N-1}` where `N` is the number of publishers
specified through the CLI.

The behavior of how these topics are published depends on if the configuration
is IPC or TCP (i.e. if `type` is set to `zmq_ipc` vs. `zmq_tcp` in the JSON
configuration file).

If IPC communication is being used, then each topic will be a different Unix
socket file in the `socket_dir` directory specified in the configuration.

If TCP communication is being used, then each message will be published over
the `host` and `port` specified under the `zmq_tcp_publish` JSON object in the
configuration.

In order to subscribe to the topics published by this example, use the
`subscriber.c` example. If you are using TCP, then you will need to specify
the topic in your configuration. For example, your JSON configuration will
need to contain the following to subscribe to the `pub-0` topic:

```json
{
    "pub-0": {
        "host": "127.0.0.1",
        "port": 5569
    }
```
> **NOTE:** The `host` and `port` are assumed above, they may be different.

In order to simplify the creation of the configuration for subscribing to
topics over TCP, the `gen_tcp_sub_conf.py` helper script is provided. This
Python script will generate a JSON file for you based on your TCP JSON
configuration for the `publisher-many` example which contains all of the
topics specified so you can subscribe to any of them.

This helper script can be ran as follows:

```sh
$ python3.6 ./gen_tcp_sub_conf.py ./configs/tcp_publisher_no_security.json output.json 5
```

The command above uses the `tcp_publisher_no_security.json` for the `publisher-many`
configuration. Then it generates all 5 topics and outputs them into the
`output.json` file.

After generating this configuration, you can use the `subscribe.c` example as
shown below to subscribe to the `pub-1` topic:

```sh
$ ./subscribe output.json pub-1
```

### Python Examples

> **NOTE:** The Python examples will only be present if the `WITH_EXAMPLES=ON`
> and `WITH_PYTHON=ON` flags are set when CMake is executed during compilation.

There are currently 4 Python examples:

1. `python/examples/publisher.py`
2. `python/examples/subscriber.py`
3. `python/examples/echo_server.py`
4. `python/examples/echo_client.py`

To run the Python examples, go to the `build/examples/` directory. Then source
the `source.sh` script that is in the examples directory.

```sh
$ source ./source.sh
```

Then, execute one of the following commands:

```sh
$ python3 ./publisher.py ./configs/ipc_example_config.json
```

> **NOTE:** The `tcp_example_config.json` can also be used in lieu of the IPC
> configuration file.

All of the examples follow the same command structure as the `publisher.py`
script, i.e. `python3 <python-script>.py <json-config-file>.json`.

### Go Examples

> **IMPORANT NOTE:** It is assumed that when compiling the C library prior to
> running the examples that the `WITH_GO=ON` flag was specified when executing
> the `cmake` command. It is also assume that, `sudo make install` has been
> ran. If it has not and you do not wish to install the library, see the
> "Running Go Examples without Installing" section below.

When the `sudo make install` command is executed on your system, the Go binding
will be copied to your system's `$GOPATH`. To execute the examples provided
with the EIS Message Bus Go binding go to the `$GOPATH/src/EISMessageBus/examples`
directory on your system in a terminal window.

Once you are in this directory choose an example (i.e. publisher, subscriber,
etc.) and `cd` into that directory. Then, to run the example execute the
following command:

```sh
$ go run main.go -configFile <CONFIG-FILE>.json -topic publish_test
```

The example command above will run either the subscriber or publisher examples.
For the echo-client and echo-server examples the `-topic` flag should be
`-serviceName`.

Additionally, there are example configurations provided in the
`build/examples/configs/` directory after building the EIS Message Bus library.

### Running Go Examples without Installing

If you wish to run the Go binding examples with out installing the EIS Message
Bus library, then this can be accomplished by either copying or creating a
soft-link to the `go/EISMessageBus` directory in your `$GOPATH`. This can be
accomplished with one of the commands shown below.

```sh
$ cp -r go/EISMessageBus/ $GOPATH/src

# OR

$ ln -s go/EISMessageBus/ $GOPATH/src
```

> **NOTE:** The command above assumes that you are currently in the
> EISMessageBus source root directory.

Since it is assumed you have not ran the `sudo make install` command to install
the EIS Message Bus library, you must set the environmental variables specified
below prior to running the examples.

```sh
$ export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:$MSGBUS_DIR/build
$ export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$MSGBUS_DIR/build
```

Note that in the `export` commands above the `$MSGBUS_DIR` variable represents
the absolute path to the `libs/EISMessageBus` directory. It is very important
that this is the absolute path.

Once you have exported these variables, once you have done these steps, you can
run any of the Go examples as specified in the previous section.

## Security

> **IMPORTANT NOTE:** Security is only available for TCP communications. If IPC
> is being used, then all access must be controlled using Linux file
> permissions.

> **NOTE:** Example configurations using for enabling security in the examples
> are provided in the `examples` directory.

The ZeroMQ protocol for the EIS Message Bus enables to usage of
[CurveZMQ](http://curvezmq.org/) for encryption and authentication where the
[ZAP](https://rfc.zeromq.org/spec:27/ZAP/) protocol is used for the
authentication.

The ZeroMQ protocol for the message bus allows for using both CurveZMQ and ZAP
together, only CurveZMQ encryption, or no encryption/authentication for TCP
communication.

Enabling the security features is done through the configuration object which
is given to the `msgbus_initialize()` method. The example configurations below
showcase how to use the security features enabled in the message bus. It is
important to note that although the examples below use JSON to convey the
configurations it is not required that you use a JSON configuration for the
message bus. However, utilities are provided in the C library for the EIS
message bus for using a JSON file to configure the bus.

### Using Only CurveZMQ Encryption

If you wish to use the message bus with only CurveZMQ encryption, then you
specify the following keys for the communication types specified in the
sections below.

**IMPORTANT NOTE:** All keys must be Z85 encoded (see ZeroMQ documentation for
more information).

#### Publish/Subscribe

For publications over TCP, the configuration must contain a `server_secret_key`
value which the secret key of the Curve key pair that is Z85 encoded (see
the ZeroMQ documentation for more information).

Additionally, every subscriber configuration object (which is specified under
the key for the topic it is subscribing to) must contain the following three
keys: `server_public_key`, `client_public_key`, and `client_secret_key`.

**Example:**

Below is an example configuration in JSON (note: the keys are not Z85 encoded,
but are more clear text to help the example).

**Publisher Config:**

```json
{
    "type": "zmq_tcp",
    "zmq_tcp_publish": {
        "host": "127.0.0.1",
        "port": 3000,
        "server_secret_key": "publishers-secret-key"
    }
}
```

**Subscriber Config:**

```json
{
    "type": "zmq_tcp",
    "pub-sub-topic": {
        "host": "127.0.0.1",
        "port": 3000,
        "server_public_key": "publishers-public-key",
        "client_secret_key": "subscriber-secret-key",
        "client_public_key": "subscriber-public-key"
    }
}
```

In the example configurations above, it is assumed that the publisher is
sending messages on the `pub-sub-topic` topic.

#### Request/Response

For every service which is going to accept and respond to requests, there must
exist the `server_secret_key` in the configuration object for the service. The
key for the configuration of the service is its service name.

For every service which is going to issue requests to another service, there
must exist a configuration object for the destination service name which
contains the following three keys: `server_public_key`, `client_public_key`,
and `client_secret_key`.

**Example:**

**Service Config:**

```json
{
    "type": "zmq_tcp",
    "example-service": {
        "host": "127.0.0.1",
        "port": 3000,
        "server_secret_key": "service-secret-key"
    }
}
```

**Service Requester Config:**

```json
{
    "type": "zmq_tcp",
    "example-service": {
        "host": "127.0.0.1",
        "port": 3000,
        "server_public_key": "service-public-key",
        "client_secret_key": "service-requester-secret_key",
        "client_public_key": "service-requester-public-key"
    }
}
```

In the example above, the service requester will connect to the `example-service`
and issue requests to it on the port: `127.0.0.1:3000`.

### Using ZAP Authentication

To enable ZAP authentication protocol using CurveZMQ on top of the encryption,
then in the configuration specify the key `allowed_clients`. This key must have
a value which is a list of Z85 encoded strings which are the public keys of the
clients which are allowed to connect to the application.

For example, using the publish/subscribe example from before, to make it so
that only the subscriber client can connect to the publisher the publisher's
configuration would be modified to be the following:

```json
{
    "type": "zmq_tcp",
    "allowed_clients": ["subscriber-public-key"],
    "zmq_tcp_publish": {
        "host": "127.0.0.1",
        "port": 3000,
        "server_secret_key": "publishers-secret-key"
    }
}
```

### Disabling Security

To disable all encryption and authentication for TCP communication do not
specify any of the configuration keys documented above. This will cause the
message bus to initialize the ZeroMQ protocol without any of the CurveZMQ
security primitives.