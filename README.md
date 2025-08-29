# PEM Certificate Expiration Checker

## Overview and Objectives

This command-line tool, written in Go, checks the expiration status of X.509 certificates in PEM format. The primary objective is to provide a simple, efficient, and reliable way to monitor certificate expiration. By scanning either a single file or an entire directory, the tool helps prevent service disruptions caused by expired certificates.

checkcert is a certificate expiration check utility. It will recurse all sub directories and alert if any certificate is exceeding the renewal threshold specified.

The application is designed with performance in mind, using concurrency to check multiple certificates in parallel, which makes it fast and efficient for large numbers of certificates.

## Command-Line Syntax

The application is configured using the following command-line flags:

| Flag      | Description                                                | Default |
|-----------|------------------------------------------------------------|---------|
| `-file`   | Path to a single PEM certificate file.                     |         |
| `-dir`    | Path to a directory containing PEM certificate files.      |         |
| `-ext`    | File extension to search for in the directory.             | `.pem`  |
| `-days`   | The number of days to warn before a certificate expires.   | `30`    |
| `-header` | Display a header in the output report.                     | `false` |
| `-version`| Display the application's version number.                  | `false` |

## Usage Examples

### 1. Check a single certificate file:
```bash
./checkcert -file /path/to/your/certificate.pem
```

### 2. Check all certificates in a directory:
This will scan the `/path/to/certs` directory for all files with the `.pem` extension.
```bash
./checkcert -dir /path/to/certs
```

### 3. Check certificates with a custom warning period and display the report header:
This command will check all certificates in the `/path/to/certs` directory and flag any that are expiring within the next 60 days. It will also include a header in the output for better readability.
```bash
./checkcert -dir /path/to/certs -days 60 -header
```

### 4. Check certificates with a custom extension
This will check all certs  in the current directory, with a .crt extension, that will expire in less than 60 days
```checkcert -dir . -ext crt -days 60```
