package exifworker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
)

type ExifDataMap map[string]*Metadata

// ExifWorker wraps a persistent ExifTool process.
type ExifWorker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	reader *bufio.Reader
	mu     sync.Mutex // Ensures thread safety across goroutines
}

// NewExifWorker spawns a background ExifTool instance.
func NewExifWorker(executablePath string) (*ExifWorker, error) {
	if executablePath == "" {
		executablePath = "exiftool"
	}

	resolvedPath, err := exec.LookPath(executablePath)
	if err != nil {
		return nil, fmt.Errorf("exiftool executable %q is not installed or not in PATH: %w", executablePath, err)
	}

	cmd := exec.Command(resolvedPath, "-stay_open", "True", "-@", "-")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to open stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to open stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start exiftool process: %w", err)
	}

	return &ExifWorker{
		cmd:    cmd,
		stdin:  stdin,
		reader: bufio.NewReader(stdout),
	}, nil
}

// GetMetadata extracts raw JSON tags from a target file.
func (w *ExifWorker) GetMetadata(result *[]map[string]any, filePath string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Write arguments line-by-line to STDIN
	// -j = Format output as JSON
	// <filePath> = Target file
	// -execute = Run queued commands and emit {ready} on completion
	cmdPayload := fmt.Sprintf(`-j
	-ExifToolVersion
	-MIMEType
	-FileType
	-FileTypeExtension
	-FileSize#
	-ImageWidth
	-ImageHeight
	%s
	-execute
	`, filePath)

	if _, err := io.WriteString(w.stdin, cmdPayload); err != nil {
		return fmt.Errorf("failed to write command to stdin: %w", err)
	}

	var outputBuf bytes.Buffer

	// Read STDOUT line-by-line until ExifTool emits the {ready} line
	for {
		line, err := w.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading stdout stream: %w", err)
		}

		if strings.TrimSpace(line) == "{ready}" {
			break
		}

		outputBuf.WriteString(line)
	}

	if err := json.Unmarshal(outputBuf.Bytes(), result); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return nil
}

// Close gracefully terminates the running ExifTool process.
func (w *ExifWorker) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Signal ExifTool to exit
	_, _ = io.WriteString(w.stdin, "-stay_open\nFalse\n")
	_ = w.stdin.Close()

	return w.cmd.Wait()
}

func (w *ExifWorker) ReadFile(result ExifDataMap, filePath string) error {

	w.mu.Lock()
	defer w.mu.Unlock()

	// Write arguments line-by-line to STDIN
	// -j = Format output as JSON
	// <filePath> = Target file
	// -execute = Run queued commands and emit {ready} on completion
	cmdPayload := fmt.Sprintf(`-j
	-ExifToolVersion
	-SourceFile
	-FileName
	-FileModifyDate
	-MIMEType
	-FileType
	-FileTypeExtension
	-FileSize#
	-ImageWidth
	-ImageHeight
	%s
	-execute
	`, filePath)

	if _, err := io.WriteString(w.stdin, cmdPayload); err != nil {
		return fmt.Errorf("failed to write command to stdin: %w", err)
	}

	var outputBuf bytes.Buffer

	// Read STDOUT line-by-line until ExifTool emits the {ready} line
	for {
		line, err := w.reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("error reading stdout stream: %w", err)
		}

		if strings.TrimSpace(line) == "{ready}" {
			break
		}

		outputBuf.WriteString(line)
	}

	items := []map[string]any{}
	if err := json.Unmarshal(outputBuf.Bytes(), &items); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// 3. Populate the map
	for _, item := range items {
		key := getString(item, "FileName")
		if key == "" {
			key = filePath
		}
		result[key] = &Metadata{data: item, key: key}
	}

	return nil
}

// --- Modified Readers ---

type Metadata struct {
	data map[string]any
	key  string
}

func (m *Metadata) ExifToolVersion() float64 {
	return getFloat64(m.data, "ExifToolVersion")
}

func (m *Metadata) SourceFile() string {
	return getString(m.data, "SourceFile")
}

func (m *Metadata) FileName() string {
	return getString(m.data, "FileName")
}

func (m *Metadata) ModifiedDate() string {
	return getString(m.data, "FileModifyDate")
}

func (m *Metadata) MimeType() string {
	return getString(m.data, "MIMEType")
}

func (m *Metadata) FileType() string {
	return getString(m.data, "FileType")
}

func (m *Metadata) FileExt() string {
	return getString(m.data, "FileTypeExtension")
}

func (m *Metadata) FileSizeBytes() int64 {
	return getInt64(m.data, "FileSize")
}

func (m *Metadata) ImageWidth() int {
	return getInt(m.data, "ImageWidth")
}

func (m *Metadata) ImageHeight() int {
	return getInt(m.data, "ImageHeight")
}

// helpers

func getString(m map[string]any, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getFloat64(m map[string]any, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0.0
}

func getInt64(m map[string]any, key string) int64 {
	switch v := m[key].(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	default:
		return 0
	}
}

func getInt(m map[string]any, key string) int {
	switch v := m[key].(type) {
	case float64:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	default:
		return 0
	}
}
