docker build --platform linux/amd64 -t dns-backend:39 . && \
docker save -o backend-39.tar dns-backend:39 && \
scp backend-39.tar dns_vps:/tmp/backend-39.tar && \
rm -rf backend-39.tar && \
docker rmi dns-backend:39
