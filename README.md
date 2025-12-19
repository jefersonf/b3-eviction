# Laager Assessment 1 - Sistema de Votação (BBB Eviction)

Este projeto implementa um sistema de votação de alta performance, inspirado nos mecanismos de "paredão" do Big Brother Brasil (BBB). O objetivo é processar votos de forma escalável, permitindo também a consulta de estatísticas em tempo real.

### Tecnologias Utilizadas

- Linguagem: Go (Golang)
- Banco de Dados: PostgreSQL (instancia otimizada com timescaledb)
- Cache/Fila: Redis
- Proxy/Load Balancer: Nginx
- Frontend: JavaScript, HTML, CSS
- Infraestrutura: Docker & Docker Compose

### Arquitetura e Fluxos

O sistema é dividido em três fluxos principais para garantir performance e consistência:

#### 1. Votação (POST /vote)
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

#### 2. Analytics Horário (GET /analytics/hourly)

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

#### 3. Estatísticas do Paredão (GET /stats/{evictionId})

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