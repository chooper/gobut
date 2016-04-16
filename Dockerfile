FROM golang:1.6
ADD . /go/src/github.com/chooper/gobut
RUN go get github.com/chooper/gobut
RUN go install github.com/chooper/gobut
ENTRYPOINT /go/bin/gobut
