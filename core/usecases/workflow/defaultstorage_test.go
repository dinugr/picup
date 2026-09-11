package workflow

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultStorageStageSaveAndClean(t *testing.T) {
	root := t.TempDir()
	helper := NewDefaultStorage(filepath.Join(root, "temp"), filepath.Join(root, "data"))

	task, err := helper.Stage(bytes.NewBufferString("image-data"), ".jpg")
	if err != nil {
		t.Fatalf("Stage() error = %v", err)
	}
	workPath := task.WorkFilepath()

	target, err := helper.Save(task, filepath.Join("master", "image.jpg"))
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	contents, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(contents) != "image-data" {
		t.Fatalf("saved contents = %q, want %q", contents, "image-data")
	}
	if _, err := os.Stat(workPath); !os.IsNotExist(err) {
		t.Fatalf("work file still exists after save: %v", err)
	}

	task.Clean()
}

func TestDefaultStorageCopyAndDelete(t *testing.T) {
	root := t.TempDir()
	helper := NewDefaultStorage(filepath.Join(root, "temp"), filepath.Join(root, "data"))
	source := filepath.Join(root, "source.jpg")
	if err := os.WriteFile(source, []byte("source-data"), 0600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	task, err := helper.Copy(source)
	if err != nil {
		t.Fatalf("Copy() error = %v", err)
	}
	defer task.Clean()

	target, err := helper.Save(task, "variant.jpg")
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if err := helper.Delete(target); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("target still exists: %v", err)
	}
}

func TestFileTaskCleanNilSafe(t *testing.T) {
	var task *FileTask
	task.Clean()
}
