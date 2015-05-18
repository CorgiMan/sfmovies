FROM ubuntu

RUN apt-get update && \
    apt-get install -y golang && \
    mkdir /home/go && \
    export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies/gocode && \
    go install github.com/CorgiMan/sfmovies/gocode/apiserver

CMD /home/go/bin/apiserver --port 80

EXPOSE 80

