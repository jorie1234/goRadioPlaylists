# goRadioPlaylists

A Go application that fetches radio playlists from onlineradiobox.com, stores them in a MongoDB database, and serves them as JSON files.

## Features

- Fetches playlist data for multiple German radio stations from `onlineradiobox.com` (radiohannover, radioffn, njoyradio, 80s80s).
- Stores playlist data in a MongoDB database.
- Merges new playlist data with existing data to create a comprehensive history.
- Saves the full playlist history to local JSON files.
- Serves the JSON files via a simple web server.
- Runs as a cron job to periodically update the data every 30 minutes.
- Containerized using Docker for easy deployment.

## Getting Started

This project is designed to be run with Docker and Docker Compose.

### Prerequisites

- Docker
- Docker Compose

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/jorie1234/goRadioPlaylists.git
    cd goRadioPlaylists
    ```

2.  This project uses [Traefik](https://traefik.io/) as a reverse proxy. The provided `docker-compose.yml` is configured to work with a Traefik instance running on the same Docker network. You will need to have a Traefik container running and a `gateway` network created.

    If you don't have a Traefik setup, you can use the following `docker-compose.yml` for Traefik:
    ```yml
    version: '3'

    services:
      traefik:
        image: traefik:v2.2
        container_name: traefik
        command:
          - "--api.insecure=true"
          - "--providers.docker=true"
          - "--providers.docker.exposedbydefault=false"
          - "--entrypoints.web.address=:80"
        ports:
          - "80:80"
          - "8080:8080"
        volumes:
          - /var/run/docker.sock:/var/run/docker.sock:ro
        networks:
          - gateway

    networks:
      gateway:
        external: true
    ```
    You will also need to create the `gateway` network:
    ```sh
    docker network create gateway
    ```

3.  The `docker-compose.yml` file uses a `DOMAINNAME` environment variable to configure the hostname for the Traefik frontend rule. You will need to set this variable. For local development, you can add it to a `.env` file in the root of the project:
    ```
    DOMAINNAME=localhost
    ```

4.  Run the application using Docker Compose:
    ```sh
    docker-compose up -d
    ```

## Configuration

The application can be configured using environment variables.

| Variable             | Description                                                                 | Default                                 |
| -------------------- | --------------------------------------------------------------------------- | --------------------------------------- |
| `DOMAINNAME`         | The domain name for the Traefik frontend rule.                              | `localhost`                             |
| `MONGODB_USERNAME`   | The username for the MongoDB database. (Currently not used in the code)     |                                         |
| `MONGODB_PASSWORD`   | The password for the MongoDB database. (Currently not used in the code)     |                                         |
| `MONGODB_ENDPOINT`   | The endpoint for the MongoDB database. (Currently not used in the code)     |                                         |

Create a `.env` file in the root of the project to set these variables:
```
DOMAINNAME=yourdomain.com
```

## Usage

Once the application is running, the playlist data can be accessed via the web server. The server exposes the contents of the `data` directory under the `/radiodata` route.

For example, to access the playlist for `radiohannover`, you can use the following URL:

```
http://radio.yourdomain.com/radiodata/radiohannover.json
```

Replace `yourdomain.com` with the value of your `DOMAINNAME` environment variable.

The available playlists are:
- `radiohannover.json`
- `radioffn.json`
- `njoyradio.json`
- `80s80s.json`
