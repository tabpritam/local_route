# RouteLocal Documentation

🌐 **[Download Now](https://tabpritam.github.io/local_route/)**

<a href="https://buymeacoffee.com/tabpritam" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 145px !important;" ></a>

RouteLocal is a zero-configuration, cross-platform CLI tool that makes it incredibly simple to share your local development server across your network or the public internet.

## Can Layman Users Run This? (No Setup Required!)
**YES!** Because RouteLocal is built in Go, it compiles down into a single **standalone native executable file**. 

This is the biggest advantage of our architecture: your users **DO NOT** need to install Go, Node.js, Python, or configure any system PATH environment variables. They don't even need to install Cloudflare manually, because our tool downloads it for them in the background! 

You simply send them the `routelocal.exe` file, they double-click or run it in their terminal, and it works perfectly out of the box.

---

## How It Works Behind the Scenes
When you run RouteLocal, it performs the following steps instantly:
1. **Application Detection**: It actively pings your specified local port (e.g., `3000`) to ensure your app is actually running before proceeding.
2. **Reverse Proxying**: It spins up a high-performance HTTP & WebSocket reverse proxy engine that dynamically binds to your network interfaces. This securely forwards traffic from outside devices directly to your local `127.0.0.1` server, fully supporting modern features like Hot Module Replacement (HMR).
3. **mDNS Broadcasting**: It broadcasts a custom `.local` domain to your router, allowing supported devices (like iPhones and MacBooks) to connect via name rather than IP.
4. **Cloudflare Tunneling (Optional)**: If you use the `--public` flag, it silently downloads the official Cloudflare binary (if missing) to a hidden cache and spins up a secure tunnel to the public internet, extracting and providing you with the final `trycloudflare.com` URL.

---

## Install on macOS & Linux

The easiest way to install RouteLocal is using our automated installation script. It will detect your OS (macOS/Linux) and architecture (Apple Silicon/Intel) and securely install the latest version.

```bash
curl -fsSL https://raw.githubusercontent.com/tabpritam/local_route/main/install.sh | bash
```

### Check version

```bash
routelocal --version
```

### Update

To update to the latest release, simply run the installation command again. It will safely replace your existing installation.

### Uninstall

RouteLocal includes a clean uninstaller:

```bash
curl -fsSL https://raw.githubusercontent.com/tabpritam/local_route/main/uninstall.sh | bash
```
*(Or you can use the built-in `routelocal uninstall` command).*

---

## Manual Installation (Windows & Alternate macOS/Linux)

You can download the standalone executable manually if you prefer not to use the installation script.

### Windows

1. Download the latest `routelocal-windows-amd64.exe` from the [Releases page](https://github.com/tabpritam/local_route/releases/latest).
2. Open **PowerShell** or **Command Prompt** in the folder where you saved it.
3. Run the following command:
   ```powershell
   .\routelocal-windows-amd64.exe --port 3000 --name myapp --public
   ```
*(Note: If a Windows Defender Firewall popup appears the first time you run it, click "Allow access" to ensure devices on your Wi-Fi are permitted to connect).*

### macOS / Linux

Download the correct binary for your system from the [Releases page](https://github.com/tabpritam/local_route/releases/latest), then grant it execute permissions:

```bash
chmod +x ./routelocal-darwin-arm64
./routelocal-darwin-arm64 --port 3000 --name myapp --public
```

---

## Compiling Release Binaries (For Developers)

Because RouteLocal is written in Go, you can easily cross-compile the standalone binaries for every major operating system directly from your machine.

Open PowerShell and run the following commands to generate highly optimized, compressed binaries:

**1. Build for Windows (.exe):**
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o routelocal-windows.exe ./cmd/routelocal
```

**2. Build for Mac (Apple Silicon / M1 / M2):**
```powershell
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -ldflags="-s -w" -o routelocal-mac-arm64 ./cmd/routelocal
```

**3. Build for Mac (Intel):**
```powershell
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o routelocal-mac-intel ./cmd/routelocal
```

**4. Build for Linux:**
```powershell
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -o routelocal-linux ./cmd/routelocal
```
