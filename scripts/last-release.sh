#!/bin/bash

# Substitua pelo dono do repositório (usuário ou organização) e o nome do repositório
OWNER="carlos-enginner"
REPO="attom"

# Substitua pelo seu token de acesso pessoal do GitHub
TOKEN="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"

# A URL da API do GitHub para obter as releases
API_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"

# Fazendo a requisição com curl
curl -H "Authorization: token $TOKEN" -H "Accept: application/vnd.github.v3+json" $API_URL