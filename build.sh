docker build --platform linux/amd64 -t dns-backend:64 . && \
docker save -o backend-64.tar dns-backend:64 && \
scp backend-64.tar dns_vps:/tmp/backend-64.tar && \
rm -rf backend-64.tar && \
docker rmi dns-backend:64