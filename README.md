# Golang_JWT_Project
Backend with authentication and validation that maintains and addresses a MongoDB Database.

C:.
│   .env
│   go.mod
│   go.sum
│   main.go
│   
├───controllers
│       userController.go
│       
├───database
│       databaseConnection.go
│       
├───helpers
│       AuthHelper.go
│       TokenHelper.go
│
├───middleware
│       authMiddleware.go
│
├───models
│       userModel.go
│
└───routes
        authRouter.go
        userRouter.go

        
.env file of the format
PORT=[..]
MONGODB_URL=[...]
SECRET_KEY="xyz"
