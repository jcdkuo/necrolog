package necrolog

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* EnforceLogRetention manages log files by deleting oversized files and removing non-whitelisted entries.
 * It cleans the log directory by:
 * 1. Deleting files exceeding 1.2MB in size
 * 2. Deleting files not in the whitelist
 * 3. Deleting compressed (.gz) files whose base name prefix is not in the whitelist
 */
func EnforceLogRetention() error {
	// Ensure directory exists
	if err := EnsurePathExists(LogDeadRisingPath); err != nil {
		Error(LogErrorFile, fmt.Sprintf("Failed to ensure log directory exists: %v", err))
		return err
	}

	// Read all files in the directory
	files, err := os.ReadDir(LogDeadRisingPath)
	if err != nil {
		Error(LogErrorFile, fmt.Sprintf("Failed to read log directory: %v", err))
		return err
	}

	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(LogDeadRisingPath, file.Name())
		
		// Handle file based on its type and whitelist status
		handleFile(file, filePath)
	}

	return nil
}

// handleFile processes a single file based on its type and whitelist status
func handleFile(file os.DirEntry, filePath string) {
	// Get file info for size checking (common operation)
	fileInfo, err := file.Info()
	if err != nil {
		Error(LogErrorFile, fmt.Sprintf("Failed to get file info for %s (removing): %v", filePath, err))
		removeFile(filePath, "Failed to get file info")
		return
	}

	// Check if file is in whitelist
	if WhitelistedFiles[filePath] {
		// Only delete if oversized
		if fileInfo.Size() > MaxLogFileSize {
			removeFile(filePath, fmt.Sprintf("Oversized whitelist file (%.2f MB)", float64(fileInfo.Size())/(1024*1024)))
		}
		return
	}
	
	// Handle compressed files (.gz)
	if strings.HasSuffix(filePath, ".gz") {
		handleCompressedFile(fileInfo, filePath)
		return
	}

	// Not in whitelist and not a compressed file, delete
	removeFile(filePath, "Non-whitelist file")
}

// handleCompressedFile processes compressed (.gz) files
func handleCompressedFile(fileInfo os.FileInfo, filePath string) {
	baseNameWithoutExt := strings.Split(filepath.Base(filePath), "-")[0] // Get base name before first hyphen
	originalLogFile := filepath.Join(LogDeadRisingPath, baseNameWithoutExt+".log")
	
	// If original log file is in whitelist, check size, otherwise delete
	if WhitelistedFiles[originalLogFile] {
		if fileInfo.Size() > MaxLogFileSize {
			removeFile(filePath, fmt.Sprintf("Oversized compressed file (%.2f MB)", float64(fileInfo.Size())/(1024*1024)))
		}
	} else {
		removeFile(filePath, "Compressed file with non-whitelisted prefix")
	}
}

// removeFile attempts to delete a file and logs the result
func removeFile(filePath, reason string) {
	if err := os.Remove(filePath); err != nil {
		Error(LogErrorFile, fmt.Sprintf("Failed to delete %s: %v", filePath, err))
	} else {
		Warn(LogWarningFile, fmt.Sprintf("Deleted %s: %s", filePath, reason))
	}
}

func EnsurePathExists(path string) error {
    // Check if the path exists
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        // Create the directory if it doesn't exist
        err := os.MkdirAll(path, 0755) // 0755 sets the directory permissions
        if err != nil {
            return fmt.Errorf("failed to create directory %s: %w", path, err)
        }
        Warn(LogWarningFile, fmt.Sprintf("Directory created: %s", path))
    } else if err != nil {
        // Handle other errors (e.g., permission issues)
        Error(LogErrorFile, fmt.Sprintf("Error checking directory %s: %v", path, err))
        return err
    } else {
        Warn(LogWarningFile, fmt.Sprintf("Directory already exists: %s", path))
    }
    return nil
}
