docker build --platform linux/amd64 -t dns-backend:77 . && \
docker save -o backend-77.tar dns-backend:77 && \
scp backend-77.tar dns_vps:/tmp/backend-77.tar && \
rm -rf backend-77.tar && \
docker rmi dns-backend:77