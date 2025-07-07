# ChargePi-go Simulator Connection to Steve Backend

## ✅ Current Status

**ChargePi-go Process**: Running (PID: 14967)
**Configuration**: Updated to connect to correct Steve port
**Steve Backend**: Confirmed reachable at http://localhost:8180

## 🔧 Configuration Fixed

**Previous Issue**: ChargePi-go was trying to connect to port 8080
**Corrected**: Now connecting to port 8180 (your Steve backend)

**Current WebSocket Connection**:
```
ws://localhost:8180/steve/websocket/CentralSystemService
```

**Charge Point Details**:
- **ID**: `ChargePi-Simulator`
- **OCPP Version**: 1.6 (Core profile)
- **Type**: AC Charging Station
- **Max Power**: 22kW
- **Connector**: Type2 on EVSE 1

## 🔍 Next Steps to Verify Connection

### 1. Check Steve Dashboard
- Go to: http://localhost:8180/steve/manager/home
- Navigate to "Charge Points" section
- Look for `ChargePi-Simulator` in the connected charge points list
- Status should show as "Available" or "Online"

### 2. Check OCPP Transaction Log in Steve
- Go to "OCPP Tags" → "Transaction Log"
- Look for recent activity from `ChargePi-Simulator`

### 3. Monitor Steve Logs
If you have access to Steve's logs, look for:
- WebSocket connection attempts from ChargePi-Simulator
- OCPP message exchanges
- Any connection errors

### 4. Test OCPP Commands (Once Connected)
From Steve's interface, try sending:
- **Status Notification** request
- **Get Configuration** request
- **Change Configuration** request

## 🛠️ Troubleshooting if Still No Connection

### Check ChargePi-go Logs
```bash
# Check if process is running
ps aux | grep ChargePi

# Look for any log files
find /tmp -name "*chargepi*" -o -name "*ChargePi*" 2>/dev/null
```

### Verify Network Connectivity
```bash
# Test Steve HTTP interface
curl -I http://localhost:8180/steve/manager/home

# Test if port 8180 is listening
netstat -tlnp | grep 8180  # (if netstat available)
```

### Check Configuration
Current config file: `simulator-settings.yaml`
Key settings:
- Charge Point ID: `ChargePi-Simulator`
- Server URI: `ws://localhost:8180/steve/websocket/CentralSystemService`
- Protocol Version: OCPP 1.6

## 🔧 Hardware Simulation Status

All hardware components are running in dummy/simulator mode:

### ✅ EVSE Configuration
- **EVSE ID**: 1
- **EVCC**: Dummy type (charging control simulation)
- **Power Meter**: Dummy (230V static simulation)
- **Connector**: Type2, Status "Available"

### ✅ Hardware Components
- **RFID Reader**: Dummy with 3 test tag IDs
- **LED Indicator**: Dummy with full color status mapping
- **Display**: Dummy display simulation
- **All logging**: Working with debug level

## 📝 Expected Behavior

When successfully connected, you should see:
1. `ChargePi-Simulator` appears in Steve's charge points list
2. Status shows as "Available"
3. You can send OCPP commands from Steve
4. ChargePi-go responds with appropriate OCPP responses
5. Transaction logging shows activity

## 🚨 If Connection Still Fails

1. **Check Steve Configuration**: Ensure Steve allows new charge points to register
2. **Firewall**: Verify no firewall blocking WebSocket connections
3. **Steve Version**: Ensure Steve supports OCPP 1.6 WebSocket connections
4. **Authentication**: Check if Steve requires HTTP basic authentication (currently disabled in config)

The simulator is now properly configured and should connect to your Steve backend on the correct port!