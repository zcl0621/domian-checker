docker build --platform linux/amd64 -t dns-backend:44 . && \
docker save -o backend-44.tar dns-backend:44 && \
scp backend-44.tar dns_vps:/tmp/backend-44.tar && \
rm -rf backend-44.tar && \
docker rmi dns-backend:44
