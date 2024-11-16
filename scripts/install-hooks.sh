#!/bin/bash

# Verifica se o diretório .git/hooks existe
if [ ! -d ".git/hooks" ]; then
  echo "Diretório .git/hooks não encontrado. Este repositório é um repositório Git?"
  exit 1
fi

# Copia todos os hooks do diretório .githooks para .git/hooks
cp githooks/* .git/hooks/

# Torna os hooks executáveis
chmod +x .git/hooks/*
echo "Hooks instalados com sucesso!"