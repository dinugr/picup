package workflow

import "io"

// FileStorage contains the filesystem operations required by an image workflow.
type FileStorage interface {
	Stage(source io.Reader, extension string) (*FileTask, error)
	Copy(sourcePath string) (*FileTask, error)
	Save(task *FileTask, targetFilename string) (string, error)
	Delete(paths ...string) error
}
