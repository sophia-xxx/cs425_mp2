package logger

import (
	"cs425_mp2/config"
	"fmt"
	"log"
	"os"
)

var (
	// console
	ConsoleLogger *log.Logger

	// InfoLogger : logging for info messages
	InfoLogger *log.Logger

	// WarningLogger : logging for warning messages
	WarningLogger *log.Logger

	// ErrorLogger : logging for error messages
	ErrorLogger *log.Logger

	// ErrorLogger : logging for error messages
	DebugLogger *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	ConsoleLogger = log.New(os.Stdout, "", 0)
	InfoLogger = log.New(file, "[INFO]", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "[WARNING]", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "[ERROR]", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "[DEBUG]", log.Ldate|log.Ltime|log.Lshortfile)
}

func PrintToConsole(args ...interface{}) {
	ConsoleLogger.Println(args)
}

func PrintInfo(args ...interface{}) {
	InfoLogger.Println(args)

	fmt.Print("[INFO] ")
	fmt.Println(args...)
}

func PrintWarning(args ...interface{}) {
	WarningLogger.Println(args)

	fmt.Print("[WARNING] ")
	fmt.Println(args...)
}

func PrintError(args ...interface{}) {
	ErrorLogger.Println(args)

	fmt.Print("[ERROR] ")
	fmt.Println(args...)
}

func PrintDebug(args ...interface{}) {
	ErrorLogger.Println(args)
	if config.DebugMode {
		fmt.Print("[DEBUG] ")
		fmt.Println(args...)
	}

}
