docker build --platform linux/amd64 -t dns-backend:40 . && \
docker save -o backend-40.tar dns-backend:40 && \
scp backend-40.tar dns_vps:/tmp/backend-40.tar && \
rm -rf backend-40.tar && \
docker rmi dns-backend:40
