# Compilers
## Windows 
* Issues can occur when working with the sqlite3 package
* What worked was using the TDM-GCC Compiler: https://jmeubank.github.io/tdm-gcc/
* GCC: gcc.exe (MinGW-W64 x86_64-ucrt-posix-seh, built by Brecht Sanders, r8) 13.2.0
* GO: go1.22.5 windows/amd64

## Linux (Ubuntu WSL)
* GCC: gcc (Ubuntu 11.4.0-1ubuntu1~22.04) 11.4.0
* GO: go1.22.5 linux/amd64

## SQLite
* Some operations differ from Postgres, such as how values are injected into the SQL string
* Or how to handle rows operations: https://stackoverflow.com/questions/32479071/sqlite3-error-database-is-locked-in-golang

# API
```
POST /login
```
* Request
    ```json
    {
        "number": 123456,
        "password": "password"
    }
    ```
* Response `200 OK`
    ```json
    {
        "token": "--jwt_token_with_payload--",
        "refresh_token": "--jwt_token--"
    }
    ```

```
POST /refresh
```
* Request
    ```json
    {
        "refresh_token": "--jwt_token--"
    }
    ```
* Response `200 OK`
    ```json
    {
        "token": "--jwt_token_with_payload--",
        "refresh_token": "--jwt_token--"
    }
    ```

```
GET /accounts 
```
* Response `200 OK`
    ```json
    [
        {
            "id": 1,
            "first_name": "Andrew",
            "last_name": "Levitt",
            "number": 123456,
            "balance": 1000000,
            "created_at": "2024-01-01T00:00:00Z"
        },
        {
            "id": 2,
            "first_name": "John",
            "last_name": "Doe",
            "number": 123457,
            "balance": 1000000,
            "created_at": "2024-01-01T00:00:00Z"
        }
    ]
    ```

```
POST /accounts 
```
* Request
    ```json
    {
        "first_name": "Andrew",
        "last_name": "Levitt",
        "password": "password"
    }
    ```
* Response `201 Created`
    ```json
    {
        "id": 1,
        "first_name": "Andrew",
        "last_name": "Levitt",
        "number": 123456,
        "balance": 1000000,
        "created_at": "2024-01-01T00:00:00Z"
    }
    ```

```
GET /account/{number}
```
* Response `200 OK`
    ```json
    {
        "id": 1,
        "first_name": "Andrew",
        "last_name": "Levitt",
        "number": 123456,
        "balance": 1000000,
        "created_at": "2024-01-01T00:00:00Z"
    }
    ```

```
DELETE /account/{number}
```
* Response `200 OK`
```
PUT /account/{number}
```
* Request
    ```json
    {
        "first_name": "Andrew",
        "last_name": "Levitt",
        "balance": 1000000
    }
    ```
*  Response `200 OK`
    ```json
    {
        "first_name": "Andrew",
        "last_name": "Levitt",
        "balance": 1000000
    }
    ```

```
POST /transfer/{number}
```
* Request
    ```json
    {
        "recipient": 123457,
        "amount": 1000000
    }
    ```
* Response `200 OK`
    ```json
    {
        "sender": 123456,
        "recipient": 123457,
        "amount": 1000000
    }
    ```


