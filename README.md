Pre-requisites to run this sample
- Go should be installed and GOPATH should be set
- get the dependencies using go get ./... from the directory where you keep jwtdemo.go
- setup the mongodb database, run mongod service
 - lanuch mongo shell
 - use books
 - db.book.insert({"PageNum":2,"content":"This is a sample page"})
- go run jwtdemo

This runs the webservice on port 3080

For testing
- echo -n 'Synerzip1:Password1' | base64
  This gives U3luZXJ6aXAxOlBhc3N3b3JkMQ==
- curl -H "Authentication : U3luZXJ6aXAxOlBhc3N3b3JkMQ==" http://localhost:3080/authenticate
  This gives a JWT token
- curl -H "Authentication : JWT token " http://localhost:3080/book/page/2
