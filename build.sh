docker build --platform linux/amd64 -t dns-backend:48 . && \
docker save -o backend-48.tar dns-backend:48 && \
scp backend-48.tar dns_vps:/tmp/backend-48.tar && \
rm -rf backend-48.tar && \
docker rmi dns-backend:48