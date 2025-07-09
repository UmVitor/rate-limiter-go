#!/bin/bash

# Cores para formatação da saída
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Iniciando Rate Limiter ===${NC}"

# Verificar se o Docker está instalado
if ! command -v docker &> /dev/null; then
    echo "Docker não está instalado. Por favor, instale o Docker e tente novamente."
    exit 1
fi

# Verificar se o Docker Compose está instalado
if ! command -v docker-compose &> /dev/null; then
    echo "Docker Compose não está instalado. Por favor, instale o Docker Compose e tente novamente."
    exit 1
fi

# Construir e iniciar os containers
echo -e "${YELLOW}Construindo e iniciando os containers...${NC}"
docker-compose up --build -d

# Verificar se os containers estão rodando
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Rate Limiter está rodando!${NC}"
    echo -e "Você pode acessar a API em: http://localhost:8080"
    echo -e "Para testar o rate limiter, execute: ./test.sh"
    echo -e "Para parar a aplicação, execute: docker-compose down"
else
    echo "Ocorreu um erro ao iniciar os containers."
    exit 1
fi
