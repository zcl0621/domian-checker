docker build --platform linux/amd64 -t dns-backend:67 . && \
docker save -o backend-67.tar dns-backend:67 && \
scp backend-67.tar dns_vps:/tmp/backend-67.tar && \
rm -rf backend-67.tar && \
docker rmi dns-backend:67