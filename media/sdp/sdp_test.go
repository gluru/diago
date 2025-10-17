// SPDX-License-Identifier: MPL-2.0
// SPDX-FileCopyrightText: Copyright (c) 2024, Emir Aganovic

package sdp

import (
	"net"
	"testing"

	psdp "github.com/pion/sdp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromString_audioTrack(t *testing.T) {
	body := `v=0
o=Z 0 2640807 IN IP4 10.6.13.138
s=Z
c=IN IP4 10.6.13.138
t=0 0
m=audio 36600 RTP/AVP 0 101 8 3
a=rtpmap:101 telephone-event/8000
a=fmtp:101 0-16
a=sendrecv
a=rtcp-mux
`
	sd, err := FromString([]byte(body))
	assert.Nil(t, err)
	assert.Equal(t, sd.ConnectionInformation.AddressType, "IP4")
	assert.Equal(t, sd.ConnectionInformation.NetworkType, "IN")
	assert.Equal(t, sd.ConnectionInformation.Address.Address, "10.6.13.138")
	assert.Equal(t, sd.ConnectionInformation.Address.String(), "10.6.13.138")

	assert.Equal(t, 1, len(sd.MediaDescriptions))
	assert.Equal(t, 36600, sd.MediaDescriptions[0].MediaName.Port.Value)
	assert.Equal(t, "audio", sd.MediaDescriptions[0].MediaName.Media)
	assert.Equal(t, []string{"RTP", "AVP"}, sd.MediaDescriptions[0].MediaName.Protos)
	assert.Equal(t, []string{"0", "101", "8", "3"}, sd.MediaDescriptions[0].MediaName.Formats)

	attrs := []string{psdp.AttrKeySendRecv, psdp.AttrKeyRTCPMux}
	for _, attr := range attrs {
		_, ok := sd.MediaDescriptions[0].Attribute(attr)
		assert.True(t, ok)
	}
}

func TestFromString_body(t *testing.T) {
	body := `v=0
o=- 3905350750 3905350750 IN IP4 192.168.100.11
s=pjmedia
b=AS:84
t=0 0
a=X-nat:0
m=audio 57797 RTP/AVP 96 97 98 99 3 0 8 9 120 121 122
c=IN IP4 192.168.100.11
a=sendrecv
a=rtpmap:96 speex/16000
a=rtpmap:97 speex/8000
a=rtpmap:98 speex/32000
a=rtpmap:99 iLBC/8000
a=fmtp:99 mode=30
a=rtpmap:120 telephone-event/16000
a=fmtp:120 0-16
a=rtpmap:121 telephone-event/8000
a=fmtp:121 0-16
a=rtpmap:122 telephone-event/32000
a=fmtp:122 0-16
a=ssrc:1204560450 cname:4585300731f880ff
a=rtcp:57798 IN IP4 192.168.100.11
a=rtcp-mux
`
	sd, err := FromString([]byte(body))
	assert.Nil(t, err)
	assert.NotNil(t, sd)

	md := sd.MediaDescriptions[0]
	assert.Equal(t, "audio", sd.MediaDescriptions[0].MediaName.Media)

	require.NoError(t, err)
	require.Equal(t, 57797, md.MediaName.Port.Value)
	require.Equal(t, []string{"RTP", "AVP"}, md.MediaName.Protos)
	require.Equal(t, []string{"96", "97", "98", "99", "3", "0", "8", "9", "120", "121", "122"}, md.MediaName.Formats)

	ci := md.ConnectionInformation
	require.Equal(t, "IN", ci.NetworkType)
	require.Equal(t, "IP4", ci.AddressType)
	require.Equal(t, net.ParseIP("192.168.100.11").String(), ci.Address.String())
}
