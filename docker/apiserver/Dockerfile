FROM sfmovies/goimage

RUN apt-get update && \
    apt-get install -y golang && \
    mkdir /home/go && \
    export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies && \
    go install github.com/CorgiMan/sfmovies/gocode/apiserver

CMD /home/go/binapiserver --port 80

EXPOSE 80

# git clone https://github.com/CorgiMan/sfmovies.git
# docker build -t sfmovies/apiserver ./sfmovies/docker/apiserver
# docker run -d -p 12001:80 sfmovies/apiserver
