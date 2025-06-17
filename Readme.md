# generate pem
openssl genrsa -out key/private.key 2048

openssl rsa -in key/private.key -pubout -out key/public.pem
