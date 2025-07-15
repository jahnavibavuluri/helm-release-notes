# Setup & Usage Instructions 

## Prerequisites

- **VPN Access:**
You must be connected to the MathWorks VPN to access the GitHub Enterprise server.

- **Go:**
Make sure you have [Go installed](https://go.dev/doc/install) on your system.

- **PAT:**
You will need a PAT (personal access token) to let the code access and download code from the GH Enterprise server.

---

## 1. Clone the Repository

```sh
git clone https://github.com/jahnavibavuluri/helm-release-notes.git
cd helm-release-notes
```

## 2. Create a Personal Access Token (PAT)

## 3. Set the GITHUB_TOKEN environment variable

For Linux, macOS, WSL:
```sh
export GITHUB_TOKEN=<your-pat-token-here>
```

For Windows Command Prompt:

```cmd
set GITHUB_TOKEN=<your-pat-token-here>
```

## 4. Run the code!

The usage for this code is: go run main.go <owner> <repo> <old_ref> <new_ref>

An example is: 
```go
go run main.go development terraform-aws-<module> v1.0.0 v1.1.0
```