in this golang project i tried to setup, i created a handler and make the server run 'go run cmd/api/main.go' and made http request

```
curl --location 'http://localhost:8080/api/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"john",
    "acceptedTermCondition":true,
    "email":"dljfl@gmail.com",
    "password":"pdj#lds@AS2"
}'
```

but i got this error

```
"error": "invalid JSON: json: Unmarshal(non-pointer models.CreateUserRequest)"
```

can you help me fix this error?

```

```
