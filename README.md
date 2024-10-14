# Muzz Explore Service

## Overview

The Explore Service is a **gRPC-based service** designed for a dating app. It manages and retrieves user interactions, focusing on users who have "liked" a recipient. The service includes methods for listing users who liked a recipient, counting the number of likes, and recording user decisions (e.g., liking or not liking someone back).

## Why Redis?

**Redis** was chosen as the data store due to its high performance and efficiency in handling real-time interactions. Redis's in-memory structure allows for rapid read/write operations, making it ideal for a dating app where quick responses are crucial. Redis also supports a variety of data types, enabling easy management of complex relationships such as mutual likes and user decisions.

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your machine.

### Running the Application

To run the Explore Service with Docker Compose, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://skartikey: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@github.com/skartikey/muzz-explore-service.git
   cd muzz-explore-service
   ```

2. Build and run the service using Docker Compose:
   ```bash
   docker-compose up --build
   ```

3. The service will be available at `localhost:50051` (or your specified host and port).

You can make requests to the gRPC server using tools like **BloomRPC**, **Postman** (with gRPC support), or **grpcurl**. Below are some example requests:

### Example gRPC Requests:

- **ListLikedYou**:
   ```json
   {
     "recipient_user_id": "1",
     "pagination_token": ""
   }
   ```
- **ListNewLikedYou**:
   ```json
   {
     "recipient_user_id": "1",
     "pagination_token": ""
   }
   ```
- **CountLikedYou**:
   ```json
   {
     "recipient_user_id": "1"
   }
   ```
- **PutDecision**:
   ```json
   {
     "actor_user_id": "1",
     "recipient_user_id": "4",
     "liked_recipient": true
   }
   ```

### Test Data Initialization

The service initializes with mock data when it starts. In the `cmd/main.go` file, there is a **Redis Initialization Script** that populates the database with sample data. This is useful for testing and verifying functionality during development.

#### Example Test Data:

- Users who have liked specific recipients.
- Mutual likes between users.
- User decisions regarding interactions (e.g., like/dislike).

This test data provides a foundation for testing the service's API endpoints.

## API Endpoints

The Explore Service exposes the following gRPC methods:

1. **ListLikedYou**: Retrieves a list of users who liked a specific recipient.
2. **ListNewLikedYou**: Retrieves a list of users who recently liked a specific recipient.
3. **CountLikedYou**: Returns the total number of users who liked a specific recipient.
4. **PutDecision**: Records a user's decision regarding a like (e.g., whether they liked the other person back).

For detailed usage and documentation, refer to the generated gRPC protobuf files.

## Performance Benchmarks

The performance of the Explore Service has been benchmarked using Go tests. Below are the results on a machine with the following specifications:

- **OS**: Linux
- **Architecture**: AMD64
- **Processor**: 11th Gen Intel(R) Core(TM) i7-11850H @ 2.50GHz

```bash
go test -bench=.
```

### Benchmark Results:

| Benchmark Method                     | Ops/Second | Time per Operation (ns/op) |
|--------------------------------------|------------|-----------------------------|
| **BenchmarkExploreService_ListLikedYou**  | 42,243     | 31,022 ns/op                |
| **BenchmarkExploreService_ListNewLikedYou** | 36,927     | 29,011 ns/op                |
| **BenchmarkExploreService_CountLikedYou**  | 44,943     | 27,471 ns/op                |
| **BenchmarkExploreService_PutDecision**    | 21,208     | 55,309 ns/op                |

The benchmarks reflect the service's efficiency, with consistently low operation times for listing and counting likes, as well as recording decisions.

## Scalability

The Explore Service is designed with scalability in mind, using Redis as the primary data store.

- **Horizontal Scaling**: As the number of user interactions grows, additional service instances can be deployed to handle increased load. Redis can be configured in clustered mode to distribute the data across multiple nodes.

- **Caching Layer**: Redis serves as a caching layer, reducing the load on primary databases. This enables fast access to frequently requested data, such as likes and user decisions.

- **Pub/Sub System**: Redis Pub/Sub capabilities can be leveraged to enable real-time updates and notifications. For example, users can be notified of new likes or mutual matches in real time.

## Future Improvements

Potential future enhancements include:

- **Rate Limiting**: Implement rate limiting to prevent abuse, especially for high-frequency actions like liking users.
- **Analytics and Monitoring**: Integrate logging, monitoring, and analytics tools to track performance and user interactions, enabling optimizations and easier debugging.
- **Data Expiration Policies**: Introduce policies to expire old data (e.g., likes and decisions) to ensure the service remains responsive in a dynamic environment like a dating app.
