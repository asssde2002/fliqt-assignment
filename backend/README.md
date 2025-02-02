# fliqt-assignment

`cd backend`

`go mod init backend`

`go get github.com/gin-gonic/gin@v1.10.0`

`go get gorm.io/gorm@v1.25.12`

`go get gorm.io/driver/mysql`

`go get github.com/gin-contrib/cors@v1.7.3`

`got get github.com/joho/godotenv`

`go get github.com/githubnemo/CompileDaemon@latest`



* Run server with hot reload
`CompileDaemon -build="go build -o app ./cmd/server" -command="./app"`