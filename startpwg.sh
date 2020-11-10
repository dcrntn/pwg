#!/bin/bash  

FILE="/etc/wsgateway/pwg"

if [[ -f "$FILE" ]]; then
    /etc/wsgateway/./pwg
else
    wget https://github.com/dcrntn/pwg/blob/main/pwg?raw=true
    mkdir /etc/wsgateway
    cp pwg?raw=true /etc/wsgateway/pwg
    rm pwg?raw=true
    chmod +x /etc/wsgateway/pwg
    /etc/wsgateway/./pwg
fi
