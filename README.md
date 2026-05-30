# Go API Lab

<img width="1024" height="559" alt="image" src="https://github.com/user-attachments/assets/452f460a-24a8-4363-a6ad-1c96892b0b2f" />

A simple, lightweight Go web server designed to serve static HTML pages with support for both HTTP and HTTPS, dynamic file-sharing, and JWT authentication.

## Features

- **Static Page Serving:** Automatically routes requests to HTML files in the `Pages/` directory.
- **Dynamic File Listing:** Automatically generates an index of files available in a configurable directory.
- **JWT Authentication:** Secure login system using JSON Web Tokens stored in HTTP-only cookies.
- **Dual Protocol Support:** Seamlessly switches between HTTP and HTTPS.
- **JSON Configuration:** Centralized setup via `conf.json`.
- **Custom Error Handling:** Styled 404 page support.

## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.22+ for modern routing syntax)
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
  "port": 1337,
  "certFile": "cert.pem",
  "certKey": "key.pem",
  "isFileServer": true,
  "fileServerRootPath": "Files"
}
```

*   `port`: The port number for the server.
*   `certFile` / `certKey`: Paths to SSL certificate files. Leave empty for HTTP.
*   `isFileServer`: If `true`, the root URL (`/`) displays the dynamic file list.
*   `fileServerRootPath`: The directory name for served assets.

### 3. Authentication

The server includes a built-in JWT authentication system:
- **Login:** Access `/logIn` to enter credentials.
- **Security:** Upon successful login, a JWT is issued as an `HttpOnly` cookie.
- **Protected Area:** Accessing `/Auth/IamSecret` requires a valid token.
- **Logout:** Use `/logout` to terminate the session.

### 4. File Server Mode

When `isFileServer` is set to `true`:
- The root URL shows a styled index of the `fileServerRootPath` directory.
- Files are directly downloadable via the UI.
- Static pages remain accessible via their specific paths.

### 5. Running the Server

```bash
# Run all files in the current package
go run .
```

## Project Structure

- `main.go`: Primary server and routing logic.
- `auth.go`: JWT generation, validation, and login logic.
- `conf.json`: Server configuration.
- `Pages/`: Static HTML pages (e.g., `main.html`, `logIn.html`).
- `Pages/Auth/`: Protected HTML pages (e.g., `IamSecret.html`).
- `Files/`: Assets for the dynamic file server.
- `FileList.html`: Template for the file listing UI.

## Development

- **Add a Protected Page:** Place new `.html` files in `Pages/Auth/`.
- **Generate Certs:**
  ```bash
  openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes -subj "/CN=localhost"
  ```
