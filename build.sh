docker build --platform linux/amd64 -t dns-backend:73 . && \
docker save -o backend-73.tar dns-backend:73 && \
scp backend-73.tar dns_vps:/tmp/backend-73.tar && \
rm -rf backend-73.tar && \
docker rmi dns-backend:73