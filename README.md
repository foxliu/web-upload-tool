#### Web upload tool
This program is used to start a web server for uploading files to the current directory remotely

The first parameter of the program is the listening port of the web server, the default is 8030

How to Use
```shell
$go build main.go
$./main 9000
Use curl like this to upload file:
curl --location --request POST 'localhost:9000/' --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTIyNTY4MTYsImlzcyI6ImZpbGVVcGxvYWRlciJ9.Rf5zD5VMNtpYDPA6LfQJNRUY-uzK-7R4du1-jev1OxU' --form 'file=@"/your-file-path"'
```
Then can use the curl command to upload file