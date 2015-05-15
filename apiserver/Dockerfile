FROM sfmovies/goimage

RUN export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies && \
    go install github.com/CorgiMan/sfmovies/apiserver && \
    cp /home/go/src/github.com/CorgiMan/sfmovies/apiserver/apidata.json /home/go/bin/

CMD export GOPATH=/home/go && \
    cd /home/go/bin && \
    ./apiserver --port 80

EXPOSE 80

# docker build -t sfmovies/apiserver //c/Users/Wieger/Desktop/sfmovies/apiserver
# docker run -d -p 12001:80 sfmovies/apiserver
# docker run -d -p 12001:80 -v //c/Users/Wieger/Desktop/go/src/github.com/CorgiMan/sfmovies:/home/go/src/github.com/CorgiMan/sfmovies sfmovies/apiserver


#RUN export GOPATH=/home/go
    #go install github.com/CorgiMan/sfmovies/apiserver 
    #echo $GOPATH && \
    #mkdir -p /home/go/src/github.com/CorgiMan && \
    #cp -r /c/Users/Wieger/Desktop/go/src/github.com/CorgiMan/sfmovies/ /home/go/src/github.com/CorgiMan/sfmovies && \
    #go install github.com/CorgiMan/sfmovies/apiserver 
    #cp /home/go/src/github.com/CorgiMan/sfmovies/apiserver/apidata.json /home/go/bin

