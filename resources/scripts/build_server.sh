cd /home/ec2-user/go/src/github.com/PrincetonOBO/OBOBackend/
export GOPATH=/home/ec2-user/go
cp -f resources/swagger/indexes/index.html resources/swagger
go get
go build 