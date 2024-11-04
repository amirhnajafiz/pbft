#!/bin/bash

openssl genpkey -algorithm RSA -out ca.key -pkeyopt rsa_keygen_bits:2048
openssl req -new -x509 -key ca.key -out ca.crt -days 365

openssl genpkey -algorithm RSA -out nodes/S1/server.key -pkeyopt rsa_keygen_bits:2048
openssl req -new -key nodes/S1/server.key -out nodes/S1/server.csr
openssl x509 -req -in nodes/S1/server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out nodes/S1/server.crt -days 365

openssl genpkey -algorithm RSA -out clients/A/client.key -pkeyopt rsa_keygen_bits:2048
openssl req -new -key clients/A/client.key -out clients/A/client.csr
openssl x509 -req -in clients/A/client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out clients/A/client.crt -days 365
