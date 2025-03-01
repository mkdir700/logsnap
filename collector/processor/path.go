package processor

import (
	"os/user"
	"path/filepath"
	"strings"
)

// Path represents a filesystem path
type Path string

// expandHomeDir replaces ~ or ~user with the corresponding home directory
func expandHomeDir(path string) (string, error) {
    if !strings.HasPrefix(path, "~") {
        return path, nil
    }

    var username string
    var rest string

    // Split path into ~user and rest of path
    parts := strings.SplitN(path[1:], "/", 2)
    if len(parts) == 1 {
        username = ""
        rest = parts[0]
    } else {
        username = parts[0]
        rest = parts[1]
    }

    var homeDir string
    if username == "" {
        // Get current user's home directory
        currentUser, err := user.Current()
        if err != nil {
            return "", err
        }
        homeDir = currentUser.HomeDir
    } else {
        // Get specified user's home directory
        u, err := user.Lookup(username)
        if err != nil {
            return "", err
        }
        homeDir = u.HomeDir
    }

    // Join home directory with rest of path
    if rest == "" {
        return homeDir, nil
    }
    return filepath.Join(homeDir, rest), nil
}

// NewPath creates a new Path with proper formatting
func NewPath(path string) Path {
    // Clean the path to remove any redundant separators and resolve ".." and "."
    return Path(filepath.Clean(path))
}

// String returns the string representation of the path
func (p Path) String() string {
    return string(p)
}

// Get Absolute Path
func (p Path) GetAbsolutePath() (string, error) {
	// First expand the home directory if path contains ~
	expandedPath, err := expandHomeDir(string(p))
	if err != nil {
		return "", err
	}
	return expandedPath, nil
}

// Join joins path elements with this path
func (p Path) Join(elem ...string) Path {
    elements := make([]string, 0, len(elem)+1)
    elements = append(elements, string(p))
    elements = append(elements, elem...)
    return Path(filepath.Join(elements...))
} 