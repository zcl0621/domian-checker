docker build --platform linux/amd64 -t dns-backend:34 . && \
docker save -o backend-34.tar dns-backend:34 && \
scp backend-34.tar dns_vps:/tmp/backend-34.tar && \
rm -rf backend-34.tar && \
docker rmi dns-backend:34
