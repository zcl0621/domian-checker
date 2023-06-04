docker buildx build --platform linux/amd64 -t dns-backend:18 . && \
docker save -o backend-18.tar dns-backend:18 && \
scp backend-18.tar dns_vps:/tmp/backend-18.tar && \
rm -rf backend-18.tar && \
docker rmi dns-backend:18 && \
ssh dns_vps
