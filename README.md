# Proxier

**Proxier** is a high-performance **reverse proxy** that routes external requests to predefined origin servers.
It is lightweight, configurable via `config.yaml`, and supports multiple endpoints.

## 🚀 Features
- 🌍 **Dynamic Proxy Routing** – Easily define multiple proxy rules in `config.yaml`
- 🛠 **Simple Configuration** – No database required, just YAML-based settings
- 🚀 **Fast & Efficient** – Optimized request forwarding
- 🏗 **Cross-Platform** – Works on Linux, macOS, and Windows

---

## 📦 Installation
### **Using Release**
Download latest version from [Release](https://github.com/ezex-io/proxier/releases).

### **Using Go**
```sh
go install github.com/ezex-io/proxier@latest
```

### **Using Docker**
```sh
docker pull ezexio/proxier:latest
docker run -p 8080:8080 -e EZEX_PROXIER_ADDRESS=127.0.0.1:8081 -e EZEX_PROXIER_PROXY_RULES=/foo|https://httpbin.org/get,/bar|https://google.com ezexio/proxier
```

### **Build from Source**
```sh
git clone https://github.com/ezex-io/proxier.git
cd proxier
go build -o proxier ./cmd/proxier/main.go
```

---

Here’s an improved and clarified version of your README section to reflect the current configuration via **environment variables**, including formatting, examples, and explanation:

---

## ⚙️ Configuration

The proxy server can be configured using **environment variables**. Here's a breakdown of the supported variables:

### ✅ Environment Variables

| Variable Name                  | Description                                                                    | Example           |
| ------------------------------ | ------------------------------------------------------------------------------ | ----------------- |
| `PROXIER_ADDRESS`         | Host and port for the proxy server to bind to                                  | `127.0.0.1:8081`  |
| `PROXIER_ENABLE_FASTHTTP` | Enable [`fasthttp`](https://github.com/valyala/fasthttp) instead of `net/http` | `true` or `false` |
| `PROXIER_RULES`     | Comma-separated list of proxy rules in \`key                                   | `[{"endpoint":"/foo","destination":"https://httpbin.org/get"}, {"endpoint":"/bar","destination":"https://google.com"}]` |

* **key**: The path endpoint to intercept (e.g., `/foo`)
* **val**: The destination URL to which the request should be proxied

### 🔁 Example

```env
PROXIER_ADDRESS=127.0.0.1:8081
PROXIER_ENABLE_FASTHTTP=false
PROXIER_RULES=[{"endpoint":"/foo","destination":"https://httpbin.org/get"}, {"endpoint":"/bar","destination":"https://google.com"}]
```

This setup creates the following proxy routes:

| Endpoint | Proxies To                |
| -------- | ------------------------- |
| `/foo`   | `https://httpbin.org/get` |
| `/bar`   | `https://google.com`      |

---


---

## 🚀 Running the Server
### **Start Proxier**
```sh
./proxier
```

### **Check if Proxier is Running**
```sh
curl -i http://localhost:8080/
```
**Response:**
```
HTTP/1.1 200 OK
Proxier is running
```

### **Health Check API**
```sh
curl -i http://localhost:8080/livez
```
**Response:**
```
HTTP/1.1 200 OK
OK
```

### **Proxy Requests**
Example request to `dex` proxy:
```sh
curl -i http://localhost:8080/dex/
```

---

## 🛠 Development & Contribution
### **Setup Development Environment**
```sh
git clone https://github.com/ezex-io/proxier.git
cd proxier
go mod tidy
```

### **Run Tests**
```sh
make test
```

### **Build Proxier**
```sh
make build_linux
```

### **Code Formatting & Linting**
```sh
make check
```

---

## 📜 License
Proxier is licensed under the **MIT License**. See [LICENSE](./LICENSE) for details.

---

## 💡 Contributing
We welcome contributions! Please follow these steps:
1. **Fork the repository**
2. **Create a new branch** (`feature/my-feature`)
3. **Commit your changes** (`git commit -m "Add new feature"`)
4. **Push your branch** (`git push origin feature/my-feature`)
5. **Create a Pull Request**
