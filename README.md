# Go API Lab

<img width="1024" height="559" alt="image" src="https://github.com/user-attachments/assets/452f460a-24a8-4363-a6ad-1c96892b0b2f" />


A simple, lightweight Go web server designed to serve static HTML pages with support for both HTTP and HTTPS.

## Features

- **Static Page Serving:** Automatically routes requests to HTML files in the `Pages/` directory.
- **Dual Protocol Support:** Seamlessly switches between HTTP and HTTPS based on configuration.
- **JSON Configuration:** Easy setup via `conf.json`.
- **404 Handling:** Custom 404 page support.
- **Basic Caching:** Infrastructure for caching page paths.

## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.26.2 or later recommended)
- [OpenSSL](https://www.openssl.org/) (optional, for generating self-signed certificates)

## Getting Started

### 1. Installation

Clone the repository and navigate to the project directory:

```bash
git clone <repository-url>
cd api
```

### 2. Configuration

Create or modify `conf.json` in the root directory:

```json
{
  "port": 8080,
  "certFile": "cert.pem",
  "certKey": "key.pem"
}
```

*Note: Leave `certFile` and `certKey` empty to run as a standard HTTP server.*

### 3. Generating SSL Certificates (Optional)

To test HTTPS locally, generate a self-signed certificate:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes -subj "/CN=localhost"
```

### 4. Running the Server

```bash
go run main.go
```

The server will start on the port specified in `conf.json`.

## Project Structure

- `main.go`: The core server logic and request handling.
- `conf.json`: Configuration settings for the server.
- `Pages/`: Directory containing `.html` files (e.g., `main.html`, `test.html`, `404.html`).
- `.gitignore`: Standard Go and environment exclusions.

## Development

To add a new page, simply create a new `.html` file inside the `Pages/` directory. For example, creating `Pages/about.html` will make it accessible at `http://localhost:port/about`.
