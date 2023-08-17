docker build --platform linux/amd64 -t dns-backend:83 . && \
docker save -o backend-83.tar dns-backend:83 && \
scp backend-83.tar dns_vps:/tmp/backend-83.tar && \
rm -rf backend-83.tar && \
docker rmi dns-backend:83