# DNS Updater API

Simple Go API to update `dnsmasq.conf` with hostname-to-IP mappings, dockerized.

## Requirements

- Docker

## Installation with Docker Compose

1. **Clone the Repository**

2. **Configure Environment**
   - copy the `.env.exemple` to `.env`
   - adapt the parameters

   ```.env
    AUTH_TOKEN=api-token # The API token used to authenticate with the API
    DNSMASQ_FILE_PATH=dnsmasq.conf # The path to the dnsmasq configuration file that will be used to add the new subdomain
    DOMAIN_BASE=dns.local #The dns base that will be used to append to the subdomain (hostname) from the query to the api
   ```

3. **Run with Docker Compose**

   ```bash
   docker-compose up
   ```

   - Uses image: `ghcr.io/GrimalDev/dns-updater:latest`
   - API runs on `http://localhost:8080`.

## Self-Build

1. **Requirements**
   - Go 1.23+
   - Docker and Docker Compose

2. **Build and Run**
   - Modify `docker-compose.yml` to use `build` instead of `image`:
     ```yaml
     version: "3.8"
     services:
       dns-updater:
         build: .
         ports:
           - "8080:8080"
         volumes:
           - ./dnsmasq.conf:/app/dnsmasq.conf
         env_file:
           - .env
     ```
   - Run:
     ```bash
     docker-compose up --build
     ```

## API Usage

- **Endpoint**: `POST /update-dns`
- **Headers**:
  - `Authorization: your-secret-token` (must match `.env` token)
  - `Content-Type: application/json`
- **Body**:
  ```json
  {
    "ip": "<IP_ADDRESS>",
    "hostname": "<HOSTNAME>"
  }
  ```
- **Response**:
  - Success: `200 OK` with `"DNS updated for <hostname>.nsa.local to <ip>"`
  - Errors: `401 Unauthorized`, `400 Bad Request`, or `500 Internal Server Error`

- **Example**:
  ```bash
  curl -X POST http://localhost:8080/update-dns \
    -H "Authorization: your-secret-token" \
    -H "Content-Type: application/json" \
    -d '{"ip":"192.168.1.100","hostname":"server1"}'
  ```

## Bash Command

- **Script**: `update-dns.sh`
- **Purpose**: Updates `dnsmasq.conf` with the local machine's IP and hostname.
- **Usage**:
  1. Ensure `Authorization` header in `update-dns.sh` matches `.env` `AUTH_TOKEN`:
     ```bash
     #!/bin/bash
     IP=$(hostname -I | awk '{print $1}')
     HOSTNAME=$(hostname)
     curl -X POST http://localhost:8080/update-dns \
       -H "Authorization: your-secret-token" \
       -H "Content-Type: application/json" \
       -d "{\"ip\":\"$IP\",\"hostname\":\"$HOSTNAME\"}"
     ```
  2. Make executable:
     ```bash
     chmod +x update-dns.sh
     ```
  3. Run:
     ```bash
     ./update-dns.sh
     ```

## Notes

- The API updates `/app/dnsmasq.conf` with `address=/<hostname>.<base dns>.local/<ip>`.
- The `dnsmasq.conf` file is persisted in the project directory.
