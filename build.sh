docker build --platform linux/amd64 -t dns-backend:35 . && \
docker save -o backend-35.tar dns-backend:35 && \
scp backend-35.tar dns_vps:/tmp/backend-35.tar && \
rm -rf backend-35.tar && \
docker rmi dns-backend:35
