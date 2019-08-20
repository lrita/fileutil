// Copyright 2016 The etcd Authors
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

// +build linux

package fileutil

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

// Fsync is a wrapper around file.Sync(). Special handling is needed on darwin platform.
func Fsync(f *os.File) error {
	return f.Sync()
}

// Fdatasync is similar to fsync(), but does not flush modified metadata
// unless that metadata is needed in order to allow a subsequent data retrieval
// to be correctly handled.
func Fdatasync(f *os.File) error {
	return syscall.Fdatasync(int(f.Fd()))
}

// Frangesync invokes sync_file_range syscall for the given offset and nbytes.
func Frangesync(f *os.File, off, n int64, strict bool) error {
	flag := unix.SYNC_FILE_RANGE_WRITE
	if strict {
		flag |= unix.SYNC_FILE_RANGE_WAIT_BEFORE
	}
	return syscall.SyncFileRange(int(f.Fd()), off, n, flag)
}

// IsSyncFileRangeSupported returns true when the filesystem which the given
// file is located in is support sync_file_range syscall.
func IsSyncFileRangeSupported(f *os.File) bool {
	var st syscall.Statfs_t
	syscall.Fstatfs(int(f.Fd()), &st)
	if st.Type == FsMagicZfs {
		// Testing on ZFS showed the writeback did not happen asynchronously when
		// `sync_file_range` was called, even though it returned success.
		return false
	}
	err := syscall.SyncFileRange(int(f.Fd()), 0, 0, 0)
	return err != syscall.ENOSYS
}
