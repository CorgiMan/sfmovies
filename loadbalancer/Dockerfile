FROM ubuntu

RUN \
  apt-get update && \
  apt-get install -y software-properties-common python-software-properties && \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get install -y nginx && \
  apt-get install -y nano && \
  rm -rf /var/lib/apt/lists/* && \
  echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  chown -R www-data:www-data /var/lib/nginx

VOLUME ["/etc/nginx/sites-enabled", "/etc/nginx/certs", "/etc/nginx/conf.d", "/var/log/nginx", "/var/www/html"]
CMD nginx -c /home/nginx.conf && \
    /bin/bash
EXPOSE 80

# docker run -dit -p 80:80 -v //c/Users/Wieger/Desktop/sfmovies/loadbalancer:/home/ sfmovies/nginx
# download mynginx.conf from github with wget