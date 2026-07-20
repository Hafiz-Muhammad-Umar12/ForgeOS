# Multi-Agent Orchestration Systems: Technical Architecture Deep Dive (2025-2026)

**Research Date:** July 2026  
**Sources:** Official GitHub repositories, documentation, source code analysis, and technical papers

---

## Table of Contents

1. [Major Orchestration Frameworks](#1-major-orchestration-frameworks)
   - 1.1 AutoGen / Microsoft Agent Framework
   - 1.2 LangGraph
   - 1.3 CrewAI
   - 1.4 Anthropic Claude Multi-Agent Patterns
2. [Software Engineering Agent Systems](#2-software-engineering-agent-systems)
   - 2.1 Devin (Cognition Labs)
   - 2.2 SWE-agent
   - 2.3 OpenHands (formerly OpenDevin)
   - 2.4 MetaGPT
   - 2.5 ChatDev
3. [Message Bus and Event-Driven Patterns](#3-message-bus-and-event-driven-patterns)
4. [Agent Memory and Context Management](#4-agent-memory-and-context-management)
5. [Workspace Isolation and Sandboxing](#5-workspace-isolation-and-sandboxing)
6. [Architectural Comparison Matrix](#6-architectural-comparison-matrix)
7. [Key Design Patterns and Trends](#7-key-design-patterns-and-trends)

---

## 1. Major Orchestration Frameworks

### 1.1 AutoGen / Microsoft Agent Framework

**Status:** AutoGen is now in maintenance mode; superseded by Microsoft Agent Framework (MAF)

**Architecture Philosophy:** Layered, extensible design built on the Actor model for message-passing between agents.

**Core Architecture (3 layers):**

1. **Core API** (autogen-core):
   - Event-driven, distributed agent system using the Actor model
   - Message passing between agents via topics and subscriptions
   - Local and distributed runtime support
   - Cross-language support (.NET and Python)
   - Agents publish/subscribe to typed messages through topic-based routing

2. **AgentChat API** (autogen-agentchat):
   - Higher-level API built on Core API
   - Supports two-agent chat, group chat patterns
   - Pre-built agent types: AssistantAgent, UserProxyAgent
   - Team-based orchestration (RoundRobinGroupChat, SelectorGroupChat, Swarm)
   - AgentTool pattern for nested agent-as-tool delegation

3. **Extensions API** (autogen-ext):
   - Pluggable LLM clients (OpenAI, Azure, etc.)
   - Code execution capabilities
   - MCP (Model Context Protocol) integration

**Microsoft Agent Framework (MAF) - The Successor:**
- Production-grade, enterprise-ready successor
- Graph-based orchestration patterns: sequential, concurrent, handoff, group collaboration
- Built-in checkpointing, streaming, human-in-the-loop, time-travel
- Middleware system for request/response processing
- OpenTelemetry integration for observability
- A2A (Agent-to-Agent) protocol support for cross-runtime interoperability
- Supports both Python and .NET with consistent APIs

**Key Design Decisions:**
- Message-driven communication (not direct RPC)
- Topic-based pub/sub for decoupled agent communication
- Typed message schemas for safety
- Runtime abstraction allows local → distributed scaling

---

### 1.2 LangGraph

**Architecture Philosophy:** Graph-based (state machine) orchestration inspired by Google Pregel and Apache Beam.

**Core Concepts:**

1. **StateGraph:**
   - Nodes represent processing steps (functions or classes)
   - Edges define transitions between nodes
   - Shared state object flows through the graph
   - State schema defined with optional reducer functions for aggregating updates from multiple nodes
   - Signature pattern: `State -> Partial<State>`

2. **Channels:**
   - Communication mechanism between nodes
   - Channel types: LastValue, BinaryOperatorAggregate, DeltaChannel, EphemeralValue, NamedBarrierValue
   - Reducer functions handle concurrent updates to the same state key

3. **Checkpointer:**
   - Durable execution - agents persist through failures
   - Can resume from exactly where they left off
   - Supports time-travel debugging
   - Thread-level state management for conversational memory

4. **Compilation and Execution:**
   - StateGraph is a builder; must call `.compile()` to create executable graph
   - Supports `invoke()`, `stream()`, `astream()`, `ainvoke()`
   - Retry policies, cache policies, timeout policies configurable per node
   - Error handler nodes for fault tolerance

5. **Multi-Agent Patterns:**
   - Subgraphs for hierarchical agent composition
   - Send() API for dynamic fan-out to multiple nodes
   - Command objects for explicit control flow
   - Human-in-the-loop via interrupts at any point

**Key Differentiators:**
- Lowest-level of control among major frameworks
- Explicit state management (no magic)
- Cycle support (unlike pure DAG systems)
- Built-in persistence and time-travel
- Production deployment via LangSmith

---

### 1.3 CrewAI

**Architecture Philosophy:** Role-based agent design simulating a "crew" of specialists with autonomous collaboration.

**Two Core Abstractions:**

1. **Crews:**
   - Teams of AI agents with role-based specialization
   - Each agent has: role, goal, backstory, tools, LLM configuration
   - Dynamic task delegation between agents
   - Three process types:
     - Sequential: tasks execute in order
     - Parallel: independent tasks run simultaneously
     - Hierarchical: manager agent delegates to worker agents
   - Built-in memory system (short-term, long-term, entity memory)
   - Context sharing between tasks

2. **Flows:**
   - Event-driven, production-ready workflows
   - Fine-grained control over execution paths
   - State management between steps (typed state objects)
   - Conditional branching for business logic
   - Can embed Crews as steps within Flows
   - Decorator-based step definitions (`@start`, `@listen`, `@router`)

**Architecture Details:**
- Agent configuration via YAML or Python classes
- Tool integration pattern (similar to LangChain tools)
- Guardrails for output validation
- Human review checkpoints
- Structured output support (Pydantic models, JSON schemas)

**Key Design Decisions:**
- High-level abstractions for rapid development
- Separation of concerns: Crews for autonomous AI, Flows for deterministic logic
- Memory persistence across executions
- Enterprise control plane (CrewAI AMP) for observability and governance

---

### 1.4 Anthropic Claude Multi-Agent Patterns

**Architecture Philosophy:** Tool-use and API-first approach with Model Context Protocol (MCP).

**Multi-Agent Patterns:**

1. **Orchestrator-Worker Pattern:**
   - Single orchestrator agent decomposes tasks
   - Dispatches subtasks to specialized worker agents
   - Workers return results to orchestrator for synthesis

2. **Tool-Use Agents:**
   - Agents equipped with tools (code execution, web search, file operations)
   - Extended thinking for complex reasoning
   - Tool use loop: think → act → observe → repeat

3. **Model Context Protocol (MCP):**
   - Standardized protocol for agent-tool communication
   - JSON-RPC based message format
   - Supports tool discovery, invocation, and resource access
   - Transport-agnostic (stdio, SSE, HTTP)

4. **Agent-Client Protocol (ACP):**
   - For agent-to-agent communication
   - OpenHands and other frameworks adopting this standard

**Key Differentiators:**
- Focus on safety and alignment
- Constitutional AI principles
- Extended thinking for complex multi-step reasoning
- MCP as an open standard for tool integration

---

## 2. Software Engineering Agent Systems

### 2.1 Devin (Cognition Labs)

**Architecture:** Autonomous AI software engineer with integrated development environment.

**Key Components:**
- **Planning Module:** Decomposes tasks into steps
- **Code Editor:** Full IDE-like environment
- **Terminal:** Shell access for execution
- **Browser:** Web browsing for research and testing

**Orchestration Pattern:**
- Single-agent with multiple tool interfaces
- Self-correcting loop: plan → code → test → debug → repeat
- Maintains context across long-running tasks
- Can deploy and verify solutions

**Sandboxing:**
- Runs in isolated containerized environment
- File system isolation
- Network access controls
- Resource limits (CPU, memory, time)

---

### 2.2 SWE-agent

**Architecture:** Research-oriented agent with Agent-Computer Interface (ACI).

**Core Components:**

1. **Agent Loop:**
   - Thought → Action → Observation cycle
   - Configurable via YAML files
   - Supports multiple LLMs (GPT-4o, Claude, etc.)

2. **Agent-Computer Interface (ACI):**
   - Custom shell commands optimized for LLM interaction
   - File viewing, searching, editing primitives
   - Search and replace operations
   - Context window management

3. **History Processing:**
   - Truncation strategies for long conversations
   - Demonstration injection
   - Thought/action parsing

4. **SWE-ReX (Runtime Environment):**
   - Remote execution of agent commands
   - Docker-based sandboxing
   - Supports local and remote environments

**Key Design Decisions:**
- YAML-driven configuration
- Minimal abstractions (100 lines for mini-swe-agent)
- Research-first: easy to modify and experiment
- State-of-the-art on SWE-bench (65% with mini-swe-agent)

---

### 2.3 OpenHands (formerly OpenDevin)

**Architecture:** Self-hosted developer control center for coding agents.

**Core Components:**

1. **Agent Canvas (Frontend):**
   - Web-based UI for agent interaction
   - Supports multiple backend agents
   - Automation creation and management

2. **Agent Server:**
   - REST API for running multiple agents
   - Single host/port per server
   - Can connect to multiple Agent Servers

3. **Automation Server:**
   - Scheduled agent execution
   - Event-driven triggers (webhooks)
   - Integration with Slack, GitHub, Linear, etc.

4. **Runtime Backends:**
   - Local execution (direct filesystem access)
   - Docker containers (sandboxed)
   - Virtual machines (full isolation)
   - Cloud infrastructure (OpenHands Cloud/Enterprise)

**Agent Types:**
- OpenHands Agent (built-in)
- Claude Code integration
- Codex integration
- Gemini integration
- Any ACP-compatible agent

**Key Design Decisions:**
- Agent-Client Protocol (ACP) for agent interoperability
- Multiple backend support from single frontend
- Self-hosted by default (privacy-first)
- Open-source with commercial cloud option

---

### 2.4 MetaGPT

**Architecture:** Software company simulation with role-based agents.

**Core Philosophy:** `Code = SOP(Team)` - standard operating procedures applied to LLM teams.

**Agent Roles:**
- Product Manager
- Architect
- Project Manager
- Engineer(s)
- QA Engineer

**Orchestration Pattern:**
- Sequential pipeline: Requirements → Design → Implementation → Testing
- Role-based communication (each role has specific responsibilities)
- Structured outputs at each stage (user stories, APIs, code, tests)

**Key Components:**

1. **Role System:**
   - Each role has defined actions and responsibilities
   - Roles can observe and react to others' outputs
   - Message passing between roles

2. **SOP Engine:**
   - Configurable workflow sequences
   - Phase-based execution
   - Quality gates between phases

3. **Workspace Management:**
   - Organized project structure
   - Version control integration
   - Incremental development support

**MGX (MetaGPT X):**
- Commercial product for AI agent development teams
- Natural language programming interface
- #1 Product of the Week on ProductHunt (March 2025)

---

### 2.5 ChatDev

**Architecture:** Virtual software company with communicative agent collaboration.

**Evolution:**
- ChatDev 1.0 (Legacy): Virtual software company
- ChatDev 2.0 (DevAll): Zero-code multi-agent platform

**ChatDev 1.0 Architecture:**

1. **Agent Roles:**
   - CEO, CTO, Programmer, Tester, Designer
   - Each role participates in specialized "functional seminars"

2. **Chat Chain:**
   - Sequential phases of collaboration
   - Phase-specific agent interactions
   - Structured dialogue patterns

3. **Collaboration Patterns:**
   - Instructor-Assistant dynamic
   - Experience Co-Learning (agents learn from past interactions)
   - Iterative Experience Refinement

**ChatDev 2.0 (DevAll) Architecture:**

1. **Zero-Code Platform:**
   - YAML/JSON configuration for agent definitions
   - Visual workflow builder
   - No coding required

2. **Multi-Agent Orchestration:**
   - Puppeteer pattern: Central orchestrator activates and sequences agents
   - Reinforcement learning for orchestration optimization
   - Dynamic, context-aware reasoning paths

3. **MacNet (Multi-Agent Collaboration Networks):**
   - Directed acyclic graph (DAG) topology
   - Supports 1000+ agents without context overflow
   - Linguistic interactions between agents
   - Beyond software development (reasoning, analysis, generation)

**Key Innovations:**
- Experiential Co-Learning: Agents accumulate "shortcut-oriented experiences"
- Docker support for safe execution
- Git integration for version control
- Human-in-the-loop via reviewer role

---

## 3. Message Bus and Event-Driven Patterns

### 3.1 Pattern Taxonomy

| Pattern | Description | Examples |
|---------|-------------|----------|
| **Direct Messaging** | Agent A sends message directly to Agent B | SWE-agent, simple chatbots |
| **Publish/Subscribe** | Agents publish events; subscribers receive relevant ones | AutoGen topics, event-driven systems |
| **Message Queue** | Asynchronous message passing with durability | RabbitMQ, Kafka, Redis Streams |
| **State Machine** | Transitions between states trigger actions | LangGraph StateGraph, CrewAI Flows |
| **Actor Model** | Agents as actors with mailboxes | AutoGen Core, Akka-style |
| **Pub/Sub Bus** | Central event bus for all communications | Apache Kafka, NATS, Redis Streams |

### 3.2 Implementation Examples

**AutoGen (Actor Model + Pub/Sub):**
```
Agent A --[topic:research]--> Message Bus --[topic:research]--> Agent B
                                         --[topic:research]--> Agent C
```
- Topic-based routing
- Typed message schemas
- Decoupled producers and consumers

**LangGraph (State Machine):**
```
[Node A] --state--> [Node B] --state--> [Node C]
    ^                            |
    +----[conditional edge]------+
```
- Shared state object
- Reducer functions for concurrent updates
- Checkpointing at each transition

**CrewAI (Role-Based Delegation):**
```
Manager Agent --[task assignment]--> Worker Agent A
                                  --> Worker Agent B
                                  --> Worker Agent C
Worker Agent A --[result]--> Manager Agent
```
- Dynamic task delegation
- Sequential/parallel/hierarchical processes
- Context passing between tasks

**MetaGPT (SOP Pipeline):**
```
PM --[requirements]--> Architect --[design]--> Engineer --[code]--> Tester
```
- Sequential phase transitions
- Structured artifacts between phases
- Quality gates

### 3.3 Emerging Standards

1. **A2A (Agent-to-Agent Protocol):**
   - Google-led initiative for agent interoperability
   - HTTP-based with JSON-RPC
   - Agent discovery and capability negotiation
   - Adopted by Microsoft Agent Framework

2. **MCP (Model Context Protocol):**
   - Anthropic-led standard for agent-tool communication
   - JSON-RPC based
   - Tool discovery and invocation
   - Resource access patterns

3. **ACP (Agent-Client Protocol):**
   - For agent-to-client communication
   - Adopted by OpenHands
   - RESTful interface
   - Agent discovery and management

---

## 4. Agent Memory and Context Management

### 4.1 Memory Taxonomy

| Type | Scope | Duration | Examples |
|------|-------|----------|----------|
| **Working Memory** | Current task | Session | Conversation context, scratchpad |
| **Short-Term Memory** | Recent interactions | Minutes-hours | Conversation history buffer |
| **Long-Term Memory** | Persistent knowledge | Permanent | Vector stores, databases |
| **Episodic Memory** | Past experiences | Permanent | Interaction logs, trajectories |
| **Semantic Memory** | Distilled knowledge | Permanent | Fact stores, knowledge graphs |

### 4.2 Implementation Patterns

**LangGraph:**
- Thread-level state for conversational memory
- Checkpointer for durable execution
- Long-term memory via store abstraction
- Human-in-the-loop state inspection

**CrewAI:**
- Short-term memory: Current task context
- Long-term memory: Persistent across executions
- Entity memory: Knowledge about entities mentioned
- Context sharing between tasks in a crew

**AutoGen:**
- Topic-based memory organization
- Message history per conversation
- Distributed runtime for shared state

**SWE-agent:**
- History truncation strategies
- Demonstration injection
- Context window management

**OpenHands:**
- Session persistence
- Agent state management
- Cross-session continuity

### 4.3 Key Techniques

1. **RAG (Retrieval-Augmented Generation):**
   - Pull relevant past context from vector stores
   - Semantic search for related memories
   - Ranking by relevance

2. **Summarization:**
   - Compress long histories into summaries
   - Extract key facts and decisions
   - Preserve important context

3. **Hierarchical Memory:**
   - Tiered storage: buffer → working → long-term
   - Automatic promotion/demotion
   - Importance-based retention

4. **KV-Cache Persistence:**
   - Cache LLM activations across sessions
   - Avoid reprocessing unchanged context
   - Reduce latency and cost

### 4.4 Challenges

- **Context window limits:** Even with 1M+ token windows, noise vs. relevance tradeoff
- **Memory staleness:** Conflicting or outdated information over time
- **Privacy:** Sensitive data in persistent memories
- **Retrieval quality:** Finding truly relevant past context

---

## 5. Workspace Isolation and Sandboxing

### 5.1 Isolation Levels

| Level | Mechanism | Isolation | Overhead | Use Cases |
|-------|-----------|-----------|----------|-----------|
| **Process** | Separate processes | Low | Minimal | Simple agents |
| **Container** | Docker/containerd | Medium | Low | Most agents |
| **MicroVM** | Firecracker/gVisor | High | Medium | Multi-tenant |
| **VM** | Full virtualization | Very High | High | Maximum security |
| **Browser** | WebAssembly/sandboxed iframe | Medium | Low | Web agents |

### 5.2 Implementation Examples

**OpenHands:**
- **Docker Sandbox:** Containerized agent execution
  - Volume mounts for project access
  - Network isolation options
  - Resource limits (CPU, memory)
- **VM Backend:** Full virtual machine isolation
  - For untrusted code execution
  - Cloud infrastructure support
- **Local Execution:** Direct filesystem access (with warning)

**SWE-agent (SWE-ReX):**
- Docker-based sandboxing
- Remote execution environment
- Network isolation
- Resource constraints

**Devin:**
- Containerized development environment
- File system isolation
- Network access controls
- Time and resource limits

**ChatDev:**
- Docker support for safe execution
- Isolated workspace per project
- Git integration for version control

### 5.3 Sandboxing Technologies

1. **Docker/containerd:**
   - Most common approach
   - Good balance of isolation and performance
   - Volume mounts for controlled file access

2. **Firecracker (AWS Lambda):**
   - MicroVM-based isolation
   - Fast startup (~125ms)
   - Used by AWS Lambda, Fly.io

3. **gVisor (Google):**
   - User-space kernel
   - Container-level sandboxing
   - Used by Google Cloud Run

4. **E2B:**
   - AI-focused sandboxing platform
   - SDK for agent code execution
   - Pre-built sandboxes for common tasks

5. **Browser Sandboxes:**
   - WebAssembly for safe code execution
   - Sandboxed iframes for web agents
   - Network isolation via CSP

### 5.4 Best Practices

1. **Principle of Least Privilege:**
   - Grant minimal necessary permissions
   - Read-only filesystem where possible
   - No network access unless required

2. **Resource Limits:**
   - CPU time limits
   - Memory caps
   - Disk space quotas
   - Execution timeouts

3. **Audit Logging:**
   - Log all agent actions
   - Track file system changes
   - Monitor network activity

4. **Cleanup:**
   - Destroy sandboxes after use
   - Wipe sensitive data
   - Release resources promptly

---

## 6. Architectural Comparison Matrix

| Framework | Orchestration | State Mgmt | Memory | Sandboxing | Protocol | Maturity |
|-----------|---------------|------------|--------|------------|----------|----------|
| **AutoGen/MAF** | Actor model, graph-based | Checkpointing, distributed | Topic-based, distributed | Runtime-dependent | A2A, MCP | Production (MAF) |
| **LangGraph** | State machine (DAG) | Checkpointer, threads | Thread + store | None built-in | LangChain ecosystem | Production |
| **CrewAI** | Role-based delegation | Flow state, crew context | Short/long/entity | None built-in | Custom | Production |
| **SWE-agent** | Agent loop (ACI) | History buffer | Truncation + demos | Docker (SWE-ReX) | Custom | Research |
| **OpenHands** | ACP-based | Session persistence | Session state | Docker, VM, Cloud | ACP | Beta/Production |
| **MetaGPT** | SOP pipeline | Phase artifacts | Workspace files | Docker | Custom | Production |
| **ChatDev** | Chat chain / Puppeteer | Phase context | Experience co-learning | Docker | Custom | Production |

---

## 7. Key Design Patterns and Trends

### 7.1 Dominant Patterns (2025-2026)

1. **Graph-Based Orchestration:**
   - LangGraph pioneered this approach
   - Microsoft Agent Framework adopting graph patterns
   - Explicit control flow over implicit delegation

2. **Actor Model:**
   - AutoGen's foundation
   - Good for distributed systems
   - Message-passing between isolated agents

3. **Role-Based Collaboration:**
   - CrewAI, MetaGPT, ChatDev
   - Simulates human team dynamics
   - Specialized agents with defined responsibilities

4. **Tool-Use Agents:**
   - Claude, GPT-4, Gemini
   - Agents equipped with external tools
   - Think → Act → Observe loops

### 7.2 Emerging Trends

1. **Protocol Standardization:**
   - A2A (Google) for agent-to-agent
   - MCP (Anthropic) for agent-to-tool
   - ACP for agent-to-client
   - Goal: Interoperable multi-agent systems

2. **Durable Execution:**
   - Checkpointing and resumption
   - Long-running agents (hours/days)
   - Fault tolerance and recovery

3. **Human-in-the-Loop:**
   - Approval checkpoints
   - State inspection and modification
   - Graduated autonomy

4. **Observability:**
   - OpenTelemetry integration
   - Distributed tracing
   - Agent trajectory visualization

5. **Multi-Provider Support:**
   - Framework-agnostic LLM integration
   - Provider failover and load balancing
   - Cost optimization across models

### 7.3 Architecture Selection Guide

| Use Case | Recommended Framework | Rationale |
|----------|----------------------|-----------|
| Complex workflows with branching | LangGraph | Explicit state machine control |
| Rapid prototyping | CrewAI | High-level abstractions |
| Enterprise production | Microsoft Agent Framework | Enterprise support, A2A |
| Software engineering research | SWE-agent | Minimal abstractions, research-friendly |
| Self-hosted coding assistant | OpenHands | ACP interoperability, multiple backends |
| Software company simulation | MetaGPT | SOP-based workflows |
| Zero-code agent orchestration | ChatDev 2.0 | Visual workflow builder |

---

## References

1. **AutoGen/MAF:** https://github.com/microsoft/autogen, https://github.com/microsoft/agent-framework
2. **LangGraph:** https://github.com/langchain-ai/langgraph
3. **CrewAI:** https://github.com/crewAIInc/crewAI
4. **SWE-agent:** https://github.com/SWE-agent/SWE-agent
5. **OpenHands:** https://github.com/OpenHands/OpenHands
6. **MetaGPT:** https://github.com/geekan/MetaGPT
7. **ChatDev:** https://github.com/OpenBMB/ChatDev
8. **A2A Protocol:** https://github.com/google/A2A
9. **MCP:** https://modelcontextprotocol.io/
10. **ACP:** https://agentclientprotocol.org/

---

*Research compiled from official GitHub repositories, documentation, and source code analysis.*
