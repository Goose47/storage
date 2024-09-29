# key-value persistent storage

## Description

This project provides a key-value persistent containerized microservice with mongo.

It allows to set, retrieve and delete
various files. Files are removed after specified ttl if present. Files are stored in local filesystem but it is easy to move to a remote storage.

Requests are authenticated with [grpc sso](https://github.com/Goose47/go-grpc-sso.protos/tree/main). App checks whether user is
authenticated and then performs gRPC request so SSO server to check users permissions.

## Prerequisites

Running gRPC SSO application is required. App must be registered in SSO. JWT tokens are obtained via sso mentioned above.

## Steps to install

- Configure .env files, located in ./mongo/.env and ./api/.env based on corresponding .env.example files
- run ``docker compose up``

## API Description

Token must be present in Authorization header:
```
Authorization: Bearer {jwt-token-here}
```

------------------------------------------------------------------------------------------

#### Upload a file and save it by :key parameter

<details>
 <summary><code>POST</code> <code><b>/api/v1/storage/:key</b></code> <code>(accepts multipart/form-data and overwrites previously set key)</code></summary>

##### Parameters

> | name | type     | data type | description                                                                                                        |
> |------|----------|-----------|--------------------------------------------------------------------------------------------------------------------|
> | file | required | file      | A file to be stored: photo, audio, video, etc                                                                      |
> | ttl  | optional | int       | Key will be removed and file will be deleted after after now + ttl. A default value of 0 will set ttl to infinity. |


##### Responses

> | http code | content-type                      | response                                |
> |-----------|-----------------------------------|-----------------------------------------|
> | `200`     | `application/json; charset=utf-8` | `{"message": "Set :key"}`               |
> | `422`     | `application/json; charset=utf-8` | `{"message": "Error description here"}` |

</details>

------------------------------------------------------------------------------------------

#### Retrieve file by :key parameter

<details>
 <summary><code>GET</code> <code><b>/api/v1/storage/:key</b></code></summary>

##### Responses

> | http code | content-type                      | response                           |
> |-----------|-----------------------------------|------------------------------------|
> | `200`     | `depend on file mime`             | `{"message": "Set :key"}`          |
> | `404`     | `application/json; charset=utf-8` | `{"message": ":key is not found"}` |

</details>

------------------------------------------------------------------------------------------

#### Delete file by :key parameter

<details>
 <summary><code>DELETE</code> <code><b>/api/v1/storage/:key</b></code></summary>

##### Responses

> | http code | content-type                       | response                           |
> |-----------|------------------------------------|------------------------------------|
> | `200`     | `application/json; charset=utf-8`  | `{"message": "ok"}`                |
> | `404`     | `application/json; charset=utf-8`  | `{"message": ":key is not found"}` |

</details>

------------------------------------------------------------------------------------------