# Graph Report - .  (2026-08-28)

## Corpus Check
- Corpus is ~2,442 words - fits in a single context window. You may not need a graph.

## Summary
- 29 nodes · 26 edges · 9 communities (6 shown, 3 thin omitted)
- Extraction: 96% EXTRACTED · 4% INFERRED · 0% AMBIGUOUS · INFERRED: 1 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- Cloudflare Tunnel Management
- Reverse Proxy Server
- mDNS Hostname Broadcaster
- Project Package Root

## God Nodes (most connected - your core abstractions)
1. `Server` - 5 edges
2. `CloudflareTunnel` - 5 edges
3. `Broadcaster` - 4 edges
4. `Start()` - 3 edges
5. `NewServer()` - 3 edges
6. `Start()` - 3 edges
7. `ensureCloudflared()` - 2 edges
8. `routelocal` - 0 edges

## Surprising Connections (you probably didn't know these)
- `Start()` --calls--> `NewServer()`  [INFERRED]
  internal/hostname/mdns.go → internal/proxy/proxy.go
- `Broadcaster` --references--> `Server`  [EXTRACTED]
  internal/hostname/mdns.go → internal/proxy/proxy.go

## Import Cycles
- None detected.

## Communities (9 total, 3 thin omitted)

### Community 0 - "Cloudflare Tunnel Management"
Cohesion: 0.38
Nodes (5): CancelFunc, Cmd, ensureCloudflared(), Start(), CloudflareTunnel

## Knowledge Gaps
- **1 isolated node(s):** `routelocal`
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Server` connect `Reverse Proxy Server` to `mDNS Hostname Broadcaster`?**
  _High betweenness centrality (0.044) - this node is a cross-community bridge._
- **Why does `Broadcaster` connect `mDNS Hostname Broadcaster` to `Reverse Proxy Server`?**
  _High betweenness centrality (0.032) - this node is a cross-community bridge._
- **What connects `routelocal` to the rest of the system?**
  _1 weakly-connected nodes found - possible documentation gaps or missing edges._