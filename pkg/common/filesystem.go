package common

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bthompso/engx-ergonomics-poc/pkg/common/interfaces"
)

// FilesystemManager implements the interfaces.FilesystemManager interface
type FilesystemManager struct {
	operations []FileOperation
	logger     interfaces.Logger
	mutex      sync.RWMutex
}

// FileOperation represents a filesystem operation for tracking
type FileOperation struct {
	Type      string `json:"type"`      // "create", "read", "write", "delete", "copy", "move"
	Path      string `json:"path"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
	Size      int64  `json:"size,omitempty"`
	Timestamp string `json:"timestamp"`
}

// NewFilesystemManager creates a new filesystem manager instance
func NewFilesystemManager(logger interfaces.Logger) interfaces.FilesystemManager {
	return &FilesystemManager{
		operations: make([]FileOperation, 0),
		logger:     logger.WithComponent("filesystem"),
		mutex:      sync.RWMutex{},
	}
}

// AppendToFile appends data to an existing file
func (fm *FilesystemManager) AppendToFile(path string, data []byte) error {
	fm.logger.Debug("Appending to file: %s (%d bytes)", path, len(data))

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fm.recordOperation("append", path, false, err)
		fm.logger.Error("Failed to open file for append %s: %v", path, err)
		return fmt.Errorf("failed to open file for append %s: %w", path, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	fm.recordOperation("append", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to append to file %s: %v", path, err)
		return fmt.Errorf("failed to append to file %s: %w", path, err)
	}

	fm.logger.Debug("Successfully appended to file: %s", path)
	return nil
}

// CreateDir creates a directory with specified permissions
func (fm *FilesystemManager) CreateDir(path string, perm uint32) error {
	fm.logger.Debug("Creating directory: %s", path)

	err := os.MkdirAll(path, os.FileMode(perm))
	fm.recordOperation("create_dir", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to create directory %s: %v", path, err)
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	fm.logger.Debug("Successfully created directory: %s", path)
	return nil
}

// RemoveDir removes a directory and all its contents
func (fm *FilesystemManager) RemoveDir(path string) error {
	fm.logger.Debug("Removing directory: %s", path)

	err := os.RemoveAll(path)
	fm.recordOperation("remove_dir", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to remove directory %s: %v", path, err)
		return fmt.Errorf("failed to remove directory %s: %w", path, err)
	}

	fm.logger.Debug("Successfully removed directory: %s", path)
	return nil
}

// ListDir lists the contents of a directory
func (fm *FilesystemManager) ListDir(path string) ([]string, error) {
	return fm.ListDirectory(path)
}

// Exists checks if a file or directory exists
func (fm *FilesystemManager) Exists(path string) bool {
	return fm.FileExists(path)
}

// IsDir checks if the path is a directory
func (fm *FilesystemManager) IsDir(path string) bool {
	return fm.IsDirectory(path)
}

// IsFile checks if the path is a regular file
func (fm *FilesystemManager) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	isFile := !info.IsDir()
	fm.logger.Debug("File check: %s = %t", path, isFile)
	return isFile
}

// GetWorkingDir returns the current working directory
func (fm *FilesystemManager) GetWorkingDir() (string, error) {
	return fm.GetWorkingDirectory()
}

// ChangeDir changes the current working directory
func (fm *FilesystemManager) ChangeDir(path string) error {
	return fm.ChangeDirectory(path)
}

// MoveFile moves a file from source to destination
func (fm *FilesystemManager) MoveFile(src, dst string) error {
	fm.logger.Debug("Moving file: %s -> %s", src, dst)

	// Create destination directory if needed
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		fm.recordOperation("move", fmt.Sprintf("%s->%s", src, dst), false, err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	err := os.Rename(src, dst)
	fm.recordOperation("move", fmt.Sprintf("%s->%s", src, dst), err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to move file %s to %s: %v", src, dst, err)
		return fmt.Errorf("failed to move file %s to %s: %w", src, dst, err)
	}

	fm.logger.Debug("Successfully moved file: %s -> %s", src, dst)
	return nil
}

// GetTempDir creates and returns a temporary directory
func (fm *FilesystemManager) GetTempDir() (string, error) {
	fm.logger.Debug("Creating temporary directory")

	tmpDir, err := os.MkdirTemp("", "engx-*")
	if err != nil {
		fm.recordOperation("temp_dir", "", false, err)
		fm.logger.Error("Failed to create temporary directory: %v", err)
		return "", fmt.Errorf("failed to create temporary directory: %w", err)
	}

	fm.recordOperation("temp_dir", tmpDir, true, nil)
	fm.logger.Debug("Successfully created temporary directory: %s", tmpDir)
	return tmpDir, nil
}

// CreateDirectory creates a directory and all necessary parent directories
func (fm *FilesystemManager) CreateDirectory(path string) error {
	fm.logger.Debug("Creating directory: %s", path)

	err := os.MkdirAll(path, 0755)

	fm.recordOperation("create_dir", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to create directory %s: %v", path, err)
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	fm.logger.Debug("Successfully created directory: %s", path)
	return nil
}

// WriteFile writes data to a file, creating directories as needed
func (fm *FilesystemManager) WriteFile(path string, data []byte, perm uint32) error {
	fm.logger.Debug("Writing file: %s (%d bytes)", path, len(data))

	// Create parent directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fm.recordOperation("write", path, false, err)
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	// Write the file using the specified permissions
	err := os.WriteFile(path, data, os.FileMode(perm))

	fm.recordOperation("write", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to write file %s: %v", path, err)
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	fm.logger.Debug("Successfully wrote file: %s", path)
	return nil
}

// ReadFile reads data from a file
func (fm *FilesystemManager) ReadFile(path string) ([]byte, error) {
	fm.logger.Debug("Reading file: %s", path)

	data, err := os.ReadFile(path)

	size := int64(0)
	if err == nil {
		size = int64(len(data))
	}

	fm.recordOperationWithSize("read", path, err == nil, err, size)

	if err != nil {
		fm.logger.Error("Failed to read file %s: %v", path, err)
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	fm.logger.Debug("Successfully read file: %s (%d bytes)", path, len(data))
	return data, nil
}

// CopyFile copies a file from source to destination
func (fm *FilesystemManager) CopyFile(src, dst string) error {
	fm.logger.Debug("Copying file: %s -> %s", src, dst)

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		fm.recordOperation("copy", fmt.Sprintf("%s->%s", src, dst), false, err)
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer srcFile.Close()

	// Create destination directory if needed
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		fm.recordOperation("copy", fmt.Sprintf("%s->%s", src, dst), false, err)
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		fm.recordOperation("copy", fmt.Sprintf("%s->%s", src, dst), false, err)
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer dstFile.Close()

	// Copy data
	_, err = io.Copy(dstFile, srcFile)

	fm.recordOperation("copy", fmt.Sprintf("%s->%s", src, dst), err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to copy file %s to %s: %v", src, dst, err)
		return fmt.Errorf("failed to copy file: %w", err)
	}

	fm.logger.Debug("Successfully copied file: %s -> %s", src, dst)
	return nil
}

// DeleteFile deletes a file
func (fm *FilesystemManager) DeleteFile(path string) error {
	fm.logger.Debug("Deleting file: %s", path)

	err := os.Remove(path)

	fm.recordOperation("delete", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to delete file %s: %v", path, err)
		return fmt.Errorf("failed to delete file %s: %w", path, err)
	}

	fm.logger.Debug("Successfully deleted file: %s", path)
	return nil
}

// FileExists checks if a file or directory exists
func (fm *FilesystemManager) FileExists(path string) bool {
	_, err := os.Stat(path)
	exists := !os.IsNotExist(err)

	fm.logger.Debug("File existence check: %s = %t", path, exists)
	return exists
}

// IsDirectory checks if the path is a directory
func (fm *FilesystemManager) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	isDir := info.IsDir()
	fm.logger.Debug("Directory check: %s = %t", path, isDir)
	return isDir
}

// ListDirectory lists the contents of a directory
func (fm *FilesystemManager) ListDirectory(path string) ([]string, error) {
	fm.logger.Debug("Listing directory: %s", path)

	entries, err := os.ReadDir(path)
	if err != nil {
		fm.recordOperation("list", path, false, err)
		fm.logger.Error("Failed to list directory %s: %v", path, err)
		return nil, fmt.Errorf("failed to list directory %s: %w", path, err)
	}

	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}

	fm.recordOperation("list", path, true, nil)
	fm.logger.Debug("Successfully listed directory: %s (%d entries)", path, len(names))
	return names, nil
}

// GetFileInfo returns information about a file
func (fm *FilesystemManager) GetFileInfo(path string) (map[string]interface{}, error) {
	fm.logger.Debug("Getting file info: %s", path)

	info, err := os.Stat(path)
	if err != nil {
		fm.recordOperation("stat", path, false, err)
		fm.logger.Error("Failed to get file info for %s: %v", path, err)
		return nil, fmt.Errorf("failed to get file info for %s: %w", path, err)
	}

	fileInfo := map[string]interface{}{
		"name":    info.Name(),
		"size":    info.Size(),
		"mode":    info.Mode().String(),
		"modtime": info.ModTime(),
		"is_dir":  info.IsDir(),
	}

	fm.recordOperation("stat", path, true, nil)
	fm.logger.Debug("Successfully got file info: %s", path)
	return fileInfo, nil
}

// CleanPath cleans and validates a file path
func (fm *FilesystemManager) CleanPath(path string) string {
	cleaned := filepath.Clean(path)
	fm.logger.Debug("Cleaned path: %s -> %s", path, cleaned)
	return cleaned
}

// JoinPath joins path elements together
func (fm *FilesystemManager) JoinPath(elements ...string) string {
	joined := filepath.Join(elements...)
	fm.logger.Debug("Joined path: %v -> %s", elements, joined)
	return joined
}

// GetWorkingDirectory returns the current working directory
func (fm *FilesystemManager) GetWorkingDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		fm.logger.Error("Failed to get working directory: %v", err)
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	fm.logger.Debug("Current working directory: %s", wd)
	return wd, nil
}

// ChangeDirectory changes the current working directory
func (fm *FilesystemManager) ChangeDirectory(path string) error {
	fm.logger.Debug("Changing directory to: %s", path)

	err := os.Chdir(path)

	fm.recordOperation("chdir", path, err == nil, err)

	if err != nil {
		fm.logger.Error("Failed to change directory to %s: %v", path, err)
		return fmt.Errorf("failed to change directory to %s: %w", path, err)
	}

	fm.logger.Debug("Successfully changed directory to: %s", path)
	return nil
}

// CreateTempFile creates a temporary file
func (fm *FilesystemManager) CreateTempFile(pattern string) (string, error) {
	fm.logger.Debug("Creating temporary file with pattern: %s", pattern)

	tmpFile, err := os.CreateTemp("", pattern)
	if err != nil {
		fm.recordOperation("temp", pattern, false, err)
		fm.logger.Error("Failed to create temporary file: %v", err)
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}

	path := tmpFile.Name()
	tmpFile.Close()

	fm.recordOperation("temp", path, true, nil)
	fm.logger.Debug("Successfully created temporary file: %s", path)
	return path, nil
}

// GetOperationHistory returns the history of filesystem operations
func (fm *FilesystemManager) GetOperationHistory() []map[string]interface{} {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	var history []map[string]interface{}
	for _, op := range fm.operations {
		history = append(history, map[string]interface{}{
			"type":      op.Type,
			"path":      op.Path,
			"success":   op.Success,
			"error":     op.Error,
			"size":      op.Size,
			"timestamp": op.Timestamp,
		})
	}

	return history
}

// GetStatistics returns filesystem operation statistics
func (fm *FilesystemManager) GetStatistics() map[string]interface{} {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_operations": len(fm.operations),
		"operations_by_type": make(map[string]int),
		"successful_operations": 0,
		"failed_operations": 0,
	}

	operationsByType := make(map[string]int)
	successfulOps := 0
	failedOps := 0

	for _, op := range fm.operations {
		operationsByType[op.Type]++
		if op.Success {
			successfulOps++
		} else {
			failedOps++
		}
	}

	stats["operations_by_type"] = operationsByType
	stats["successful_operations"] = successfulOps
	stats["failed_operations"] = failedOps

	if len(fm.operations) > 0 {
		stats["success_rate"] = float64(successfulOps) / float64(len(fm.operations))
	} else {
		stats["success_rate"] = 0.0
	}

	return stats
}

// recordOperation records a filesystem operation for tracking
func (fm *FilesystemManager) recordOperation(opType, path string, success bool, err error) {
	fm.recordOperationWithSize(opType, path, success, err, 0)
}

// recordOperationWithSize records a filesystem operation with size information
func (fm *FilesystemManager) recordOperationWithSize(opType, path string, success bool, err error, size int64) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}

	operation := FileOperation{
		Type:      opType,
		Path:      path,
		Success:   success,
		Error:     errorMsg,
		Size:      size,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}

	fm.operations = append(fm.operations, operation)

	// Keep only the last 1000 operations to prevent memory leaks
	if len(fm.operations) > 1000 {
		fm.operations = fm.operations[len(fm.operations)-1000:]
	}
}

// Cleanup performs cleanup of filesystem manager resources
func (fm *FilesystemManager) Cleanup() error {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// Clear operation history
	fm.operations = make([]FileOperation, 0)

	fm.logger.Debug("Filesystem manager cleaned up")
	return nil
}