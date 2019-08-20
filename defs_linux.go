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

const (
	// FsMagicAufs filesystem id for Aufs
	FsMagicAufs = 0x61756673
	// FsMagicBtrfs filesystem id for Btrfs
	FsMagicBtrfs = 0x9123683E
	// FsMagicCramfs filesystem id for Cramfs
	FsMagicCramfs = 0x28cd3d45
	// FsMagicEcryptfs filesystem id for eCryptfs
	FsMagicEcryptfs = 0xf15f
	// FsMagicExtfs filesystem id for Extfs
	FsMagicExtfs = 0x0000EF53
	// FsMagicF2fs filesystem id for F2fs
	FsMagicF2fs = 0xF2F52010
	// FsMagicGPFS filesystem id for GPFS
	FsMagicGPFS = 0x47504653
	// FsMagicJffs2Fs filesystem if for Jffs2Fs
	FsMagicJffs2Fs = 0x000072b6
	// FsMagicJfs filesystem id for Jfs
	FsMagicJfs = 0x3153464a
	// FsMagicNfsFs filesystem id for NfsFs
	FsMagicNfsFs = 0x00006969
	// FsMagicRAMFs filesystem id for RamFs
	FsMagicRAMFs = 0x858458f6
	// FsMagicReiserFs filesystem id for ReiserFs
	FsMagicReiserFs = 0x52654973
	// FsMagicSmbFs filesystem id for SmbFs
	FsMagicSmbFs = 0x0000517B
	// FsMagicSquashFs filesystem id for SquashFs
	FsMagicSquashFs = 0x73717368
	// FsMagicTmpFs filesystem id for TmpFs
	FsMagicTmpFs = 0x01021994
	// FsMagicVxFS filesystem id for VxFs
	FsMagicVxFS = 0xa501fcf5
	// FsMagicXfs filesystem id for Xfs
	FsMagicXfs = 0x58465342
	// FsMagicZfs filesystem id for Zfs
	FsMagicZfs = 0x2fc12fc1
	// FsMagicOverlay filesystem id for overlay
	FsMagicOverlay = 0x794C7630
)
