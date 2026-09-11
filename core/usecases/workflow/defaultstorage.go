package workflow

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// DefaultStorage is the local filesystem implementation of FileStorage.
type DefaultStorage struct {
	workdir string
	datadir string
}

// FileTask owns the temporary files used by one workflow operation.
type FileTask struct {
	workdir  string
	workpath string
	previous string
}

func NewDefaultStorage(workdir, datadir string) *DefaultStorage {
	return &DefaultStorage{workdir: workdir, datadir: datadir}
}

func (f *DefaultStorage) Stage(source io.Reader, extension string) (*FileTask, error) {
	task, err := f.newTask(extension)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(task.workpath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		task.Clean()
		return nil, fmt.Errorf("open staged file: %w", err)
	}

	_, copyErr := io.Copy(file, source)
	closeErr := file.Close()
	if copyErr != nil {
		task.Clean()
		return nil, fmt.Errorf("write staged file: %w", copyErr)
	}
	if closeErr != nil {
		task.Clean()
		return nil, fmt.Errorf("close staged file: %w", closeErr)
	}
	return task, nil
}

func (f *DefaultStorage) Copy(sourcePath string) (*FileTask, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("open source file: %w", err)
	}
	defer source.Close()

	task, err := f.newTask(filepath.Ext(sourcePath))
	if err != nil {
		return nil, err
	}

	destination, err := os.OpenFile(task.workpath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		task.Clean()
		return nil, fmt.Errorf("open work file: %w", err)
	}
	_, copyErr := io.Copy(destination, source)
	closeErr := destination.Close()
	if copyErr != nil {
		task.Clean()
		return nil, fmt.Errorf("copy source file: %w", copyErr)
	}
	if closeErr != nil {
		task.Clean()
		return nil, fmt.Errorf("close work file: %w", closeErr)
	}
	return task, nil
}

func (f *DefaultStorage) newTask(extension string) (*FileTask, error) {
	if err := os.MkdirAll(f.workdir, 0755); err != nil {
		return nil, fmt.Errorf("create work directory: %w", err)
	}
	path := filepath.Join(f.workdir, uuid.NewString()+extension)
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("create work file: %w", err)
	}
	if err := file.Close(); err != nil {
		os.Remove(path)
		return nil, fmt.Errorf("close work file: %w", err)
	}
	return &FileTask{workdir: f.workdir, workpath: path}, nil
}

// WorkFilepath returns the temporary file path for this specific worker task.
func (f *DefaultStorage) Save(task *FileTask, targetFilename string) (string, error) {
	if task == nil || task.workpath == "" {
		return "", fmt.Errorf("no work file available to save")
	}

	targetPath := filepath.Join(f.datadir, targetFilename)
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", fmt.Errorf("create data directory: %w", err)
	}

	if err := os.Rename(task.workpath, targetPath); err == nil {
		task.workpath = ""
		return targetPath, nil
	}

	source, err := os.Open(task.workpath)
	if err != nil {
		return "", fmt.Errorf("open work file for saving: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("create data file: %w", err)
	}
	if _, err := io.Copy(destination, source); err != nil {
		destination.Close()
		os.Remove(targetPath)
		return "", fmt.Errorf("copy work file to data: %w", err)
	}
	if err := destination.Close(); err != nil {
		os.Remove(targetPath)
		return "", fmt.Errorf("close data file: %w", err)
	}

	task.Clean()
	return targetPath, nil
}

func (f *DefaultStorage) Delete(paths ...string) error {
	var firstErr error
	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) && firstErr == nil {
			firstErr = fmt.Errorf("delete %q: %w", path, err)
		}
	}
	return firstErr
}

func (t *FileTask) WorkFilepath() string { return t.workpath }

func (t *FileTask) NewWorkFilepath(filename string) string {
	t.previous = t.workpath
	t.workpath = filepath.Join(t.workdir, filename)
	return t.workpath
}

func (t *FileTask) Clean() {
	if t == nil {
		return
	}
	if t.workpath != "" {
		os.Remove(t.workpath)
		t.workpath = ""
	}
	if t.previous != "" {
		os.Remove(t.previous)
		t.previous = ""
	}
}
