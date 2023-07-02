docker build --platform linux/amd64 -t dns-backend:36 . && \
docker save -o backend-36.tar dns-backend:36 && \
scp backend-36.tar dns_vps:/tmp/backend-36.tar && \
rm -rf backend-36.tar && \
docker rmi dns-backend:36
