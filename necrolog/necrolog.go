package necrolog

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "runtime"
    "strings"
    "sync"
    "time"

    "github.com/joho/godotenv"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    buffer "go.uber.org/zap/buffer"
    "gopkg.in/natefinch/lumberjack.v2"
)

var (
    loggers         = make(map[string]*zap.Logger)
    loggersMu       sync.RWMutex
    logTimezoneStr  string
)

func init() {
    _ = godotenv.Load()
    printBanner()
}

func getZapLevel() zapcore.Level {
    if val := strings.ToLower(os.Getenv("ZAP_LOG_LEVEL")); val != "" {
        switch val {
        case "debug":
            return zapcore.DebugLevel
        case "info":
            return zapcore.InfoLevel
        case "warn":
            return zapcore.WarnLevel
        case "error":
            return zapcore.ErrorLevel
        }
    }
    return zapcore.DebugLevel
}

func (c *colorConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
    levelStr := strings.ToUpper(entry.Level.String())
    colorMsg := formatConsoleMessage(levelStr, entry.Message)

    newEntry := entry
    newEntry.Message = colorMsg

    // Call the underlying Encoder's EncodeEntry method
    // Returns a buffer containing the encoded log entry and potential error
    // buffer.Buffer is zapcore's internal buffer type
    // Used to store encoded log data
    return c.Encoder.EncodeEntry(newEntry, fields)
}

func GetLogger(path string) *zap.Logger {

    logTimezoneStr = os.Getenv("ZAP_LOG_TIMEZONE")

    loggersMu.RLock()
    logger, exists := loggers[path]
    loggersMu.RUnlock()
    if exists {
        return logger
    }

    loggersMu.Lock()
    defer loggersMu.Unlock()
    if logger, exists := loggers[path]; exists {
        return logger
    }

    _ = os.MkdirAll(LogDeadRisingPath, 0755)

    lj := &lumberjack.Logger{
        Filename:   path,
        MaxSize:    1,
        MaxBackups: 1,
        Compress:   true,
    }

    fileEncoderCfg := zapcore.EncoderConfig{
        TimeKey:    "time",
        LevelKey:   "level",
        MessageKey: "msg",
        EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000" + " " + logTimezoneStr),
        EncodeLevel: zapcore.CapitalLevelEncoder,
    }
    jsonEncoder := zapcore.NewJSONEncoder(fileEncoderCfg)

    consoleEncoderCfg := zapcore.EncoderConfig{
        MessageKey: "msg",
        LevelKey:   "",
        TimeKey:    "",
        EncodeTime: nil,
        EncodeLevel: nil,
    }
    baseConsoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderCfg)
    colorConsoleEncoder := &colorConsoleEncoder{Encoder: baseConsoleEncoder}

    level := getZapLevel()
    core := zapcore.NewTee(
        zapcore.NewCore(jsonEncoder, zapcore.AddSync(lj), level),
        zapcore.NewCore(colorConsoleEncoder, zapcore.AddSync(os.Stdout), level),
    )

    logger = zap.New(core)
    loggers[path] = logger
    return logger
}

func LogMessage(path string, level LogLevel, msg string) {
    logger := GetLogger(path)

    emit(logger, level, msg)
    _ = logger.Sync()
}

func emit(logger *zap.Logger, level LogLevel, msg string) {
    switch level {
    case LevelInfo:
        logger.Info(msg)
    case LevelWarn:
        logger.Warn(msg)
    case LevelError:
        logger.Error(msg)
    default:     
        logger.Debug(msg)
    }
}

func Debug(path, msg string) { LogMessage(path, LevelDebug, msg) }
func Info(path, msg string)  { LogMessage(path, LevelInfo, msg)  }
func Warn(path, msg string)  { LogMessage(path, LevelWarn, msg)  }
func Error(path, msg string) { LogMessage(path, LevelError, msg) }

func formatConsoleMessage(level string, msg string) string {
    timestamp := time.Now().Format("2006-01-02 15:04:05.000")
    var color string
    switch level {
    case "DEBUG":
        color = "[32mDEBUG[0m"
    case "INFO":
        color = "[34mINFO[0m"
    case "WARN":
        color = "[33mWARN[0m"
    case "ERROR":
        color = "[31mERROR[0m"
    default:
        color = level
    }
    return fmt.Sprintf("%s [%s] %s %s", color, timestamp, logTimezoneStr, msg)
}

func getCallerInfo() string {
    pc := make([]uintptr, 10)
    n := runtime.Callers(3, pc)
    frames := runtime.CallersFrames(pc[:n])
    for {
        frame, more := frames.Next()
        if !strings.Contains(frame.Function, "necrolog.") {
            return fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
        }
        if !more {
            break
        }
    }
    return "unknown"
}

func LogHandler(w http.ResponseWriter, r *http.Request) {
    type reqBody struct {
        Path  string   `json:"path"`
        Level LogLevel `json:"level"`
        Msg   string   `json:"msg"`
    }

    var req reqBody
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid json", http.StatusBadRequest)
        return
    }
    if req.Path == "" || req.Msg == "" {
        http.Error(w, "missing path or msg", http.StatusBadRequest)
        return
    }
    LogMessage(req.Path, req.Level, req.Msg)
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("ok"))
}
