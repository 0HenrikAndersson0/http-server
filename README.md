# Go API Lab

<img width="1024" height="559" alt="image" src="https://github.com/user-attachments/assets/452f460a-24a8-4363-a6ad-1c96892b0b2f" />

A simple, lightweight Go web server designed to serve static HTML pages with support for both HTTP and HTTPS, and a dynamic file-sharing mode.

## Features

- **Static Page Serving:** Automatically routes requests to HTML files in the `Pages/` directory.
- **Dynamic File Listing:** Automatically generates an index of files available in the `Files/` directory.
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
  "certKey": "key.pem",
  "isFileServer": false,
  "fileServerRootPath": "Files"
}
```

*   `port`: The port number for the server to listen on.
*   `certFile` / `certKey`: Paths to SSL certificate files for HTTPS. Leave empty for HTTP.
*   `isFileServer`: If `true`, the root URL (`/`) will display a dynamic list of files.
*   `fileServerRootPath`: The name of the directory containing the files to serve (defaults to "Files").

### 3. File Server Mode

When `isFileServer` is set to `true`:
- Navigating to `/` will show a styled list of all assets inside the directory specified by `fileServerRootPath`.
- Files can be downloaded directly by clicking their names.
- Static pages in `Pages/` remain accessible via their specific paths (e.g., `/main` or `/test`).

### 4. Generating SSL Certificates (Optional)

To test HTTPS locally, generate a self-signed certificate:

```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes -subj "/CN=localhost"
```

### 5. Running the Server

```bash
go run main.go
```

## Project Structure

- `main.go`: The core server logic and request handling.
- `conf.json`: Configuration settings for the server.
- `Pages/`: Directory containing static `.html` files.
- `Files/`: Directory containing assets available for download.
- `FileList.html`: The HTML template used for the dynamic file listing.
- `.gitignore`: Standard Go and environment exclusions.

## Development

- **Add a Page:** Create a new `.html` file in `Pages/` (e.g., `Pages/about.html` -> `/about`).
- **Add an Asset:** Drop any file into `Files/` to make it appear in the File Server index.
