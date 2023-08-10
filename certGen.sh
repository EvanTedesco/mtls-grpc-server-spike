# Remove existing certs/keys
# Delete all files in ./certs directory
find ./certs -type f -exec rm {} +

# Generate self-signed root CA
cfssl selfsign -config cfssl.json --profile rootca "My Root CA" csr.json | cfssljson -bare root

# Create your server and client certificate
cfssl genkey csr.json | cfssljson -bare server
cfssl genkey csr.json | cfssljson -bare client

# Sign new certificates with self-signed root CA
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile server server.csr | cfssljson -bare server
cfssl sign -ca root.pem -ca-key root-key.pem -config cfssl.json -profile client client.csr | cfssljson -bare client

# Move generated certificates to certs folder
mv root.pem ./certs/ca.crt
mv root-key.pem ./certs/ca.key
mv client.pem ./certs/client.crt
mv client-key.pem ./certs/client.key
mv server.pem ./certs/server.crt
mv server-key.pem ./certs/server.key

# Delete unused root/server .csr files
rm *.csr