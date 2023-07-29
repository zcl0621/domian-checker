docker build --platform linux/amd64 -t dns-backend:47 . && \
docker save -o backend-47.tar dns-backend:47 && \
scp backend-47.tar dns_vps:/tmp/backend-47.tar && \
rm -rf backend-47.tar && \
docker rmi dns-backend:47