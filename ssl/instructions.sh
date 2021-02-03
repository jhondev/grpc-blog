#!/bin/bash
INFO='\033[0;34m'
LINFO='\033[1;34m'
NORMAL='\033[1;37m'
SUCCESS='\033[1;32m'
QUESTION='\033[1;36m'
# Summary
# Private files: ca.key, server.key, server.pem, server.crt
# "Share" files: ca.crt (needed by the client), server.csr (needed by the CA)

echo -e "${LINFO}Certificate Generation${NORMAL}"
mkdir files
# Changes these CN's to match your hosts in your environment if needed
SERVER_CN=localhost
PASS=B781135297CDAC85112168541C211


echo "Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)"
openssl genrsa -passout pass:${PASS} -des3 -out files/ca.key 4096
openssl req -passin pass:${PASS} -new -x509 -days 365 -key files/ca.key -out files/ca.crt -subj "/CN=${SERVER_CN}" -addext "subjectAltName=DNS:${SERVER_CN}"
# openssl req -x509 -nodes -newkey rsa:2048 -days 3650 -sha256 -keyout test.key -out test.cert -reqexts SAN -extensions SAN -subj '/CN=test.example.com' -config <(cat /etc/pki/tls/openssl.cnf; printf "[SAN]\nsubjectAltName=DNS:test.example.com,DNS:test2.example.com")

echo "Step 2: Generate the Server Private Key (server.key)"
openssl genrsa -passout pass:${PASS} -des3 -out files/server.key 4096

echo "Step 3: Get a certificate signing request from the CA (server.csr)"
openssl req -passin pass:${PASS} -new -key files/server.key -subj "/CN=${SERVER_CN}" -addext "subjectAltName=DNS:${SERVER_CN}" -out files/server.csr

echo "Step 4: Sign the certificate with the CA we created (it's called self signing) - server.crt"
openssl x509 -req -passin pass:${PASS} -days 365 -in files/server.csr -CA files/ca.crt -CAkey files/ca.key -set_serial 01 -out files/server.crt

echo "Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC"
openssl pkcs8 -topk8 -nocrypt -passin pass:${PASS} -in files/server.key -out files/server.pem 

echo "Step 6: Convert the server certificate to .pfx format - usable by dotnet"
openssl pkcs12 -export -passin pass:${PASS} -in files/server.crt -inkey files/server.pem -passout pass:${PASS} -out files/server.pfx

echo -e "${SUCCESS}Done${NORMAL}"