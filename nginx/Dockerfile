FROM nginx

RUN apt-get update && apt-get upgrade -y && apt-get install -y curl git gnupg

RUN cp nginx.conf /etc/nginx/

EXPOSE 80