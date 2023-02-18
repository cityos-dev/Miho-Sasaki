# Specification of this video storage server

## Overview
- language
  - Golang
- Libraries
  - [xorm.io/xorm](https://xorm.io/)
    - ORM library for communicating with DB.
    - Used this to make it easier to Get/Insert/Delete operations.
  - [gin](https://github.com/gin-gonic/gin)
    - Web framework for Go.
    - Easy to use.
    - Easier to response Video Data.
  - [mysql](github.com/go-sql-driver/mysql)
    -  Driver for communicating with mysql.


## Architecture
I define an architecture of this application like the below directory structure.

The purpose that I separate the directories like this is to define a responsibility for each operation.
In each file, the interfaces are defined to be called by other layers. 
So that the other layer doesn't need to know the implementation inside the layer.

Also, there is a merit that we can easily convert to new library/db.

```
├─┬ contents/
| └─ video/
├─ helpers/
├─ infra/
├─ service/
├─ testing/
├─ docker-compose.yml
├─ Dockerfile
├─ README.md
├─ README_SPECIFICATION.md
├─ main.go
└─ testing/
```

### contents/
- Path to store video and other files.
### helpers/
- Store errors and other functions to be used inside of the application.
### infra/
- Store db related functions.
- Currently mysql related operations are in this directory.
If we decided to apply new db(such as S3), we can put the implementation in here and define interfaces to hide db specific operations.
### service/
- Define business logics of the application.
- There are not so much complex business logics existed so it's very simple. 
But this service layer will be needed if this application will be getting more large and has more feature specific logics in this application.
### testing/
- This directory is used when test operations are executed for file server. I want to separate directory from `/contents` as this directory is only used for the testing purpose.

## How to run locally.
### Pre-requests
- docker-compose is installed.

1. Execute `docker-compose build`
2. Execute `docker-compose up -d`

You can see db and app are running
```shell
[+] Running 2/2
 ⠿ Container storage-db    Healthy                                                                                                                                           0.5s
 ⠿ Container video-server  Running
```
3. Execute `docker-compose logs -f`

You can check logs from the containers.

Now you can see these logs belows. You can test with these endpoints with `http://localhost:8080`

```shell
video-server  | [GIN-debug] GET    /v1/files                 --> videoservice/handler.ServerHandler.GetFiles-fm (4 handlers)
video-server  | [GIN-debug] POST   /v1/files                 --> videoservice/handler.ServerHandler.PostFiles-fm (4 handlers)
video-server  | [GIN-debug] DELETE /v1/files/:fileid         --> videoservice/handler.ServerHandler.DeleteFilesFileId-fm (4 handlers)
video-server  | [GIN-debug] GET    /v1/files/:fileid         --> videoservice/handler.ServerHandler.GetFilesFileId-fm (4 handlers)
video-server  | [GIN-debug] GET    /v1/health                --> videoservice/handler.ServerHandler.GetHealth-fm (4 handlers)
```

## Others
### Why I use file system instead of DB?
Video contents is very large and it's not suitable for storing all byte data into the database. It's possible that a performance will be very slow and it's possible that the data will be broken when inserting and getting from it.
On the other hands, file system is just putting the file to the directory, and getting a file from the directory. So we don't need to worry about the content might be broken.

The demerit of file system is a performance of complex operations. It's possible to be slower than the Database.
In this application, to resolve this issue, I store the information of path, file name, id to the video contents in the RDB, and it makes it easier to search path of the video contents.
### Challenges

- Github actions.
In the Github actions flow, the test runner starts running as soon as the container created.
So in the first setting, the test runner starts running before db is setup.

To resolve this issue, I added the condition columns to the application definition.
```yml
    depends_on:
      db:
        condition: service_healthy
```

By adding this column, the application container starts running after the db container is setup and healthy condition, not just after the db container is created.

### Future considerations
- To expand this application, I think it's better to use S3 or other cloud file system to be capable of storing more contents with good performance.
- Also if we want to stream this contents, [MediaConvert](https://aws.amazon.com/jp/mediaconvert/) is preferable to endure more traffic and store caches with [AWS ClodFront](https://aws.amazon.com/jp/cloudfront/).
Sample architecture that I think is like [this](https://aws.amazon.com/jp/cdp/cdn/) if we choose AWS.