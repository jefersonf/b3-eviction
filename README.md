# Laager Assessment 1 - Voting System (BBB Eviction)

## Project dependencies

- Postgres
- Go
- Nginx
- Redis

## Use case workflows

### Voting

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

### Fetching houtly analytics data

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

### Fetching eviction stats

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