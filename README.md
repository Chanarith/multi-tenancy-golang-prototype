## Ecommerce SAAS

This project is a simplified multi-tenant ecommerce web application api that supports 1 database per tenant. It is designed to demonstrate the data separation this approach gives and highlights the theory of enterprise applications development.

## The Stack

- **Go**
- **Postgresql**
- **Redis**
- **Supertokens**
- **Tools**: [web-framework: Gin Gonic, cli: Cobra, orm: Gorm, queue: go-craft]

## Commands

Here are some commands available in this project:

`go run core {command1} {-f1}` where f1 is a flag for command1

- `migrate`: Migrates the database. Runs the migrations for the central database by default.
  - `-t` or `--tenants`: Migrates all databases of existing tenants.
- `run_server`: Runs the Go API server

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
