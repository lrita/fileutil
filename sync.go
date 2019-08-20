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

// +build !linux,!darwin

package fileutil

import "os"

// Fsync is a wrapper around file.Sync(). Special handling is needed on darwin platform.
func Fsync(f *os.File) error {
	return f.Sync()
}

// Fdatasync is a wrapper around file.Sync(). Special handling is needed on linux platform.
func Fdatasync(f *os.File) error {
	return f.Sync()
}

// Frangesync invokes sync_file_range syscall for the given offset and nbytes.
func Frangesync(f *os.File, off, n int64, strict bool) error {
	return syscall.ENOSYS
}

// IsSyncFileRangeSupported returns true when the filesystem which the given
// file is located in is support sync_file_range syscall.
func IsSyncFileRangeSupported(f *os.File) bool {
	return false
}
