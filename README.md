# Event-Driven

This repository is a hands-on learning space for Event-Driven Architecture (EDA) using RabbitMQ.

The goal is to understand how systems communicate through events, how message brokers enable asynchronous processing, and how decoupling improves reliability and scalability. The focus is on practical exploration rather than production-ready frameworks.

## What this repository explores

- Event-driven communication
- RabbitMQ fundamentals (exchanges, queues, bindings)
- Routing patterns (direct, fanout, topic)
- Asynchronous processing
- Message-based system design
- Reliability and failure isolation

## Why RabbitMQ

RabbitMQ provides a clear and practical way to understand the core ideas behind Event-Driven Architecture without the operational complexity of large streaming platforms. It is well-suited for learning how events flow through a system and how producers and consumers are decoupled.

## Status

This repository is actively evolving as concepts are explored and refined.

## Learning Notes

The root-level `main.go` is kept on purpose as a small concurrency sandbox.
It is not the main entrypoint for the RabbitMQ example application.

That file exists to practice:

- `errgroup.WithContext`
- goroutine coordination
- cancellation when one task fails
- waiting for concurrent work with `g.Wait()`

The RabbitMQ app entrypoints live in:

- `cmd/producer/main.go`
- `cmd/consumer/main.go`
