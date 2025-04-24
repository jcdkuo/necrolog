# 🪦 necrolog

**Log from the underworld. Summon your system’s truth, one entry at a time.**

`necrolog` is a dark-themed, structured logging utility for Go, built on top of `zap` and `lumberjack`.  
It is designed to capture logs like digital necromancy: resurrecting lost traces from the system abyss.

---

## 💡 Features

- ⚡ Blazing fast logging via [zap](https://github.com/uber-go/zap)
- 🌲 Automatic log rotation via [lumberjack](https://github.com/natefinch/lumberjack)
- 🎨 Colored terminal output + structured JSON logs
- 🔖 Auto-tagged logs with origin metadata (`origin: edotensei`, `summoned_by: necrolog`)
- 🧾 Minimal HTTP log API endpoint (`/log`)
- 🎭 Terminal banner with random necrolog quotes
- ✅ Compatible with **Go 1.18+**

---

## 🔧 Example

```go
import "github.com/yourname/necrolog"

func main() {
    necrolog.Info(necrolog.LogSystemInfo, "System initialized")
    necrolog.Warn(necrolog.LogConfigUpdate, "Missing default value")
    necrolog.Error(necrolog.LogReaderAuth, "Unauthorized access detected")
}
```

## 🌐 HTTP Usage

```bash
curl -X POST http://localhost:8080/log \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/var/log/uah_log/edotensei/system_info.log",
    "level": "info",
    "msg": "Summoning log from the web"
}'
```

## 🧙 Output Example

```json
{
  "timestamp": "2025-04-15T16:06:00Z",
  "level": "INFO",
  "origin": "edotensei",
  "summoned_by": "necrolog",
  "message": "System initialized"
}
```

## 📁 Predefined Log Paths

```go
necrolog.LogSystemInfo
necrolog.LogReaderAuth
necrolog.LogAccessLog
// ...and more (see source)

```

## 🔮 Philosophy

“Every log is a whisper from the underworld.”

We believe logging is a ritual.
You summon the truth not by print statements, but by incanting structured invocations.
necrolog is your grimoire.

## 🛠️ Coming Soon

- Prometheus metrics (Go 1.20+)
- CLI tool (necrologctl)
- Configurable output formats
- Dark terminal dashboard UI 😈

## 🧱 Requirements

- Go 1.18+
- Linux (log paths default to /var/log/uah_log/edotensei/)

## 🧞‍♂️ License

MIT. Summon freely.
