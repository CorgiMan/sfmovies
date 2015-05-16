FROM sfmovies/goimage

RUN apt-get update && \
    apt-get install -y golang && \
    mkdir /home/go && \
    export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies && \
    go install github.com/CorgiMan/sfmovies/gocode/dbupdate && \
    touch /home/apiservers

CMD /home/go/bin/dbupdate && \
    /bin/bash

    # run as cronjob
    # also do monitoring

# git clone https://github.com/CorgiMan/sfmovies.git
# docker build -t sfmovies/jobserver ./sfmovies/docker/managerserver
# docker run -d -p 12001:80 sfmovies/jobserver
