# MINIO-Cache-Web-Server

[DOC](README.md) | [文档](README_zh.md) 

A minio web application.  
Based on minio and go-gin, adding the cache and jwt for the web server, now we get a useful resource server.

## How to use

*Before we get a file, we need to sign a new token.*
```
GET /login 

Header:Content-Type=application/json // data format
Data:
{
    "UserName":"odm", // for the user which is allowed to access
    "UserKey":"odmKey", // an authorization key
    "ServiceName":"test" // for what service in need
}

Response:
{
    "code": 200,
    "message": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyIiOiJvZG0iLCJleHAiOjE2NDAzMzg4NTgsIm9yaWdfaWF0IjoxNjQwMzM1MjU4fQ.AFvLERGMAkI5ht5PX9EwznrEBDZtB2WDi-nuGAvX8yE"
}
```
*Now we get a token to use other methods.*
**Note: in Header Authorization, token need add 'Bearer ' at the first**  
*Get file by the method `getObject`.*
```
GET /auth/getObject?bucket={bucketName}&name={fileName}

Header:Authorization=Bearer eyJhbGciOiJIUzI...

Response::FileStream
```
*Upload file by the method `uploadObject`.*
```
POST /auth/putObject?bucket={bucketName}

Header:Authorization=Bearer eyJhbGciOiJIUzI...
       Content-Type=multipart/form-data
Form-Data:select your file

Response:
{
    "message": {fileSize}
}
```
*In addition, add `/quick` to use the cache accessing the method.*
* /auth/quick/getObject
## Configuration
[minio]  
[gin]  
[token]  
[key]  
`user group`
## API

### POST   /login
### GET    /auth/quick/getBucketList
### GET    /auth/quick/getObject
### GET    /auth/refresh_token
### GET    /auth/getObject
### POST   /auth/puttObject
