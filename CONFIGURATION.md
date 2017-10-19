# Configuration for multimissile

A configuration file format for multimissile is [TOML](https://github.com/toml-lang/toml).

A configuration for multimissile has some sections. A example is [here](config/example.toml).

 * [Global Section](#core-section)
 * [Endpoints Section](#endpoints-section)

## Global Section

|name               |type  |description                                 |default         |note                                                              |
|-------------------|------|--------------------------------------------|----------------|------------------------------------------------------------------|
|Port               |string|port number or unix socket path             |29300           |e.g.)29300, unix:/tmp/msl.sock <br/> `-p` option can overwrite    |
|LogLevel           |string|log-level                                   |error           |                                                                  |
|Timeout            |int   |timeout for proxying request                 |5               |unit is second                                                    |
|MaxIdleConnsPerHost|int   |maximum idle to keep per-host               |100             |                                                                  |
|DisableCompression |bool  |delete `Accept-Encoding: gzip` in header    |false           |                                                                  |

## Endpoints Section

|name           |type          |description                        |default|note|
|---------------|--------------|-----------------------------------|-------|----|
|Name           |string        |Endpoint name                      |       |    |
|To             |string        |Endpoint URL                       |       |    |
|ProxySetHeaders|array of array|Headers appended on proxying request|       |    |

As a scheme, **http** and **https** are available for **To**.

```
To = "http://example.com"
# or
To = "https://example.com"
```

If a scheme is not specified, **http** is used.

```
To = "example.com"
```

* example.com

# About API

See [SPEC.md](SPEC.md) about details for APIs.
