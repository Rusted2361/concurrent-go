# Go Channels: Buffered vs Unbuffered — Usage in Worker Pool & Fan-In/Fan-Out Patterns

## 1. Channel Basics
Go channels enable safe communication between goroutines.  
They can be:
- **Unbuffered:** No internal capacity. Send blocks until another goroutine receives.
- **Buffered:** Has internal queue of fixed size. Send blocks only when the buffer is full.

---

## 2. Unbuffered Channels
**Behavior:**
- Acts as a *synchronization point* — sender and receiver must meet.
- Guarantees that message is received immediately when sent.

**Use Cases:**
- When strict ordering or synchronization is required.
- When workload is small or predictable.
- Ensures deterministic control flow (safer for correctness).

**Pros:**
- Simpler logic — less risk of hidden race conditions.
- Implicit synchronization between goroutines.
- Low memory usage.

**Cons:**
- Can easily lead to **deadlocks** if sends/receives are not perfectly balanced.
- Reduced throughput due to blocking behavior.
- Workers often idle waiting for messages.

**Example (Worker Pool - Unbuffered):**
```go
jobs := make(chan Job)
results := make(chan Result)

for w := 0; w < 3; w++ {
    go worker(jobs, results)
}

for _, job := range jobList {
    jobs <- job // blocks until a worker picks it
}
close(jobs)
```
➡️ Ensures that no jobs are queued unless a worker is ready.

---

## 3. Buffered Channels
**Behavior:**
- Adds a *message queue* between sender and receiver.
- Sends only block when the buffer is full; receives block when empty.

**Use Cases:**
- When producers are faster than consumers.
- For parallel pipelines or batched workloads.
- When slight decoupling improves throughput.

**Pros:**
- Higher concurrency and throughput.
- Smooths out producer-consumer speed differences.
- Avoids deadlocks in many fan-out patterns.

**Cons:**
- Harder to reason about timing and ordering.
- Risk of memory buildup if buffer size is large or unbounded producers.
- Doesn’t provide natural backpressure if buffer is too big.

**Example (Worker Pool - Buffered):**
```go
jobs := make(chan Job, 10)
results := make(chan Result, 10)

for w := 0; w < 3; w++ {
    go worker(jobs, results)
}

for _, job := range jobList {
    jobs <- job // won’t block until buffer fills
}
close(jobs)
```
➡️ Decouples producer and workers — improves throughput.

---

## 4. Fan-In / Fan-Out with Channels

| Pattern | Description | Channel Type Choice |
|----------|--------------|--------------------|
| **Fan-Out** | Multiple workers read from one channel. | Either works — buffered improves throughput; unbuffered ensures tight sync. |
| **Fan-In** | Multiple producers send to one channel, one consumer reads. | Buffered preferred to prevent producer blocking if consumer is slow. |

**Example (Fan-In using Buffered Channel):**
```go
merged := make(chan string, 100)
for _, src := range sources {
    go func(ch <-chan string) {
        for msg := range ch {
            merged <- msg
        }
    }(src)
}
```
Buffered channel prevents producer goroutines from blocking if consumer lags.

---

## 5. Design Guidelines

| Scenario | Recommended Channel Type | Reason |
|-----------|--------------------------|---------|
| Tight synchronization between goroutines | Unbuffered | Ensures handoff safety and timing accuracy |
| High-throughput worker pools | Buffered | Avoids blocking, improves concurrency |
| Fan-in from multiple producers | Buffered | Prevents blocking producers |
| Fan-out to multiple workers with stable load | Unbuffered | Keeps coordination strict and predictable |
| Burst or uneven workloads | Buffered | Provides a queue to absorb spikes |

---

## 6. Summary

| Aspect | Unbuffered | Buffered |
|--------|-------------|----------|
| Synchronization | Strong | Loose |
| Performance | Lower | Higher (with proper sizing) |
| Risk of Deadlocks | Higher | Lower (but logic complexity increases) |
| Recommended for | Controlled flow, testing | Production pipelines, real workloads |
