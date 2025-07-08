# ChargePi-go Issues Resolution & Steve Connection Guide

## ✅ **ISSUE 1 RESOLVED: UI Not Accessible**

### Problem
- ChargePi-go UI was not accessible at http://localhost:3000
- Connection refused errors

### Root Cause
- UI was a Svelte application that needed to be built before serving
- Static assets were missing

### Solution Applied ✅
1. **Built the UI**:
   ```bash
   cd ui
   npm install
   npm run build
   ```
2. **Restarted ChargePi-go** to serve the new static files

### Result
- **✅ UI is now working at http://localhost:3000**
- Shows proper ChargePi dashboard with HTML response
- HTTP 200 status confirmed

---

## 🚨 **ISSUE 2 DISCOVERED: Steve Connectivity Problem**

### Problem
- ChargePi-go cannot connect to Steve backend
- WebSocket connection failing

### Root Cause Discovered
- **Steve is not accessible from ChargePi-go environment**
- HTTP test to http://localhost:8180 fails
- This means websocket at ws://localhost:8180/steve/websocket/CentralSystemService will also fail

### Possible Reasons

#### 1. **Steve Not Running**
- Steve service might be stopped
- Check: `docker ps` or service status

#### 2. **Different Host/Port**
- Steve might be running on a different port
- Steve might be in a different container/host

#### 3. **Network Isolation**
- Docker networking issues
- Firewall blocking connections
- Different network namespaces

#### 4. **Steve Configuration**
- Steve WebSocket not enabled
- Different WebSocket endpoint path

---

## 🔧 **TROUBLESHOOTING STEPS**

### Step 1: Verify Steve is Running
```bash
# Check if Steve container is running (if using Docker)
docker ps | grep steve

# Check if Steve process is running
ps aux | grep steve

# Check what's listening on port 8180
netstat -tln | grep 8180
# or
ss -tln | grep 8180
```

### Step 2: Test Steve Accessibility from Your Browser
- Try accessing http://localhost:8180/steve/manager/home from your browser
- If this doesn't work, Steve has a fundamental connectivity issue

### Step 3: Identify Correct Steve URL
If Steve is running but not on localhost:8180, check:
```bash
# Look for Steve in Docker containers
docker ps --format "table {{.Names}}\t{{.Ports}}" | grep steve

# Check all listening ports
netstat -tln | grep LISTEN
```

### Step 4: Update ChargePi-go Configuration
Once you find the correct Steve URL, update `simulator-settings.yaml`:
```yaml
chargePoint:
  connectionSettings:
    uri: ws://[CORRECT_HOST]:[CORRECT_PORT]/steve/websocket/CentralSystemService
```

### Step 5: Common Steve WebSocket Paths
Try these common paths in your configuration:
- `ws://localhost:8180/steve/websocket/CentralSystemService`
- `ws://localhost:8180/steve/websocket/CentralSystem`
- `ws://localhost:8180/steve/websocket`
- `ws://localhost:8080/steve/websocket/CentralSystemService` (if on port 8080)

---

## 📋 **CURRENT STATUS**

### ✅ Working Components
- **ChargePi-go**: Running successfully (PID: 16810)
- **UI**: Accessible at http://localhost:3000
- **API**: Available on port 4269
- **EVSE Simulation**: Dummy hardware configured
- **OCPP Core Profile**: Loaded and ready

### ❌ Not Working
- **Steve Connection**: Cannot reach Steve backend
- **WebSocket**: Cannot establish OCPP connection

### 🎯 **Immediate Next Steps**

1. **Check Steve Status**: Verify Steve is actually running and accessible
2. **Find Correct Steve URL**: Determine the right host:port for Steve
3. **Update Configuration**: Fix the Steve URI in simulator-settings.yaml
4. **Restart ChargePi-go**: Apply the corrected configuration
5. **Verify Connection**: Check Steve dashboard for ChargePi-Simulator

---

## 🛠️ **Quick Test Commands**

```bash
# Test Steve accessibility
curl -I http://localhost:8180/steve/manager/home

# Test different ports if 8180 doesn't work
curl -I http://localhost:8080/steve/manager/home
curl -I http://localhost:9000/steve/manager/home

# Check what's running on common ports
for port in 8080 8180 9000 3000; do 
  echo "Testing port $port:"
  curl -s --max-time 2 http://localhost:$port/ | head -2
done
```

---

## 🎉 **Once Steve is Fixed**

When Steve connectivity is restored, you should see:
1. `ChargePi-Simulator` appears in Steve's charge points list
2. Status shows as "Available" or "Connected"
3. OCPP messages flow between ChargePi-go and Steve
4. You can send commands from Steve to the simulator

The ChargePi-go simulator is **100% ready** - it just needs to reach Steve!