# Assignment 2 - Corona case information service

## Overview

## Endpoints

1. /corona/v1/country/
2. /corona/v1/policy/
3. /corona/v1/diag/
4. /corona/v1/notifications/

## Development

This project targets Go 1.15 and 1.16 and I will assume `$GO111MODULE` is set to `on` (or empty if you are running GO 1.16 or newer).

The project uses two frameworks for code consistency and quality.
The first is golangci-lint which I will assume any Go programmer is familiar with.
The other is pre-commit; which is a python based git hook framework for running stuff (golangci-lint in this case) everytime you commit.
These can be installed via the following commands:
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.37.1
# Install pre-commit (also available from traditional package managers)
pip install pre-commit
# Install pre-commit hooks
pre-commit install --install-hooks
```

## Deployment

The server is accessible using the following url: http://10.212.142.242:3000

I chose a slightly different setup then what was shown in the lectures. I went with a fedora based instance, simply because I am much more familiar with fedora than with ubuntu or debian.

### Firebase authentication

In order to deploy this project, you need to generate a service account key and point the environment variable GOOGLE_APPLICATION_CREDENTIALS to it.
The server then picks up on the environment variable, and uses the key to authenticate against Firebase.

I recommend setting up a `.env` file for all your environment variables when developing on projects like this.

### Run as a systemd service

The server is currently deployed on skyhigh / openstack as a systemd service in user mode. This is achieved using the service unit included under `systemd`.
As written the service unit makes a bunch of assumptions about the instance that it will run on; like the user being called fedora and the server being located at `/home/fedora/server`; and will have to be customized before being deployed to any other instance.

Using a systemd service unit provides the following benefits:
1. It provides a relatively standardized way of running the server as a service on linux distributions.
2. The server can be auto-restarted on failure.
3. It enables admins to use systemd's tools to monitor and control the server as a service (like viewing logs and status using `systemctl status server.service`).

The service can be installed using the following commands:
```bash
# Install the service unit
cp <path_to_repo>/systemd/server.service ~/.config/systemd/user/server.service
# Reload units
sudo systemctl daemon-reload
# Enable and start the service
systemctl --user enable --now server.service
# Check the status of the server
systemctl --user status server.service
# Follow the logs from the server as they are written
journalctl --user-unit=server.service --follow
```

## Third party

For this project I used a library called [chi][1], which is a express-like routing library that simplifies specifying endpoints and their routs.

[1]: https://github.com/go-chi/chi
