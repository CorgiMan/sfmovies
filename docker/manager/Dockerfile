FROM ubuntu

RUN apt-get update && \
    apt-get install -y golang && \
    mkdir /home/go && \
    export GOPATH=/home/go && \
    apt-get install -y git && \
    go get github.com/CorgiMan/sfmovies/gocode && \
    cp /home/go/src/github.com/CorgiMan/sfmovies/gocode/manager/apiservers /home && \
    go install github.com/CorgiMan/sfmovies/gocode/dbupdate && \
    apt-get install -y ssh && \
    ssh-keygen -t rsa -f /root/.ssh/id_rsa -P ""


CMD echo install this pub key in all the nodes that run apiserver containers. && \
    cat /root/.shh/id_rsa.pub && \
    cat apiservers | while read ip; do ssh-copy-id docker@$ip; done && \
    /home/go/bin/dbupdate && \
    /bin/bash
    "* * * * * echo Hello World" >> /etc/cron.d/update_apiservers
    crontab -f /etc/cron.d/update_apiservers
    # run as cronjob
    # also do monitoring
