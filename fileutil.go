// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package fileutil implements utility functions related to files and paths.
package fileutil

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

const (
	// PrivateFileMode grants owner to read/write a file.
	PrivateFileMode = 0600
	// PrivateDirMode grants owner to make/remove files inside the directory.
	PrivateDirMode = 0700
)

// IsDirWriteable checks if dir is writable by writing and removing a file
// to dir. It returns nil if dir is writable.
func IsDirWriteable(dir string) error {
	f := filepath.Join(dir, ".touch")
	if err := ioutil.WriteFile(f, []byte(""), PrivateFileMode); err != nil {
		return err
	}
	return os.Remove(f)
}

// ReadDir returns the filenames in the given directory in sorted order.
func ReadDir(dirpath string) ([]string, error) {
	dir, err := os.Open(dirpath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	names, err := dir.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}

// TouchDirAll is similar to os.MkdirAll. It creates directories with 0700 permission if any directory
// does not exists. TouchDirAll also ensures the given directory is writable.
func TouchDirAll(dir string) error {
	// If path is already a directory, MkdirAll does nothing
	// and returns nil.
	err := os.MkdirAll(dir, PrivateDirMode)
	if err != nil {
		// if mkdirAll("a/text") and "text" is not
		// a directory, this will return syscall.ENOTDIR
		return err
	}
	return IsDirWriteable(dir)
}

// CreateDirAll is similar to TouchDirAll but returns error
// if the deepest directory was not empty.
func CreateDirAll(dir string) error {
	err := TouchDirAll(dir)
	if err == nil {
		var ns []string
		ns, err = ReadDir(dir)
		if err != nil {
			return err
		}
		if len(ns) != 0 {
			err = fmt.Errorf("expected %q to be empty, got %q", dir, ns)
		}
	}
	return err
}

func isErrInvalid(err error) bool {
	if err == os.ErrInvalid {
		return true
	}
	// Go < 1.8
	if syserr, ok := err.(*os.SyscallError); ok && syserr.Err == syscall.EINVAL {
		return true
	}
	// Go >= 1.8 returns *os.PathError instead
	if patherr, ok := err.(*os.PathError); ok && patherr.Err == syscall.EINVAL {
		return true
	}
	return false
}

// SyncDir call fsync() on a directory.
// Calling fsync() does not necessarily ensure that the entry in
// the directory containing the file has also reached disk. For
// that an explicit fsync() on a file descriptor for the directory
// is also needed.
func SyncDir(dir string) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil && !isErrInvalid(err) {
		return err
	}
	return nil
}

// Exist returns true if a file or directory exists.
func Exist(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

// ZeroToEnd zeros a file starting from SEEK_CUR to its SEEK_END. May temporarily
// shorten the length of the file.
func ZeroToEnd(f *os.File) error {
	// TODO: support FALLOC_FL_ZERO_RANGE
	off, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	lenf, lerr := f.Seek(0, io.SeekEnd)
	if lerr != nil {
		return lerr
	}
	if err = f.Truncate(off); err != nil {
		return err
	}
	// make sure blocks remain allocated
	if err = Preallocate(f, lenf, true); err != nil {
		return err
	}
	_, err = f.Seek(off, io.SeekStart)
	return err
}

// AtomicWriteFile **atomically** writes data to a file named by filename.
// If an error occurs, the target file is guaranteed to be either fully written,
// or not written at all.
// If the file does not exist, WriteFile creates it with permissions perm;
// otherwise WriteFile truncates it before writing.
func AtomicWriteFile(filename, tmpext string, data []byte, perm os.FileMode) (err error) {
	d, file := filepath.Split(filename)
	f, err := ioutil.TempFile(d, file+"*"+tmpext)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.Remove(f.Name()) // ignore
		}
	}()
	// write & fsync
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if ex := f.Sync(); err == nil {
		err = ex
	}
	if ex := f.Close(); err == nil {
		err = ex
	}
	if err != nil {
		return err
	}
	// chmod
	fi, err := os.Stat(filename)
	if err == nil {
		if err = os.Chmod(f.Name(), fi.Mode()); err != nil {
			return err
		}
	} else if os.IsNotExist(err) {
		if err = os.Chmod(f.Name(), perm); err != nil {
			return err
		}
	} else {
		return err
	}
	if err = os.Rename(f.Name(), filename); err != nil {
		return err
	}
	return SyncDir(d)
}
