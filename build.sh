docker build --platform linux/amd64 -t dns-backend:55 . && \
docker save -o backend-55.tar dns-backend:55 && \
scp backend-55.tar dns_vps:/tmp/backend-55.tar && \
rm -rf backend-55.tar && \
docker rmi dns-backend:55