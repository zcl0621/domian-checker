docker buildx build --platform linux/amd64 -t dns-backend:21 . && \
docker save -o backend-21.tar dns-backend:21 && \
scp backend-21.tar dns_vps:/tmp/backend-21.tar && \
rm -rf backend-21.tar && \
docker rmi dns-backend:21
