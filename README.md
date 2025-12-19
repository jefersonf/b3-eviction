# Laager Assessment 1 - Sistema de Votação (BBB Eviction)

Este projeto implementa um sistema de votação de alta performance, inspirado nos mecanismos de "paredão" do Big Brother Brasil (BBB). O objetivo é processar votos de forma escalável, permitindo também a consulta de estatísticas em tempo real.

## Tecnologias Utilizadas

- Linguagem: Go (Golang)
- Banco de Dados: PostgreSQL (instancia otimizada com timescaledb)
- Cache/Fila: Redis
- Proxy/Load Balancer: Nginx
- Frontend: JavaScript, HTML, CSS
- Infraestrutura: Docker & Docker Compose

## Arquitetura e Fluxos

O sistema é dividido em três fluxos principais para garantir performance e consistência:

### 1. Votação (POST /vote)
Processa o voto do usuário e o enfileira para contabilização assíncrona.

<details open>
      <summary> Mostrar fluxo de voto </summary>

```mermaid
---
config:
      theme: redux
---
flowchart TD
        V(["POST /vote"])
        V --> B["load-balancer"]
        B --> LB{"round-robin"}
        LB --> C["voting-api-instance-1"]
        LB --> D["voting-api-instance-2"]
        C --> Q["voting-queue"]
        D --> Q
        Q <--> A["vote-aggregator"]
        A --> P["postgres"] 
```

</details>

### 2. Analytics Horário (GET /analytics/hourly)

Recupera dados consolidados por hora diretamente do banco de dados persistente.

<details>
      <summary> Mostrar fluxo de consulta das estatísticas por hora. </summary>

```mermaid
---
config:
      theme: redux
---
flowchart TD
        V(["GET /analytics/hourly"])
        V <--> B["load-balancer"]
        B <--> LB{"round-robin"}
        LB <--> C["voting-api-instance-1"]
        LB <--> D["voting-api-instance-2"]
        C & D <--> P["postgres"] 
```

</details>

### 3. Estatísticas do Paredão (GET /stats/{evictionId})

Consulta o estado atual da votação, geralmente acessando dados em cache (Redis) para rapidez.


<details>
      <summary> Mostrar fluxo de consulta das estatísticas da votação </summary>

```mermaid
---
config:
      theme: redux
---
flowchart TD
        V(["GET /{evictionId}"])
        V <--> B["load-balancer"]
        B <--> LB{"round-robin"}
        LB <--> C["voting-api-instance-1"]
        LB <--> D["voting-api-instance-2"]
        C & D <--> P["voting-queue"] 
```

</details>


## Como Executar

O projeto utiliza `Docker` para orquestração dos serviços e um `Makefile` (opcional) para automação de tarefas de build/execução.

### Pré-requisitos

- Docker e Docker Compose instalados/configurados.

#### Passos

1. Clone o repositório:

```bash
git clone https://github.com/jefersonf/b3-eviction.git
cd b3-eviction
```

2. Suba os serviços: Utilize o Docker Compose para iniciar a aplicação, banco de dados, redis e nginx.

```bash
docker-compose up -d
```
Alternativamente, e o autor RECOMENDA, podes subir os serviços via Makefile pela praticidade, pois o comando a seguir constroi e executa a aplicação, isto é, equivale `make build` seguido de `make run` (checar arquivo Makefile).

```bash
make
```

1. Acesse a aplicação: O serviço deve estar acessível através das portas configuradas para acesso externo ao Docker Compose. São três os serviços ao todo:
   - Frontend (`localhost:3000`)
   - Nginx/LoadBalancer/Entrypoint da API (`localhost:8080`) e o
   - Loscust (`localhost:8089`), serviçõ extra para fins de testes de alto volume de requisiões contra a API da votação.


## Documentação da API

|Verbo HTTP|Recurso|Descrição Simples|
|-|-|-|
|`POST`|`/vote`|Solicita um registro de um voto em um dos indicados.|
|`GET`|`/stats/`|Retorna informação sobre a saúde da API de votação.|
|`GET`|`/stats/{evictionId}`|Retorna as estatísticas de um paredão específico.|
|`GET`|`/analytics/hourly`|Retorna o sumário de votos totais e por cada indicado a cada hora nas últimas 24 horas.|
|`GET`|`/analytics/minutely`|Retorna o sumário de votos totais e por cada indicado a cada minuto nas últimas 24 horas.|

## Documentação da UI (Área de votação + Dashboard)

A interface de usuário tanto para votação quanto para quem assume o papel de administrador do dashboard, compreende duas visões, respectivamente.

### 1. Área de votação

![Tela de Votação](/docs/artifacts/fluxo-area-de-votacao.png)

### 2. Dashboard para acompanhamento das votação ao vivo.

*TODO adcionar descricao e screeshots aqui*

## Modelo de dados

### Redis (in-memory/cached)

- `EvictionStats` - Sumário da votação em memória.

### Postgres (persisted)

- `votes_minutely` - Tabela de votos agregados e consolidados por minuto.
- `votes_hourly` - Tabela de Visão de votos agregados e consolidados por hora.

## Demonstações em vídeo

- Clone do projeto, construção, demonstração de voto e execução do teste de carga: [link-video-1]()
- Verificando dos votos armazendos: [link-video-2]()