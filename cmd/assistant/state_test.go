package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestCombinationPreset_IsEqualFalse1(t *testing.T) {
	cp := CombinationPreset{
		Codes: []uint16{5, 32, 234},
	}

	codes := make(KeysState)
	codes[4] = KeyPressed
	codes[5] = KeyPressed

	isEqual := cp.IsEqual(codes)
	assert.False(t, isEqual)
}

func TestCombinationPreset_IsEqualFalse2(t *testing.T) {
	cp := CombinationPreset{
		Codes: []uint16{5, 6},
	}

	codes := make(KeysState)
	codes[4] = KeyPressed
	codes[5] = KeyPressed

	isEqual := cp.IsEqual(codes)
	assert.False(t, isEqual)
}

func TestCombinationPreset_IsEqualTrue1(t *testing.T) {
	cp := CombinationPreset{
		Codes: []uint16{5, 4},
	}

	codes := make(KeysState)
	codes[4] = KeyPressed
	codes[5] = KeyPressed

	isEqual := cp.IsEqual(codes)
	assert.True(t, isEqual)
}
