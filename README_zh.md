# MINIO-Cache-Web-Server

[DOC](README.md) | [文档](README_zh.md)

基于minio和go-gin，并加入缓存和jwt的web服务器，提供可靠的资源服务。

## 快速使用

在使用前，需要根据用户获取一个token
```
GET /login 

Header:Content-Type=application/json // data format
Data:
{
    "UserName":"test", // 用户名
    "UserKey":"testKey", // 用户key
    "ServiceName":"test" // 用于服务名
}

Response:
{
    "code": 200,
    "message": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyIiOiJvZG0iLCJleHAiOjE2NDAzMzg4NTgsIm9yaWdfaWF0IjoxNjQwMzM1MjU4fQ.AFvLERGMAkI5ht5PX9EwznrEBDZtB2WDi-nuGAvX8yE"
}
```
token加入鉴权方法的header即可。  
**注意，token前要加'Bearer '**  
获取文件`getObject`.
```
GET /auth/getObject?bucket={bucketName}&name={fileName}

Header:Authorization=Bearer eyJhbGciOiJIUzI...

Response::FileStream
```
上传文件`uploadObject`.
```
POST /auth/postObject?bucket={bucketName}

Header:Authorization=Bearer eyJhbGciOiJIUzI...
       Content-Type=multipart/form-data
Form-Data:select your file

Response:
{
    "message": {fileSize}
}
```
## Configuration

## API

### /login
### /auth/getObject
### /auth/postObject