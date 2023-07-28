docker build --platform linux/amd64 -t dns-backend:43 . && \
docker save -o backend-43.tar dns-backend:43 && \
scp backend-43.tar dns_vps:/tmp/backend-43.tar && \
rm -rf backend-43.tar && \
docker rmi dns-backend:43
