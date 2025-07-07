# Simple Bank API
A simple bank API built with Go. Based on Anthony GG's [Building A JSON API In Golang](https://www.youtube.com/watch?v=pwZuNmAzaH8&list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2) series.
Uses:
* Gorilla Mux for routing
* JWT for authentication
    * Currently NOT secure, uses access token and refresh tokens
    * At the moment, access token only, plan to use rotating refresh tokens
    * But better: Use a third-party authentication service instead
* SQLite for storage
* SSE for real-time updates
* Simple vanilla JS for the frontend

## Requirements
* Go 1.22.5 (windows/amd64)
* [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) (for sqlite3 package)
* Optional: Docker for Postgres (if you want to use Postgres instead of SQLite)

## Backend Setup
### Locally
* Assuming you've installed everything required, you can use the commands in the Makefile
```sh
cd backend
make run
```

### Docker
```sh
cd backend
make docker-build
make docker-run
```

### Hosting
It is possible to host the server on Render, make sure to set a `CLIENT` and `JWT_SECRET` environment variable.

## Frontend Setup
### Locally
* Since it is vanilla JS and HTML, you can use the VSCode extension `Live Server` to serve the files

### Hosting
It is possible to host the frontend using Netlify, you can uploud the folder and it does the rest. 
