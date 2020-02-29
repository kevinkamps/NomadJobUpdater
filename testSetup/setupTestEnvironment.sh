#!/usr/bin/env bash

NOMAD_VERSION=0.10.4
CONSUL_VERSION=1.7.1

if [ ! -f "nomad" ]; then
    wget "https://releases.hashicorp.com/nomad/$NOMAD_VERSION/nomad_${NOMAD_VERSION}_linux_amd64.zip"
    unzip nomad_${NOMAD_VERSION}_linux_amd64.zip
    rm -f nomad_${NOMAD_VERSION}_linux_amd64.zip
fi

if [ ! -f "consul" ]; then
    wget "https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_amd64.zip"
    unzip consul_${CONSUL_VERSION}_linux_amd64.zip
    rm -f consul_${CONSUL_VERSION}_linux_amd64.zip
fi

if [ ! -f "certificates/cfssl" ]; then
    cd certificates
    wget "https://pkg.cfssl.org/R1.2/cfssl_linux-amd64"
    wget "https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64"
    chmod +x cfssl_linux-amd64
    chmod +x cfssljson_linux-amd64
    mv cfssl_linux-amd64 cfssl
    mv cfssljson_linux-amd64 cfssljson

    ./cfssl print-defaults csr | ./cfssl gencert -initca - | ./cfssljson -bare nomad-ca

    echo '{}' | ./cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -config=cfssl.json \
        -hostname="server.global.nomad,192.168.168.112,localhost,127.0.0.1" - | ./cfssljson -bare server

    # Generate a certificate for the Nomad client
    echo '{}' | ./cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -config=cfssl.json \
    -hostname="client.global.nomad,192.168.168.112,localhost,127.0.0.1" - | ./cfssljson -bare client

    # Generate a certificate for the CLI
    echo '{}' | ./cfssl gencert -ca=nomad-ca.pem -ca-key=nomad-ca-key.pem -profile=client \
    - | ./cfssljson -bare cli
    cd ..
fi


echo "Test setup is installed!"

echo "You can now start all dependencies by using the following commands:"
echo "./startConsul.sh"
echo "./startNomad.sh"
echo "docker-compose up"

