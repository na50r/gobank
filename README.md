# Simple Bank API
A simple bank API built with Go. Based on Anthony GG's [Building A JSON API In Golang](https://www.youtube.com/watch?v=pwZuNmAzaH8&list=PL0xRBLFXXsP6nudFDqMXzrvQCZrxSOm-2) series. This is a practice repositoriy, so it has also other functionalities, such as:
* A game where you combine elements to create new ones (InfiniteCraft)
* Real-time updates of transactions; so each time money is transferred, the account page of the recipient is updated automatically (uses Server-Sent Events)

https://github.com/user-attachments/assets/c955ee04-d720-4267-8fc9-f2f2dec1f279

Uses:
* Gorilla Mux for routing
* JWT for authentication
    * NOT SECURE! Uses rotating refresh tokens but recommended to use third-party options instead
* SQLite for storage
* SSE for real-time updates
* Simple vanilla JS for the frontend

## Requirements
* Go 1.22.5 (windows/amd64)
* [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) (for sqlite3 package)
* Optional: Docker for Postgres (if you want to use Postgres instead of SQLite)

# Additional Files
* You need a folder named `icons' in the backend folder that contains pngs (I used [these](https://github.com/wayou/anonymous-animals))
* You need a file named `Recipes.csv` in the backend folder that contains the recipes for the elements (I used [this](https://docs.google.com/spreadsheets/d/14aPnIQt252SYvmjP7iKQM3_rkLvqFntcqmxZNwq_tkI/edit?gid=0#gid=0))
    * Note that `Recipes.csv' should adhere to the following format:
    ```csv
    A,B,Result
    Fire,Water,Steam
    ...
    ```

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
