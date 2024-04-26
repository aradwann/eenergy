# eEnergy

![Go workflow](https://github.com/aradwann/eenergy/actions/workflows/test.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/aradwann/eenergy/badge.svg?branch=main)](https://coveralls.io/github/aradwann/eenergy?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/aradwann/eenergy)](https://goreportcard.com/report/github.com/aradwann/eenergy)

eEnergy is a cutting-edge platform designed to revolutionize the way individuals and businesses exchange and source energy. By leveraging the power of renewable energy technology, eEnergy aims to empower households and energy producers to not only generate but also sell and exchange energy units efficiently. This platform facilitates users in finding the nearest energy providers or generators, making the transition to renewable energy sources seamless and convenient.

## Key Features

- **Energy Exchange Platform**: Enables users to buy, sell, and exchange energy units with ease.
- **Find Nearby Energy Resources**: Utilizes advanced algorithms to help users locate the nearest energy providers or generators.
- **User Verification**: Implements a robust verification process for new users, enhancing security and trust.
- **High-Performance Backend**: Built with Go and gRPC for a fast, efficient, and scalable backend service.
- **Job Queueing**: Leverages Redis for queueing jobs, such as sending verification emails to newly registered users.
- **Containerization**: Utilizes Docker for easy deployment and scaling.
- **Observability**: Adopt OpenTelemetry observability standard that incorporates Loki, Promtail for logs, Jaeger for traces, Prometheus for metrics, and Grafana for a comprehensive real-time monitoring dashboard.

- **Reliable Data Storage**: Uses PostgreSQL as the database, enhanced with stored procedures for efficient data management.

## Tools to Install

To contribute to eEnergy, you will need to set up your development environment with the following tools:

- **Buf**: Provides Protocol Buffer support and linting in VSCode. Essential for working with gRPC services.
- **Migrate**: A CLI tool for handling database migrations. Ensures that your local development database schema is up-to-date.
- **Protocol Buffer Compiler and Plugins**: Necessary for compiling `.proto` files into Go code.
  - Install on Fedora/RHEL/CentOS:

    ```bash
    sudo dnf install protobuf-compiler protobuf-devel
    ```

## Getting Started

1. **Clone the Repository**

    ```bash
    git clone https://github.com/aradwann/eenergy.git
    ```

2. **Set Up Your Local Development Environment**
    - Ensure all required tools are installed.

3. **Generate Certificates For Local Use**

- ```bash
      make init
  ```

- ```bash
      make gen-cert
  ```

4.**Build and Run Using Docker**

- Build the Docker images:

    ```bash
    docker compose build
    ```

- Run the containers:

    ```bash
    docker compose up
    ```

### using postman with mTLS enabled

1. add CA pem
2. add client certificate with crt, key files and localhost:9091 as domain

3. **Access the Application**
    - The application and its services are now accessible on your local machine.

## Contributing

We welcome contributions! Please see the `CONTRIBUTING.md` file for more information on how to get involved.

## Reporting Issues

If you encounter any issues or have suggestions for improvements, please file an issue on our GitHub repository.

## License

eEnergy is open-sourced under the [MIT license](LICENSE).

---

By contributing to eEnergy, you're helping to create a more sustainable and efficient future for energy consumption and distribution. Let's make a difference together!

### To-Do List

- [ ] Add account balance endpoint for admin only (or make it smart meter responsibility)
- [ ] Improve logging and configure Grafana
- [ ] Implement trace to log integration
- [ ] Explore pprof for performance profiling
- [ ] Implement nearest energy source endpoint (input: current location, output: nearest available energy source details)
- [ ] Kubernetes vs Docker Swarm
