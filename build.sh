docker build --platform linux/amd64 -t dns-backend:82 . && \
docker save -o backend-82.tar dns-backend:82 && \
scp backend-82.tar dns_vps:/tmp/backend-82.tar && \
rm -rf backend-82.tar && \
docker rmi dns-backend:82