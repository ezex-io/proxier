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
docker run -p 8080:8080 -v $(pwd)/config.yaml:/etc/proxier/config.yaml ezexio/proxier
```

### **Build from Source**
```sh
git clone https://github.com/ezex-io/proxier.git
cd proxier
go build -o proxier ./cmd/proxier/main.go
```

---

## ⚙️ Configuration (`config.yaml`)
Define your proxy routes in a **YAML config file**:
```yaml
server:
  host: "0.0.0.0"
  listen_port: "8080"
  fast_http: true

proxy:
  - endpoint: /foo1
    destination_url: "https://example.com/bar1"

  - endpoint: /foo2
    destination_url: "https://example.com/bar2"

  - endpoint: /foo3
    destination_url: "https://example.com/bar3"
```

---

## 🚀 Running the Server
### **Start Proxier**
```sh
./proxier -config ./config.yaml
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
