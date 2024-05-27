# Dicom Server

## Overview
This application is desgined for 3 specific functionalities:
1. Users can upload a DICOM file
2. Users can retrieve the DICOM file header attributes through [tags](https://www.dicomlibrary.com/dicom/dicom-tags/)
3. Users can retrieve the file as an image as PNG/JPG

This application implements the repository pattern it is a good way to separate logic between the controller, service and repository.
Since there isn't a requirement to implement a database, the repository is not implemented.

## Issues
Using the [Go DICOM Library](https://github.com/suyashkumar/dicom) causes the conversion to PNG/JPG to become pitch black in the response.
It appears that there's missing auto-scaling in the `frames` and `nativeFrames` `GetImage` function.
The solution to have this work requires going into a photo editor and manually change the brightness,contranst and exposure.
### Links
- https://github.com/suyashkumar/dicom/issues/164
- https://github.com/suyashkumar/dicom/issues/301
   
## API Endpoints
| HTTP Methods | Endpoint | Parameters | Tags | Description |
| ---- | ---- | ---- | ---- | ---- |
| GET | /ping | N/A | N/A | Pings the server, should return pong in the response |
| POST | /file/upload | N/A | N/A | Upload a DICOM file and a file ID is returned in the response |
| GET | /file/:id | id - file ID | N/A | Retrieves the header attribues of an uploaded DICOM file |
| GET | /file/:id/image | id - file ID | fileType: png/jpg | Converts an existing DICOM file to either PNG/JPG and returns it to the user | 

## Getting Started
This project will start on port `8080`
### Via Docker
Pull down the repository start the project
```bash
git clone https://github.com/w3ye/dicomserver.git
cd dicomserver
docker compose up --build
```
To stop the container
```bash
docker compose down --remove-orphans
```
### Via Local
```bash
git clone https://github.com/w3ye/dicomserver.git
cd dicomserver
go run main.go
```
## Testing the endpoints
### Via Postman
Import the collection `dicom_server.postman_collection.json` at the root of the directory
![image](https://github.com/w3ye/dicomserver/assets/33244107/be168003-88aa-470f-8314-8fb5ede0879c)

The `upload` endpoint should be called first. It will set the `id` as a global variable so other endpoints could use
To call the upload endpoint you must add the file first. Within postman `Body` -> `form-data` -> `Select Files`
![image](https://github.com/w3ye/dicomserver/assets/33244107/d44c27ae-180d-45a3-8eae-029addf1e526)

If the `GET /file/:id/image` endpoint ever returns a black box. Save the response as a png file. Open the file in a photo editor and adjust the brightness, contrast and exposure
![image](https://github.com/w3ye/dicomserver/assets/33244107/2ab81ea2-4605-45e0-a475-6e059a79eaa1)
![image](https://github.com/w3ye/dicomserver/assets/33244107/fe8300bb-51ff-4b6f-b46f-cda6989ac388)

Exmampe response after manual editing
![image](https://github.com/w3ye/dicomserver/assets/33244107/fb24ec90-f535-48a5-8270-7caf9a69f85f)

### Via CURL
`GET /ping`
```curl
curl http://localhost:8080/ping
```
Example response 
```
{"message": "pong"}
```

`POST /upload`
```curl
curl -i -X POST \
-H 'content-type: multipart/form-data' \
-F file=@<FILE_PATH> \
localhost:8080/file/upload
```
Example response
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 27 May 2024 01:14:23 GMT
Content-Length: 45

{"id":"6b3a1037-db46-4c41-bf0f-2726dd68743c"}
```

`GET /file/:id`
```
curl "localhost:8080/file/<FILE_ID>?tag=(0002,0000)"
```
Example response
```
{"tag":"(0002,0000)","name":"FileMetaInformationGroupLength","vr":"UL","value":"[186]"}
```

`GET /file/:id/image`
```
curl -H "content-type: image/png" \
"localhost:8080/file/<FILE_ID>/image?fileType=png" \
-o downloaded_image.png
```
Example response
```
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  115k    0  115k    0     0  1426k      0 --:--:-- --:--:-- --:--:-- 1426k
```

## Dependencies
go v1.22.3

air v1.52.0

docker v24.0.7

docker-compose v2.23.3

github.com/google/uuid v1.6.0
github.com/suyashkumar/dicom v1.0.7
