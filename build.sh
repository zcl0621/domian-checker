docker build --platform linux/amd64 -t dns-backend:87 . && \
docker save -o backend-87.tar dns-backend:87 && \
scp backend-87.tar dns_vps:/tmp/backend-87.tar && \
rm -rf backend-87.tar && \
docker rmi dns-backend:87