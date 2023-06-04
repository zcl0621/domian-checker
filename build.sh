docker buildx build --platform linux/amd64 -t dns-backend:17 . && \
docker save -o backend-17.tar dns-backend:17 && \
scp backend-17.tar dns_vps:/tmp/backend-17.tar && \
rm -rf backend-17.tar && \
docker rmi dns-backend:17 && \
ssh dns_vps
