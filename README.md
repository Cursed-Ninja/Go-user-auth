## Setting up the project

Clone the repository in your local machine using the below command.

```bash
git clone https://github.com/Cursed-Ninja/Go-user-auth.git
```

## Setting up the environment variables

Create a config.yml file in the `/user-auth/etc/secrets directory`. Add the following fields and update the values when necessary:

```
google_oauth:
  client_id: "your_client_id"
  client_secret: "your_client_secret"
  redirect_uri: http://localhost:8080/auth/google/callback

database:
  username: "your_username"
  password: "your_password"
  host: "your_host"
  port: "your_port"
  database: "your_database_name"

session:
  name: "your_session_name"
  secret: "your_session_secret"

```

## Starting the server

To start the server, first open the `/user-auth/cmd/main` directory in the terminal and then run the below command.

```bash
go run main.go
```

## Structure

1. The server starts from the main.go file in the `/user-auth/cmd/main` directory.
2. Database connection is establised in the `/user-auth/internal/config/config.go` file.
3. The routes are defined in the `/user-auth/internal/routes/routes.go` file.
4. The controllers are defined in the `/user-auth/internal/controllers` directory.
5. The models are defined in the `/user-auth/internal/models` directory.

## Notes

1. Uncomment the `Line 9` and comment `Line 11` of `user-auth/internal/config/viper/viper.go` file to use the config.yml file.
