# Desafio Rate Limiter

## Objetivo: 
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

## Descrição: 
O objetivo deste desafio é criar um rate limiter em Go que possa ser utilizado para controlar o tráfego de requisições para um serviço web. O rate limiter deve ser capaz de limitar o número de requisições com base em dois critérios:

Endereço IP: O rate limiter deve restringir o número de requisições recebidas de um único endereço IP dentro de um intervalo de tempo definido.
Token de Acesso: O rate limiter deve também poderá limitar as requisições baseadas em um token de acesso único, permitindo diferentes limites de tempo de expiração para diferentes tokens. O Token deve ser informado no header no seguinte formato:
API_KEY: <TOKEN>
As configurações de limite do token de acesso devem se sobrepor as do IP. Ex: Se o limite por IP é de 10 req/s e a de um determinado token é de 100 req/s, o rate limiter deve utilizar as informações do token.
## Requisitos:

- O rate limiter deve poder trabalhar como um middleware que é injetado ao servidor web
- O rate limiter deve permitir a configuração do número máximo de requisições permitidas por segundo.
- O rate limiter deve ter ter a opção de escolher o tempo de bloqueio do IP ou do Token caso a quantidade de requisições tenha sido excedida.
- As configurações de limite devem ser realizadas via variáveis de ambiente ou em um arquivo “.env” na pasta raiz.
- Deve ser possível configurar o rate limiter tanto para limitação por IP quanto por token de acesso.
- O sistema deve responder adequadamente quando o limite é excedido:
  - Código HTTP: 429
  - Mensagem: you have reached the maximum number of requests or actions allowed within a certain time frame
- Todas as informações de "limiter” devem ser armazenadas e consultadas de um banco de dados Redis. Você pode utilizar docker-compose para subir o Redis.
Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência.
- A lógica do limiter deve estar separada do middleware.

## Exemplos:
- Limitação por IP: Suponha que o rate limiter esteja configurado para permitir no máximo 5 requisições por segundo por IP. Se o IP 192.168.1.1 enviar 6 requisições em um segundo, a sexta requisição deve ser bloqueada.
- Limitação por Token: Se um token abc123 tiver um limite configurado de 10 requisições por segundo e enviar 11 requisições nesse intervalo, a décima primeira deve ser bloqueada.

Nos dois casos acima, as próximas requisições poderão ser realizadas somente quando o tempo total de expiração ocorrer. Ex: Se o tempo de expiração é de 5 minutos, determinado IP poderá realizar novas requisições somente após os 5 minutos.

## Como rodar a aplicação:
- Clone o projeto
- Vá até a raiz do projeto
- Edite as variáveis de ambiente dentro do arquivo `.env` para configurar as opções do limiter:
  - `LIMIT_REQUEST_IP` - Define o limite requests por IP dentro do tempo especificado.
  - `LIMIT_REQUEST_TOKEN` - Define o limite requests por Token dentro do tempo especificado.
  - `IP_BLOCKING_TIME` - Especifica o tempo para as requests de IP
  - `TOKEN_BLOCKING_TIME` - Especifica o tempo para as requests de Token
- Execute o seguinte comando para subir o servidor de teste com o rate-limiter configurado e o Redis:
```
docker compose up
```
- Execute uma chamada para o seguinte endereço `http://localhost:8080/` (ou conforme for configurado no arquivo `.env`) podendo ser adicionado ou não um header `API_KEY` com um token.
- Para facilitar há exemplos das chamadas dentro da pasta `api` na raiz do projeto.

## Como rodar os testes:
- Vá até a raiz do projeto
- Execute o seguinte comando para subir o serviço do Redis:
```
docker compose -f docker-compose-test.yaml up
```
- Execute o seguinte comando para executar os testes:
```
go test ./...
```
