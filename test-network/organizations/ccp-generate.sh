#!/bin/bash

function one_line_pem {
    echo "`awk 'NF {sub(/\\n/, ""); printf "%s\\\\\\\n",$0;}' $1`"
}

function json_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        organizations/ccp-template.json
}

function yaml_ccp {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    local OP=$(one_line_pem $8)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s/\${ORDERERPORT}/$7/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${CRYPTO}#$6#" \
        -e "s#\${ORDERERPEM}#$OP#" \
        organizations/ccp-template.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

function yaml_ccp2 {
    local PP=$(one_line_pem $4)
    local CP=$(one_line_pem $5)
    local OP=$(one_line_pem $8)
    sed -e "s/\${ORG}/$1/" \
        -e "s/\${P0PORT}/$2/" \
        -e "s/\${CAPORT}/$3/" \
        -e "s/\${ORDERERPORT}/$7/" \
        -e "s#\${PEERPEM}#$PP#" \
        -e "s#\${CAPEM}#$CP#" \
        -e "s#\${CRYPTO}#$6#" \
        -e "s#\${ORDERERPEM}#$OP#" \
        organizations/ccp-template-2.yaml | sed -e $'s/\\\\n/\\\n          /g'
}

ORDERERPORT=7050
ORDERERPEM=organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

ORG=1
P0PORT=7051
CAPORT=7054
PEERPEM=organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
CAPEM=organizations/peerOrganizations/org1.example.com/ca/ca.org1.example.com-cert.pem
CRYPTO=organizations/peerOrganizations/org1.example.com

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org1.example.com/connection-org1.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $CRYPTO $ORDERERPORT $ORDERERPEM)" > organizations/peerOrganizations/org1.example.com/connection-org1.yaml
echo "$(yaml_ccp2 $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $CRYPTO $ORDERERPORT $ORDERERPEM)" > organizations/peerOrganizations/org1.example.com/dconnection-org1.yaml


ORG=2
P0PORT=9051
CAPORT=8054
PEERPEM=organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
CAPEM=organizations/peerOrganizations/org2.example.com/ca/ca.org2.example.com-cert.pem
CRYPTO=organizations/peerOrganizations/org2.example.com

echo "$(json_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM)" > organizations/peerOrganizations/org2.example.com/connection-org2.json
echo "$(yaml_ccp $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $CRYPTO $ORDERERPORT $ORDERERPEM)" > organizations/peerOrganizations/org2.example.com/connection-org2.yaml
echo "$(yaml_ccp2 $ORG $P0PORT $CAPORT $PEERPEM $CAPEM $CRYPTO $ORDERERPORT $ORDERERPEM)" > organizations/peerOrganizations/org2.example.com/dconnection-org2.yaml
