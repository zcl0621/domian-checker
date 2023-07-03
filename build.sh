docker build --platform linux/amd64 -t dns-backend:41 . && \
docker save -o backend-41.tar dns-backend:41 && \
scp backend-41.tar dns_vps:/tmp/backend-41.tar && \
rm -rf backend-41.tar && \
docker rmi dns-backend:41
