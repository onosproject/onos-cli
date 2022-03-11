// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

// None returns the string <None> if the input is empty, otherwise it returns the string
func None(s string) string {
	if s == "" {
		return "<None>"
	}
	return s
}
