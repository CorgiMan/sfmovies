FROM ubuntu

RUN apt-get update && \
    apt-get install -y golang && \
    mkdir /home/go && \
    export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies/gocode && \
    cp -r /home/go/src/github.com/CorgiMan/sfmovies/docker/manager /home && \
    go install github.com/CorgiMan/sfmovies/gocode/dbupdate && \
    apt-get install -y ssh && \
    ssh-keygen -t rsa -f /root/.ssh/id_rsa -P ""

CMD echo "install this pub key in all the apiserver nodes (the server that hosts the containers)." && \
    cat /root/.ssh/id_rsa.pub && \
    echo "* * * * * sh /home/manager/monitor.sh >> /home/log" > /home/cron_jobs && \
    echo "0 4 * * * sh /home/manager/updatedb.sh >> /home/log" >> /home/cron_jobs && \
    crontab /home/cron_jobs && \
    cron && \
    /bin/bash

# for boot2docker type: ssh-copy-id docker@192.168.59.103
# password: tcuser