docker build --platform linux/amd64 -t dns-backend:70 . && \
docker save -o backend-70.tar dns-backend:70 && \
scp backend-70.tar dns_vps:/tmp/backend-70.tar && \
rm -rf backend-70.tar && \
docker rmi dns-backend:70