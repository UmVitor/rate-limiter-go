#!/bin/bash

# Cores para formatação da saída
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Testando Rate Limiter ===${NC}"

# URL base para os testes
BASE_URL="http://localhost:8080"

# Função para testar limitação por IP
test_ip_limit() {
    echo -e "\n${YELLOW}Testando limitação por IP (limite padrão: 10 req/s)${NC}"
    echo "Enviando 12 requisições sequenciais..."

    for i in {1..12}; do
        response=$(curl -s -o /dev/null -w "%{http_code}" $BASE_URL/api/test)

        if [ $response -eq 200 ]; then
            echo -e "${GREEN}Requisição $i: OK (200)${NC}"
        else
            echo -e "${RED}Requisição $i: Bloqueada ($response)${NC}"
        fi

        # Pequena pausa para não sobrecarregar o terminal
        sleep 0.1
    done
}

# Função para testar limitação por token
test_token_limit() {
    local token=$1
    echo -e "\n${YELLOW}Testando limitação por token: $token (limite padrão: 100 req/s)${NC}"
    echo "Enviando 102 requisições sequenciais com o token..."

    for i in {1..102}; do
        response=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: $token" $BASE_URL/api/test)

        # Mostrar apenas a cada 10 requisições para não poluir a saída
        if [ $((i % 10)) -eq 0 ] || [ $i -eq 1 ] || [ $i -gt 98 ]; then
            if [ $response -eq 200 ]; then
                echo -e "${GREEN}Requisição $i: OK (200)${NC}"
            else
                echo -e "${RED}Requisição $i: Bloqueada ($response)${NC}"
            fi
        fi

        # Pequena pausa para não sobrecarregar o terminal
        sleep 0.01
    done
}

# Verificar se o servidor está rodando
echo "Verificando se o servidor está rodando..."
if curl -s -o /dev/null $BASE_URL; then
    echo -e "${GREEN}Servidor está rodando!${NC}"
else
    echo -e "${RED}Servidor não está respondendo. Certifique-se de que o servidor esteja rodando na porta 8080.${NC}"
    exit 1
fi

# Executar os testes
test_ip_limit
test_token_limit "test-token-1"

echo -e "\n${YELLOW}Aguardando 5 segundos antes de testar outro token...${NC}"
sleep 5

# Testar com outro token para mostrar que os limites são independentes
test_token_limit "test-token-2"

echo -e "\n${GREEN}Testes concluídos!${NC}"
