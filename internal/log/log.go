package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Setup creates needed directories, log file(s), sets the level
// and removes old log files if needed.
func Setup() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	debug := strings.ToLower(os.Getenv("MINYLS_DEVELOPMENT")) == "true"

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		newLogger := log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		logger = &newLogger
		logger.Warn().Msg("debug logging enabled")
		return nil
	}

	err := createLogsDirIfNotExists()
	if err != nil {
		return fmt.Errorf("createLogsDirIfNotExists: %w", err)
	}

	go removeOldLogs()

	w, err := createLogFile()
	if err != nil {
		return fmt.Errorf("createLogFile: %w", err)
	}

	newLogger := zerolog.New(w)
	logger = &newLogger

	return nil
}

// Log serves as a generic function to get the logger
// we created when running Setup() and use the sub functions
// of the *zerolog.Logger.
func Log() *zerolog.Logger {
	return logger
}

var logger *zerolog.Logger

var logsDir = ""

func createLogsDirIfNotExists() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable: %w", err)
	}

	logsDir = filepath.Join(filepath.Dir(exe), "logs")

	err = os.MkdirAll(logsDir, 0700)
	if err != nil {
		return fmt.Errorf("could not create logs directory: %w", err)
	}

	return nil
}

var logFileName = "minyls_" + time.Now().Format("2006-01-02_15-04-05") + ".log.json"

func createLogFile() (io.Writer, error) {
	f, err := os.Create(filepath.Join(logsDir, logFileName))
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %w", err)
	}

	return f, nil
}

const logFileMaxAge = 7 * 24 * time.Hour

func removeOldLogs() {
	files, err := os.ReadDir(logsDir)
	if err != nil {
		fmt.Println("log.removeOldLogs: could not read logsDir:", err)
		return
	}

	for _, file := range files {
		fp := filepath.Join(logsDir, file.Name())

		if file.IsDir() {
			fmt.Println("log.removeOldLogs: unexpected dir:", fp)
			continue
		}

		info, err := file.Info()
		if err != nil {
			fmt.Println("log.removeOldLogs: could not read file info:", err)
			continue
		}

		if time.Since(info.ModTime()) <= logFileMaxAge {
			continue
		}

		err = os.Remove(fp)
		if err != nil {
			fmt.Println("could not remove file:", err)
			continue
		}
	}
}
