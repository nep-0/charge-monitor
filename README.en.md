This project aims to provide an EV charging station status monitoring service that can periodically query the status of charging stations and provide real-time data through an HTTP interface. Below is a brief description of the project:

## Project Structure

- **Dockerfile**: Provides Docker configuration for building and running the project.
- **app/app.go**: Main program logic, including HTTP service, data updating, and cross-origin middleware.
- **cache/local_cache.go**: Implements local caching functionality for storing and querying EV charging station status information.
- **config/config.go**: Configuration file parsing logic, supporting loading configurations from files.
- **query/query.go**: Provides the core functionality for querying EV charging station status.
- **main.go**: Program entry point, used to start the service.

## Features

- **Periodic Polling of Charging Station Status**: Automatically updates charging station status information based on the configured polling interval.
- **HTTP Service**: Provides a RESTful API to query real-time status of charging stations.
- **Caching Mechanism**: Uses local caching to store charging station status, improving query efficiency.
- **Configuration Support**: Supports loading parameters such as charging station addresses and polling intervals from configuration files.
- **Cross-Origin Support**: Supports cross-origin requests via middleware, making it easier for front-end applications to access.

## Usage

### Build and Run

1. **Build Docker Image**:
   ```bash
   docker build -t charge-monitor .
   ```

2. **Run Docker Container**:
   ```bash
   docker run -d -p 8000:8000 charge-monitor
   ```

### Configuration

- Modify the `config.yaml` file to configure charging station addresses and polling intervals.

### API Interface

- **Get Charging Station Status**:
  - **URL**: `/outlets`
  - **Method**: `GET`
  - **Response**: Returns the current status information of all charging stations.

## Development and Testing

- **Unit Tests**: Unit tests for caching and querying functionality are provided in `cache/local_cache_test.go` and `query/query_test.go`.
- **API Testing**: The `/outlets` interface can be tested using tools such as `curl` or Postman.

## License

This project is licensed under the MIT License. For details, please refer to the LICENSE file in the project.

For further information about the project or for custom development, please refer to the source code and related documentation.