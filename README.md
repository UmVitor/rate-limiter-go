# Rate Limiter

Um serviço de rate limiter em Go que pode ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Funcionalidades

- Limitação de requisições por IP
- Limitação de requisições por token de acesso (API_KEY)
- Configuração via variáveis de ambiente ou arquivo .env
- Armazenamento em Redis
- Design com padrão Strategy para permitir diferentes implementações de armazenamento
- Middleware HTTP para fácil integração com servidores web

## Configuração

As configurações do rate limiter são feitas através de variáveis de ambiente ou arquivo `.env`:

```
# Rate Limiter Configuration
RATE_LIMITER_IP_LIMIT=10        # Número máximo de requisições por IP
RATE_LIMITER_IP_EXPIRATION=300  # Tempo de expiração do contador de IP (segundos)
RATE_LIMITER_TOKEN_LIMIT=100    # Número máximo de requisições por token
RATE_LIMITER_TOKEN_EXPIRATION=300 # Tempo de expiração do contador de token (segundos)
RATE_LIMITER_BLOCK_DURATION=300 # Duração do bloqueio quando o limite é excedido (segundos)

# Redis Configuration
REDIS_HOST=redis                # Host do Redis
REDIS_PORT=6379                 # Porta do Redis
REDIS_PASSWORD=                 # Senha do Redis (vazio se não houver)
REDIS_DB=0                      # Banco de dados Redis a ser usado

# Server Configuration
SERVER_PORT=8080                # Porta do servidor HTTP
```

## Como Executar

### Com Docker (Recomendado)

Utilize o script `start.sh` para iniciar a aplicação com Docker:

```bash
./start.sh
```

Este script verifica se o Docker e o Docker Compose estão instalados, constrói e inicia os containers, e fornece instruções para testar a aplicação.

Alternativamente, você pode iniciar manualmente com:

```bash
docker-compose up
```

Isso iniciará o servidor na porta 8080 e o Redis.

### Sem Docker

1. Certifique-se de ter um servidor Redis em execução
2. Configure o arquivo `.env` com os dados corretos do Redis
3. Execute:

```bash
go run main.go
```

## Testando

### Usando o Script de Teste Automatizado

Utilize o script `test.sh` para testar automaticamente o rate limiter:

```bash
./test.sh
```

Este script executa testes para:
- Limitação por IP (envia 12 requisições sequenciais)
- Limitação por token (envia 102 requisições com dois tokens diferentes)

O script mostra visualmente quando as requisições são aceitas (código 200) ou bloqueadas (código 429).

### Testes Manuais

Você também pode testar manualmente usando `curl` ou qualquer cliente HTTP:

#### Teste de limitação por IP

```bash
# Envie várias requisições para testar o limite por IP
for i in {1..15}; do curl -i http://localhost:8080/api/test; done
```

#### Teste de limitação por token

```bash
# Envie várias requisições com um token para testar o limite por token
for i in {1..15}; do curl -i -H "API_KEY: seu_token" http://localhost:8080/api/test; done
```

## Implementação

O rate limiter foi implementado seguindo os princípios de design orientado a interfaces e com separação clara de responsabilidades:

1. **Interface Storage**: Define uma interface para armazenamento que pode ser implementada por diferentes backends (Redis, memória, etc.)
2. **RateLimiter**: Contém a lógica de limitação de taxa, independente do armazenamento
3. **Middleware HTTP**: Integra o rate limiter com servidores HTTP

Quando o limite de requisições é excedido, o servidor responde com o código HTTP 429 e uma mensagem informativa.
