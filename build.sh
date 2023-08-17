docker build --platform linux/amd64 -t dns-backend:85 . && \
docker save -o backend-85.tar dns-backend:85 && \
scp backend-85.tar dns_vps:/tmp/backend-85.tar && \
rm -rf backend-85.tar && \
docker rmi dns-backend:85