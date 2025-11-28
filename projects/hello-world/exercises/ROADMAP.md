# 🚀 Go Learning Roadmap

**Goal**: Become confident in Go application development (1.5 YoE equivalent)

**Your Background**: Java/Spring Boot, Redis, Cassandra, Hibernate

---

## 📈 Learning Path

### Phase 1: Go Fundamentals (Week 1-2)
| # | Exercise | Status | Key Concepts |
|---|----------|--------|--------------|
| 01 | Interfaces & Polymorphism | ✅ Complete | Implicit interfaces, method receivers |
| 02 | Error Handling | ✅ Complete | Error types, wrapping, custom errors |
| 03 | Goroutines & Channels | ✅ Complete | Concurrency, channel patterns |
| 04 | Context & Cancellation | ✅ Complete | Context propagation, timeouts |
| 05 | Testing in Go | ✅ Complete | Table-driven tests, mocking |

### Phase 2: Web Development (Week 3-4)
| # | Exercise | Status | Key Concepts |
|---|----------|--------|--------------|
| 06 | HTTP Server Basics | 🔄 In Progress | net/http, handlers, middleware |
| 07 | REST API with Routing | ⏳ Pending | Chi/Gin router, JSON handling |
| 08 | Dependency Injection | ⏳ Pending | Wire, manual DI (no Spring magic!) |
| 09 | Configuration Management | ⏳ Pending | Viper, env files |
| 10 | Logging & Observability | ⏳ Pending | Zap/Zerolog, structured logging |

### Phase 3: Data Layer (Week 5-6)
| # | Exercise | Status | Key Concepts |
|---|----------|--------|--------------|
| 11 | Database with GORM | ⏳ Pending | ORM similar to Hibernate |
| 12 | Raw SQL with sqlx | ⏳ Pending | When ORM is overkill |
| 13 | Redis Integration | ⏳ Pending | go-redis, caching patterns |
| 14 | Connection Pooling | ⏳ Pending | Managing DB connections |

### Phase 4: Production Ready (Week 7-8)
| # | Exercise | Status | Key Concepts |
|---|----------|--------|--------------|
| 15 | Graceful Shutdown | ⏳ Pending | Signal handling, cleanup |
| 16 | Health Checks | ⏳ Pending | Liveness, readiness probes |
| 17 | Docker & Deployment | ⏳ Pending | Multi-stage builds |
| 18 | Complete Microservice | ⏳ Pending | Put it all together! |

---

## 🎯 Java → Go Mental Model Shifts

| Java Concept | Go Equivalent | Key Difference |
|--------------|---------------|----------------|
| `class` | `struct` | No inheritance, use composition |
| `implements` | (implicit) | Just define the methods |
| `extends` | Embedding | Composition over inheritance |
| `try/catch` | `if err != nil` | Errors are values, not exceptions |
| `@Autowired` | Constructor injection | No magic, explicit wiring |
| `Thread` | `goroutine` | Lightweight, thousands OK |
| `synchronized` | `channels` / `mutex` | "Share by communicating" |
| `Optional<T>` | Multiple returns | `value, ok := map[key]` |
| `Stream API` | `for` loops | Go prefers explicit loops |
| Annotations | Code generation | `go generate`, no reflection magic |

---

## 📊 Progress Tracker

- [ ] Phase 1: Go Fundamentals
- [ ] Phase 2: Web Development
- [ ] Phase 3: Data Layer
- [ ] Phase 4: Production Ready
- [ ] 🏆 Build a complete microservice!

---

**Let's start with Exercise 01!** 🎉

