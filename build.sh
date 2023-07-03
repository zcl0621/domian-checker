docker build --platform linux/amd64 -t dns-backend:42 . && \
docker save -o backend-42.tar dns-backend:42 && \
scp backend-42.tar dns_vps:/tmp/backend-42.tar && \
rm -rf backend-42.tar && \
docker rmi dns-backend:42
