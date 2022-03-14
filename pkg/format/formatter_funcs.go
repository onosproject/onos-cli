// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package format

import (
	"time"

	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
)

// formats a Timestamp proto as a RFC3339 date string
func formatTimestamp(tsproto *timestamppb.Timestamp) (string, error) {
	if tsproto == nil {
		return "", nil
	}
	return tsproto.AsTime().Truncate(time.Second).Format(time.RFC3339), nil
}

// Computes the age of a timestamp and returns it in HMS format
func formatGoSince(ts time.Time) (string, error) {
	return time.Since(ts).Truncate(time.Second).String(), nil
}

// Computes the age of a timestamp and returns it in HMS format
func formatSince(tsproto *timestamppb.Timestamp) (string, error) {
	if tsproto == nil {
		return "", nil
	}
	return time.Since(tsproto.AsTime()).Truncate(time.Second).String(), nil
}
