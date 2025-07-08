# ChargePi-go Docker Configuration Solution

## Problem Discovered
The Docker container validation was failing with:
```
{"error":"Key: 'Settings.ChargePoint.ConnectionSettings.ServerUri' Error:Field validation for 'ServerUri' failed on the 'required' tag","level":"fatal","msg":"Invalid settings","time":"2025-07-08T04:26:10Z"}
```

## Root Cause Analysis
1. **Direct execution works**: Our configuration works perfectly when running `go run . run --settings=test-config.yaml`
2. **Double ws:// prefix issue**: The application automatically adds `ws://` prefix, so the URI should NOT include it
3. **Field name confusion**: The codebase has multiple ConnectionSettings structs:
   - Internal settings struct: uses `uri`, `id`, `basicAuthUser`, etc.
   - Protobuf struct: uses `url`, `chargePointId`, `basicAuthUsername`, etc.

## Solutions to Try

### Option 1: Protobuf Field Names (test-config.yaml)
```yaml
chargePoint:
  connectionSettings:
    chargePointId: ChargePi-Simulator
    protocolVersion: '1.6'
    url: host.docker.internal:8180/steve/websocket/CentralSystemService
    basicAuthUsername: ''
    basicAuthPassword: ''
```

### Option 2: Internal Struct Field Names (test-config-original-fields.yaml)
```yaml
chargePoint:
  connectionSettings:
    id: ChargePi-Simulator
    protocolVersion: '1.6'
    uri: host.docker.internal:8180/steve/websocket/CentralSystemService
    basicAuthUser: ''
    basicAuthPass: ''
```

## Docker Commands to Test

### Test Option 1
```bash
docker stop chargepi-test && docker rm chargepi-test
docker run -d --name chargepi-test \
  -p 3000:3000 \
  -p 4269:4269 \
  -v $(pwd)/test-config.yaml:/etc/ChargePi/configs/settings.yaml \
  docker-chargepi \
  run
docker logs chargepi-test
```

### Test Option 2
```bash
docker stop chargepi-test && docker rm chargepi-test
docker run -d --name chargepi-test \
  -p 3000:3000 \
  -p 4269:4269 \
  -v $(pwd)/test-config-original-fields.yaml:/etc/ChargePi/configs/settings.yaml \
  docker-chargepi \
  run
docker logs chargepi-test
```

### Debug Command (if still failing)
```bash
docker run --rm -it \
  -v $(pwd)/test-config.yaml:/etc/ChargePi/configs/settings.yaml \
  docker-chargepi \
  sh -c "echo '=== Config file content ===' && cat /etc/ChargePi/configs/settings.yaml && echo '=== Running application ===' && /app/main run"
```

## Expected Success Indicators
✅ Database initialization logs
✅ No ServerUri validation error
✅ "Trying to connect to the central system: ws://host.docker.internal:8180/steve/websocket/CentralSystemService"
✅ API and UI servers starting
✅ All dummy hardware initializing

## Key Fixes Applied
1. **Removed ws:// prefix** from URI (application adds it automatically)
2. **Created two field name variants** to test both struct formats
3. **Proper Docker networking** using `host.docker.internal:8180`
4. **Correct file mounting** to `/etc/ChargePi/configs/settings.yaml`