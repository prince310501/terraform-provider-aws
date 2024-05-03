// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kms

// Exports for use in other modules.
var (
	DiffSuppressKey             = diffSuppressKey
	FindDefaultKeyARNForService = findDefaultKeyARNForService
	ValidateKey                 = validateKey
	ValidateKeyOrAlias          = validateKeyOrAlias
)
