# Soccerbuddy

> **⚠️ Warning: This project is highly WIP and unstable. The events and event journal format will change without warning.**

## Overview

This project is an application designed for soccer clubs (and in the future potentially other sports clubs) to manage their teams, trainings, matches, and more. It aims to streamline club operations by providing tools for scheduling, communication, and record-keeping.

The application consists of:

- **Server Code**: Written in Go, implementing an event-sourced architecture with PostgreSQL as the event store.
- **Web Admin Frontend**: Developed with Svelte and TypeScript for administration and management tasks.
- **Redis with RedisJSON/RedisSearch**: Used for providing super-fast in-memory projections.
- **Permify Integration**: Handles relationship-based access control by adjusting tuples in response to certain events.
- **Connect RPC**: Facilitates communication between the frontend and backend using Protobuf and gRPC-web.

*Note: A mobile app will soon be developed using React Native.*

## Architecture

- **Event Store**: All domain events are stored in PostgreSQL.
- **Projections**: Real-time projections are maintained in Redis using RedisJSON and RedisSearch for fast read access.
- **Access Control**: We leverage Permify to handle complex relationship-based permissions. A dedicated projection listens to specific events and updates the tuples in Permify accordingly.
- **Communication**: The frontend communicates with the server using Connect RPC, which utilizes Protobuf and gRPC-web for efficient and type-safe communication.
- **Frontend**: The web admin interface interacts with the server via Connect RPC.

## Dependencies

- **Go**: Backend server implementation.
- **PostgreSQL**: Event store database.
- **Redis with RedisJSON/RedisSearch**: In-memory data storage for projections.
- **Permify**: Relationship-based access control system.
- **Buf**: Generating code and stubs based on the Protobuf definitions.
- **Node.js and pnpm**: For the frontend development with Svelte and TypeScript.

## Credits

- [**Zitadel**](https://github.com/zitadel/zitadel): For open sourcing, showcasing, and documenting their excellent event sourcing system that inspired this one.
- [**Akka**](https://akka.io/): Providing additional insights on building event-sourced systems.

## License

This project is licensed under the AGPL v3 License.
