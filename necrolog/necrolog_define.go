package necrolog

import (
	"fmt"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	LogDeadRisingPath   = "/var/log/japp/edotensei/"
	LogLinkLayer        = LogDeadRisingPath + "link_layer.log"
	LogPowerEvent       = LogDeadRisingPath + "power_event.log"
	LogUnlockSchedule   = LogDeadRisingPath + "unlock_schedule.log"
)

const (
	LogDefaultPath 		= "/var/log/japp/"
	LogDefaultFile 		= LogDefaultPath + "app.log"
	LogErrorFile 		= LogDefaultPath + "error.log"
	LogWarningFile 		= LogDefaultPath + "warning.log"
)

// 1.2MB in bytes (explicitly using integer math to avoid truncation)
const MaxLogFileSize int64 = 1228800 // 1.2 * 1024 * 1024

// WhitelistedFiles contains paths that should not be deleted unless they exceed size limits
var WhitelistedFiles = map[string]bool{
	LogLinkLayer:      true,
	LogPowerEvent:     true,
	LogUnlockSchedule: true,
}

type LogLevel string
const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type colorConsoleEncoder struct {
	zapcore.Encoder	
}

func printBanner() {
	phrases := []string{
		"Log what others fear to remember.",
		"Every log is a whisper from the underworld.",
		"Truth resurrected through error.",
		"If you forget the failure, it returns stronger.",
		"Logs never die. They haunt.",
		"Speak in logs, summon the truth.",
		"Even if systems crash, necrolog remains.",
		"Born from ashes of segfaults and stack traces.",
		"Log like you're writing your system's obituary.",
		"What you record lives beyond you.",
	}
	seed := time.Now().UnixNano() % int64(len(phrases))
	fmt.Printf(`
  _   _                 _               
 | \ | | ___  ___ _ __| | ___  ___ ___ 
 |  \| |/ _ \/ __| '__| |/ _ \/ __/ __|
 | |\  |  __/ (__| |  | |  __/\__ \__ \
 |_| \_|\___|\___|_|  |_|\___||___/___/

      :: necrolog :: log from the underworld

  "%s"

`, phrases[seed])
}
