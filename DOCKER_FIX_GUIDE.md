# 🐳 ChargePi Docker Fix Guide

## 🚨 **CURRENT PROBLEM**
Your ChargePi container is **failing and restarting continuously**, plus it has **no port forwarding** configured.

```
docker-chargepi-1    Restarting (1) 10 seconds ago    [NO PORTS MAPPED]
```

## 🔧 **STEP-BY-STEP FIX**

### **Step 1: Check the Error Logs**
Run this on your **host machine** (outside any container):
```bash
docker logs docker-chargepi-1
```
This will show you why the container keeps failing.

### **Step 2: Stop and Remove the Failing Container**
```bash
docker stop docker-chargepi-1
docker rm docker-chargepi-1
```

### **Step 3: Run ChargePi with Proper Port Forwarding**

#### Option A: Using Docker Run (Recommended)
```bash
docker run -d --name chargepi-simulator \
  -p 3000:3000 \
  -p 4269:4269 \
  -v $(pwd)/simulator-settings.yaml:/app/simulator-settings.yaml \
  docker-chargepi \
  run --settings=/app/simulator-settings.yaml
```

#### Option B: If You Have Docker Compose
Update your `docker-compose.yml` to include port mappings:
```yaml
services:
  chargepi:
    image: docker-chargepi
    ports:
      - "3000:3000"    # UI
      - "4269:4269"    # API
    volumes:
      - ./simulator-settings.yaml:/app/simulator-settings.yaml
    command: run --settings=/app/simulator-settings.yaml
```

Then run:
```bash
docker-compose up -d chargepi
```

## 🎯 **CONFIGURATION FOR STEVE**

Since your Steve is at `localhost:8180`, update your ChargePi config:

### **Create/Update simulator-settings.yaml**
```yaml
api:
  enabled: true
  address: 0.0.0.0:4269

ui:
  enabled: true
  address: 0.0.0.0:3000

chargePoint:
  connectionSettings:
    id: ChargePi-Simulator
    protocolVersion: '1.6'
    uri: ws://host.docker.internal:8180/steve/websocket/CentralSystemService
    # Note: Use host.docker.internal instead of localhost to reach Steve from container
    basicAuthUser: ''
    basicAuthPass: ''
    tls:
      isEnabled: false

  info:
    type: AC
    maxPower: 22
    maxChargingTime: 180
    freeMode: false
    ocpp:
      vendor: ChargePi
      model: ChargePi-Simulator

  hardware:
    display:
      enabled: true
      driver: dummy
      dummy: {}
    reader:
      enabled: true
      model: dummy
      dummy:
        tagIds:
          - "04142C5A2D4980"
          - "04F4275A2D4980"
    indicator:
      enabled: true
      type: dummy
      dummy: {}
```

## 📋 **COMMON CONTAINER FAILURE REASONS**

### **1. Missing Settings File**
Error: `settings file not found`
**Fix**: Mount the settings file with `-v` option

### **2. Permission Issues**
Error: `permission denied`
**Fix**: Check file permissions or run with `--user` flag

### **3. Port Already in Use**
Error: `port already allocated`
**Fix**: Use different ports or stop conflicting services

### **4. Missing Dependencies**
Error: `command not found` or library errors
**Fix**: Rebuild the Docker image with all dependencies

## ✅ **VERIFICATION STEPS**

### **1. Check Container is Running**
```bash
docker ps
```
Should show:
```
chargepi-simulator   Up X minutes   0.0.0.0:3000->3000/tcp, 0.0.0.0:4269->4269/tcp
```

### **2. Test UI Access**
```bash
curl http://localhost:3000/
```

### **3. Check Container Logs**
```bash
docker logs chargepi-simulator
```
Should show startup messages, not error loops.

### **4. Test Steve Connection**
Check your Steve dashboard at:
```
http://localhost:8180/steve/manager/chargepoints
```
Look for `ChargePi-Simulator` in the list.

## 🚨 **TROUBLESHOOTING**

### **If Container Still Fails**
1. **Check the Docker image**: `docker images | grep chargepi`
2. **Rebuild if needed**: Check your Dockerfile
3. **Run interactively**: `docker run -it docker-chargepi sh`
4. **Check logs**: `docker logs [container-name]`

### **If Steve Connection Fails**
Try different Steve URL formats:
- `ws://host.docker.internal:8180/steve/websocket/CentralSystemService`
- `ws://172.17.0.1:8180/steve/websocket/CentralSystemService`
- `ws://host.docker.internal:8080/steve/websocket/CentralSystemService`

## 🎉 **SUCCESS INDICATORS**

When everything works, you'll see:
- ✅ Container shows "Up" status (not "Restarting")
- ✅ UI accessible at `http://localhost:3000`
- ✅ `ChargePi-Simulator` appears in Steve dashboard
- ✅ Container logs show successful startup, not errors

---

**First, run the `docker logs` command to see what's causing the container to fail, then follow the appropriate fix steps above!** 🚀