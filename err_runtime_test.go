package std

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ErrorRuntime(t *testing.T) {
	// permanent
	for i, err := range []ErrorRuntime{
		// deprecated
		NewErrorRuntimeFf("default"),
		WrapErrorToRuntime(fmt.Errorf("err"), "err", "err"),
		// new
		NewErrorRuntimePermanentFf("err"),
		// default wrap + default
		ErrToRuntime(fmt.Errorf("err"), "err", "err"),
		// default wrap + permanent
		ErrToRuntime(NewErrorRuntimePermanentFf("err"), "err", "err"),
		ErrToRuntime(fmt.Errorf("wrap: %w", NewErrorRuntimePermanentFf("err")), "err", "err"),
		// force wrap + default
		ErrToRuntimePermanent(fmt.Errorf("err"), "err", "err"),
		// force wrap + permanent
		ErrToRuntimePermanent(NewErrorRuntimePermanentFf("err"), "err", "err"),
		ErrToRuntimePermanent(fmt.Errorf("wrap: %w", NewErrorRuntimePermanentFf("err")), "err", "err"),
		// force wrap + tmp
		ErrToRuntimePermanent(NewErrorRuntimeTemporaryFf("err"), "err", "err"),
		ErrToRuntimePermanent(fmt.Errorf("wrap: %w", NewErrorRuntimeTemporaryFf("err")), "err", "err"),
	} {
		assert.False(t, err.IsTemporary(), "permanent case #%d", i)
		assert.True(t, err.IsPermanent(), "permanent case #%d", i)
	}

	// temporary
	for i, err := range []ErrorRuntime{
		// new
		NewErrorRuntimeTemporaryFf("err"),
		// default wrap + tmp
		ErrToRuntime(NewErrorRuntimeTemporaryFf("err"), "err", "err"),
		ErrToRuntime(fmt.Errorf("wrap: %w", NewErrorRuntimeTemporaryFf("err")), "err", "err"),
		// force wrap + default
		ErrToRuntimeTemporary(fmt.Errorf("err"), "err", "err"),
		// force wrap + permanent
		ErrToRuntimeTemporary(NewErrorRuntimePermanentFf("err"), "err", "err"),
		ErrToRuntimeTemporary(fmt.Errorf("wrap: %w", NewErrorRuntimePermanentFf("err")), "err", "err"),
		// force wrap + tmp
		ErrToRuntimeTemporary(NewErrorRuntimeTemporaryFf("err"), "err", "err"),
		ErrToRuntimeTemporary(fmt.Errorf("wrap: %w", NewErrorRuntimeTemporaryFf("err")), "err", "err"),
	} {
		assert.True(t, err.IsTemporary(), "temporary case #%d", i)
		assert.False(t, err.IsPermanent(), "temporary case #%d", i)
	}
}
