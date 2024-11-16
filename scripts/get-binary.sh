#!/bin/bash

# Variáveis
OWNER="carlos-enginner"
REPO="attom"
TOKEN="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"  # Substitua pelo seu token

wget --auth-no-challenge --header='Accept:application/octet-stream' https://$TOKEN:@api.github.com/repos/$OWNER/$REPO/releases/assets/206845800 -O attom_1.0.0_linux_amd64

