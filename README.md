
# Golang Messaging Application

`cme_goapp` is a Go-based messaging application designed to handle efficient message exchanges with integrated caching capabilities. This application ensures rapid and reliable communication, leveraging modern technologies for scalability and performance.

## Features

- **Messaging System**: Implements a robust messaging platform.
- **Caching**: Utilizes Redis for caching mechanisms, ensuring quick message retrieval.
- **Database Integration**: Supports Cassandra for scalable data storage.
- **Monitoring**: Includes Prometheus for monitoring and metrics collection.
- **Docker Support**: Provides Docker configurations for easy deployment.

## Prerequisites

Ensure the following are installed on your system:

- [Go](https://golang.org/dl/) (version 1.18 or higher)
- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/emersonary/cme_goapp.git
   ```

2. **Navigate to the project directory**:

   ```bash
   cd cme_goapp
   ```

3. **Build and start the services using Docker Compose**:

   ```bash
   docker-compose up --build
   ```

This command will set up the following services:

- **Cassandra**: Accessible on port 9042.
- **Redis**: Accessible on its default port.
- **Messaging Application**: Manages the core messaging logic.
- **Nginx**: Serves as a reverse proxy, accessible on port 80.

## Configuration

Configuration files are located in the `config` directory. Ensure that the environment variables in the `docker-compose.yaml` file are correctly set:

- `DB_HOST`: Hostname for the Cassandra service.
- `REDIS_HOST`: Hostname for the Redis service.
- `APP_WAITFORSTARTUP`: Time in seconds to wait before the application starts.

## Usage

Once all services are running, the application can be accessed via `http://localhost`.

## Project Structure

- `cache/`: Contains caching mechanisms and configurations.
- `cmd/`: Entry points for the application.
- `config/`: Configuration files.
- `controller/`: Handles incoming requests and routes.
- `database/`: Database connection and queries.
- `dto/`: Data Transfer Objects.
- `error/`: Error handling mechanisms.
- `event/`: Event handling, including Prometheus integration.
- `handler/`: Business logic handlers.
- `init/`: Initialization scripts.
- `log/`: Logging configurations.
- `model/`: Data models.
- `pck/`: Utility packages.
- `webserver/`: Web server configurations.

## Contributing

Contributions are welcome! Please fork the repository, create a feature branch, and submit a pull request with your changes.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

Special thanks to [emersonary](https://github.com/emersonary) for developing this application.
