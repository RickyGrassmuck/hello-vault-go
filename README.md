# hello-vault-go

This is a sample application that demonstrates how to authenticate to and
retrieve secrets from HashiCorp [Vault][vault].

## Prerequisites

1. [`vault`][vault] with LDAP preconfigured
1. [`curl`][curl] to test our endpoints


## Try it out
### 1. Bring up the services

### 2. Try out `POST /payments` endpoint (static secrets workflow)

`POST /payments` endpoint is a simple example of the static secrets workflow.
Our service will make a request to another service's restricted API endpoint
using an API key value stored in Vault's static secrets engine.

```bash
curl -s -X POST http://localhost:8080/payments | jq
```

```json
{
  "message": "hello world!"
}
```

Check the logs:

```bash
docker logs hello-vault-go-app-1
```

```log
...
2021/12/10 23:20:36 getting secret api key from vault
2021/12/10 23:20:36 getting secret api key from vault: success!
[GIN] 2021/12/10 - 23:20:36 | 200 |    3.219167ms |    192.168.96.1 | POST     "/payments"
```

### 3. Examine the logs for renew logic

One of the complexities of dealing with short-lived secrets is that they must be
renewed periodically. This includes authentication tokens and database
credentials.

Examine the logs for how the Vault auth token is periodically renewed:

## Stack Design

### API

| Endpoint             | Description                                                            |
| -------------------- | ---------------------------------------------------------------------- |
| **POST** `/payments` | A simple example of Vault static secrets workflow (see example above)  |
