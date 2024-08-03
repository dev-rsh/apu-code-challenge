# Code challenge

This project is developed for code challenge by Arian Roshanzamir. This project is responsible for spinning up 3 services implementing grpc servers using proto files in the `proto` directory.

## Project Structure

The project is organized into the following directories:

- `proto`: Contains .proto files
- `service1`: Implementation of Service1
- `service2`: Implementation of Service2
- `service3`: Implementation of Service3

Each service has a `config` folder to load configs from `configs.yaml` file and a `deployment` folder for holding files necessary for deployment. 

The directory structure of Service3 is as follows:

- `dto`: Contains Data Transfer Object between layers
- `Infrastructure`: Contains repository and server implementations

## Getting Started

1. Clone this repository.
2. `cd` to the root of the project.
3. Run `docker-compose up -d` to run the containers in detached mode. 3 Golang containers spin up alongside a PostgreSQL container with the table for service3 already in it.
4. Now  service3 is running on port `50053` and exposes it.
5. PostgreSQL runs on port `5434` and exposes it.
6. To run tests, run `docker-compose exec service3 sh` and when you are inside the container run the following command `go test -C $(pwd) ./...`.