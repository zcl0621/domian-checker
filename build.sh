docker build --platform linux/amd64 -t dns-backend:45 . && \
docker save -o backend-45.tar dns-backend:45 && \
scp backend-45.tar dns_vps:/tmp/backend-45.tar && \
rm -rf backend-45.tar && \
docker rmi dns-backend:45
