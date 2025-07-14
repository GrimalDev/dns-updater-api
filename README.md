# DNS Updater API

Simple Go API to update `dnsmasq.conf` with hostname-to-IP mappings, dockerize.

## Requirements

- Docker and Docker Compose
- Go (for local development, optional)
- Bash (for running the command)

## Setup

1. **Clone the Repository**

2. **Configure Environment**
   - Create a `.env` file:
     ```
     AUTH_TOKEN=your-secret-token
     ```
   - Replace `your-secret-token` with a secure token.
   - Ensure an empty `dnsmasq.conf` exists in the project directory.

3. **Build and Run**

   ```bash
   docker-compose up --build
   ```

   - API runs on `:8080`.

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
  1. Ensure the `Authorization` header in `update-dns.sh` matches the `.env` `AUTH_TOKEN`.
  2. Make executable:
     ```bash
     chmod +x update-dns.sh
     ```
  3. Run:
     ```bash
     ./update-dns.sh
     ```

## Notes

- The API updates `/app/dnsmasq.conf` with `address=/<hostname>.nsa.local/<ip>`.
- The `dnsmasq.conf` file is persisted in the project directory.
- Ensure the API is running before executing the Bash script.
- Update `<api-host>` in `update-dns.sh` if not using `localhost`.
