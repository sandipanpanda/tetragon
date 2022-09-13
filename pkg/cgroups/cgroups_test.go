// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Tetragon

package cgroups

import (
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/unix"
)

// Test cgroup mode detection on an invalid directory
func TestDetectCgroupModeInvalid(t *testing.T) {
	mode, err := detectCgroupMode("invalid-cgroupfs-path")
	assert.Error(t, err)
	assert.Equal(t, CGROUP_UNDEF, mode)
}

// Test cgroup mode detection on default cgroup root /sys/fs/cgroup
func TestDetectCgroupModeDefault(t *testing.T) {
	var st syscall.Statfs_t

	err := syscall.Statfs(defaultCgroupRoot, &st)
	if err != nil {
		t.Skipf("TestDetectCgroupModeDefault() failed to Statfs(%s): %v, test skipped", defaultCgroupRoot, err)
	}

	mode, err := detectCgroupMode(defaultCgroupRoot)
	assert.NoError(t, err)

	if st.Type == unix.CGROUP2_SUPER_MAGIC {
		assert.Equal(t, CGROUP_UNIFIED, mode)
	} else if st.Type == unix.TMPFS_MAGIC {
		unified := filepath.Join(defaultCgroupRoot, "unified")
		err = syscall.Statfs(unified, &st)
		if err == nil && st.Type == unix.CGROUP2_SUPER_MAGIC {
			assert.Equal(t, CGROUP_HYBRID, mode)

			// Extra detection
			mode, err = detectCgroupMode(unified)
			assert.NoError(t, err)
			assert.Equal(t, CGROUP_UNIFIED, mode)
		} else {
			assert.Equal(t, CGROUP_LEGACY, mode)
		}
	} else {
		t.Errorf("TestDetectCgroupModeDefault() failed Cgroupfs %s type failed:  want:%d or %d -  got:%d",
			defaultCgroupRoot, unix.CGROUP2_SUPER_MAGIC, unix.TMPFS_MAGIC, st.Type)
	}
}