# 🚀 Quick Fix Commands for ChargePi Docker

## 🎯 **ISSUE IDENTIFIED**
Container failing because it expects `serverUri` field but config uses `uri`.

## ⚡ **IMMEDIATE FIX** 

Run these commands on your **host machine** (where you ran `docker ps`):

### **Step 1: Stop the failing container**
```bash
docker stop docker-chargepi-1
docker rm docker-chargepi-1
```

### **Step 2: Create the correct configuration file**
**Copy and save this as `docker-simulator-settings.yaml`:**

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
    serverUri: ws://host.docker.internal:8180/steve/websocket/CentralSystemService
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
      firmwareVersion: 1.0.0
      serialNumber: SIM001

  logging:
    type:
      - file
    format: text
    level: debug

  hardware:
    display:
      enabled: true
      driver: dummy
      language: en
      dummy:
        enabled: true
    reader:
      enabled: true
      model: dummy
      dummy:
        tagIds:
          - "04142C5A2D4980"
          - "04F4275A2D4980"
          - "040E2D5A2D4980"
    indicator:
      enabled: true
      type: dummy
      indicateCardRead: true
      statuses:
        available: "#00FF00"
        reserved: "#0000FF"
        preparing: "#FFFF00"
        charging: "#FF8000"
        finishing: "#800080"
        fault: "#FF0000"
        error: "#FF0000"
      dummy:
        enabled: true
```

### **Step 3: Run the fixed container**
```bash
docker run -d --name chargepi-simulator \
  -p 3000:3000 \
  -p 4269:4269 \
  -v $(pwd)/docker-simulator-settings.yaml:/app/settings.yaml \
  docker-chargepi \
  run --settings=/app/settings.yaml
```

### **Step 4: Verify it's working**
```bash
# Check container status
docker ps

# Check logs (should show no errors)
docker logs chargepi-simulator

# Test UI access
curl http://localhost:3000/
```

## ✅ **SUCCESS INDICATORS**

- ✅ `docker ps` shows container "Up" (not "Restarting")
- ✅ `http://localhost:3000` works in your browser
- ✅ `ChargePi-Simulator` appears in Steve at `http://localhost:8180/steve/manager/chargepoints`

## 🎯 **KEY CHANGES MADE**

1. **Fixed field name**: `uri` → `serverUri`
2. **Added port forwarding**: `-p 3000:3000 -p 4269:4269`
3. **Fixed Steve URL**: `host.docker.internal:8180` (to reach Steve from container)
4. **Proper file mounting**: Volume mounted correctly

That's it! Your ChargePi simulator should now work perfectly! 🚀