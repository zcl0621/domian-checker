docker build --platform linux/amd64 -t dns-backend:54 . && \
docker save -o backend-54.tar dns-backend:54 && \
scp backend-54.tar dns_vps:/tmp/backend-54.tar && \
rm -rf backend-54.tar && \
docker rmi dns-backend:54