## Run App ##
```plaintext
go run main.go
```

## Generate JWT Secret ##
```plaintext
go run generate_secret.go
```

## Structure Base ##
```plaintext
Project/
├── config/
│   └── config.go
├── controller/
│   ├── auth_controller.go
│   ├── role_controller.go
│   ├── secret_controller.go
│   └── user_controller.go
├── middleware/
│   └── auth_middleware.go
├── model/
│   ├── init.go
│   ├── role_model.go
│   └── user_model.go 
├── route/
│   └── routes.go
├── utils/
│   ├── api_response_helper.go
│   ├── blacklist_helper.go
│   ├── hash_helper.go
│   ├── input_validation_helper.go
│   └── jwt_helper.go 
├── .env-example
├── generate_secret.go
├── go.mod
├── go.sum
└── main.go
  ```

## Library ##
- gin
- golang-jwt
- godotenv
- crypto
- postgres
- dotenv
- gorm