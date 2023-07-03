docker build --platform linux/amd64 -t dns-backend:38 . && \
docker save -o backend-38.tar dns-backend:38 && \
scp backend-38.tar dns_vps:/tmp/backend-38.tar && \
rm -rf backend-38.tar && \
docker rmi dns-backend:38
