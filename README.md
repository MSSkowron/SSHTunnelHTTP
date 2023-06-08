# ssh-tunnel-http

This code provides a basic implementation of SSH tunneling using the github.com/gliderlabs/ssh package in Go. It allows establishing an SSH server that can handle incoming SSH sessions and forward their data to HTTP clients via an HTTP server. It can be used as a starting point for building more complex SSH tunneling applications.

## Prerequisites

Before running this project, ensure that you have Go installed on your system.

## Install

Follow these steps to install the project:

1. Clone the repository:

   ```
   git clone https://github.com/MSSkowron/ssh-tunnel-http.git
   ```

2. Install _gliderlabs_ SSH package:

   ```
   go get github.com/gliderlabs/ssh
   ```

3. Change to the project directory:

   ```
   cd ssh-tunnel-http
   ```

## Run

To run the project, execute the following command in the project directory:

```
go run main.go
```

The SSH server will be listening on port 2222, and the HTTP server will be listening on port 3000.

## Usage

- **SSH Tunneling**

  Once the server is running, you can establish an SSH session to transfer data through a tunnel. Use your preferred SSH client to connect to the server:

  Example (using OpenSSH client):

  ```
  ssh user@localhost -p 2222 < file.txt
  ```

  Data sent from the SSH client will be forwarded to the corresponding tunnel.

- **API Endpoint**

  The API endpoint / allows HTTP clients to access the data being forwarded through the SSH tunnel. The endpoint expects a query parameter id that represents the tunnel ID. The ID is printed in the server logs when an SSH connection is established. The data from the corresponding tunnel will be returned as the response.

  Example:

  ```
  GET /?id=12345
  ```

  This will retrieve the data from tunnel ID 12345 and return it as the response.
