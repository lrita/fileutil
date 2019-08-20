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

package fileutil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFileWriteWithAllSyncMethod(t *testing.T) {
	f, err := ioutil.TempFile("", "sync_file_testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	supportRangeSync := IsSyncFileRangeSupported(f)

	n, err := f.Write([]byte("abcdefg"))
	if err != nil {
		t.Errorf("write failed: %v", err)
	}

	if supportRangeSync {
		err = Frangesync(f, 0, int64(n), true)
	} else {
		err = Fdatasync(f)
	}
	if err != nil {
		t.Errorf("sync failed %v", err)
	}
	if err = Fsync(f); err != nil {
		t.Errorf("sync failed %v", err)
	}
	f.Close()
}
