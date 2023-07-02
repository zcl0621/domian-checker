docker build --platform linux/amd64 -t dns-backend:37 . && \
docker save -o backend-37.tar dns-backend:37 && \
scp backend-37.tar dns_vps:/tmp/backend-37.tar && \
rm -rf backend-37.tar && \
docker rmi dns-backend:37
