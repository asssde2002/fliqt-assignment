# fliqt-assignment

`cd backend`

`go mod init backend`

`go get github.com/gin-gonic/gin@v1.10.0`

`go get gorm.io/gorm@v1.25.12`

`go get gorm.io/driver/mysql`

`go get github.com/gin-contrib/cors@v1.7.3`

`got get github.com/joho/godotenv`

`go get github.com/githubnemo/CompileDaemon@latest`

`go get go get github.com/golang-jwt/jwt/v5@v5.2.1`



* Run server with hot reload
`CompileDaemon -build="go build -o app ./cmd/server" -command="./app"`