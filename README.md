
# JWT Authentication and Validation in Go

A Backend Application with authentication and validation that maintains and addresses a MongoDB Database.
## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [License](#license)
4. [Contact](#contact)

## Introduction - Features
Password Encryption - implements hashing of user passwords.

Maintain access points (User or Admin) - special access privilege using tokens.

Last updated time, created-at time, token, refresh token,userId etc.

Special access of all available users to the Admin.

Only the user themselves or the admin can access Logging in to access the user data.

## Installation

### Prerequisites
List any software, libraries, or tools that need to be installed before using your project.

```bash
# Example for Go installation
$ go version
```

### Installing the Project
Provide step-by-step instructions on how to install and set up the project.

```bash
# Clone the repository
$ git clone https://github.com/ShabnamHaque/go-jwt.git

# Navigate to the project directory
$ cd go-jwt

# Install dependencies
$ go mod tidy

# Run the project
$ go run main.go
```go
/*
### File Structure
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
  

### Configuration
PORT=[..]
MONGODB_URL=[...]
SECRET_KEY="xyz"

```

<!-- 
## API Reference
Document the project's API, including all public methods, structures, and interfaces. Provide a description, input parameters, and return values for each.
 -->


## License
```
This project is licensed under the MIT License.
```

## Contact
Mail - shabnamhaque20@gmail.com
