# aaa

| URI                       | Methods| Description           | Params
|:--------------------------|:-------|:--------------------- |:-------
| /users                    | Get    | get all users         | 
| /users/{id}               | Get    | get user by id        |
| /users/register           | POST   | add new user          | post
| /users/auth               | POST   | user authentication   | post{ login : "aaa", password : "aaa"}
| /users/{login}            | DELETE | user delete           | 
| /users/change_passwd      | PUT    | user passwd change    | post{ login : "aaa", old_password : "aaa", new_password : "aaa"}


#-----      GRPC       -----

| URI                       | Methods
|:--------------------------|:-----------------------|
| /user/register            | RegistrationUserEndPoint(AuthUserRequest) returns (AuthUserResponse) {}  port : 8447
| /users/auth               | AuthUserEndPoint(AuthUserRequest) returns (AuthUserResponse)  {}   port : 8447
| /backend/auth             | AuthBackendEndpoint(AuthRequest) returns (AuthResponse) {}   port : 8448
| :-------------------------|:-----------------------| 
| /ibis/auth                | AuthEndpoint(stream pb.AAA_AuthEndpointServer)  {}  returns (stream AuthResponse) port : 8443
| /ibis/session             | SessionEndpoint(stream pb.AAA_SessionEndpointServer) {} returns (stream SessionResponse) port : 8443


## Run aaa server

```
$ cd cmd/
$ go run main.go
```

## Run backend client 

```
$ cd cmd/
$ go run srv_client.go
```

## Ports 

```
userPort    = 8447
backendPort = 8448
ibisPort    = 8443

```


## Proto files 

```
$ cd endpoint/

```
