// SPDX-License-Identifier: MPL-2.0
// SPDX-FileCopyrightText: Copyright (c) 2024, Emir Aganovic

package sdp

import (
	psdp "github.com/pion/sdp/v3"
)

func FromString(body []byte) (*psdp.SessionDescription, error) {
	p := psdp.SessionDescription{}
	if err := p.Unmarshal(body); err != nil {
		return nil, err
	}
	return &p, nil
}
