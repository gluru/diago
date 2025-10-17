// SPDX-License-Identifier: MPL-2.0
// SPDX-FileCopyrightText: Copyright (c) 2024, Emir Aganovic

package sdp

import (
	"net"
	"time"

	psdp "github.com/pion/sdp/v2"
)

func GetCurrentNTPTimestamp() uint64 {
	var ntpEpochOffset int64 = 2208988800 // Offset from Unix epoch (January 1, 1970) to NTP epoch (January 1, 1900)
	currentTime := time.Now().Unix() + int64(ntpEpochOffset)

	return uint64(currentTime)
}

func NTPTimestamp(now time.Time) uint64 {
	var ntpEpochOffset int64 = 2208988800 // Offset from Unix epoch (January 1, 1970) to NTP epoch (January 1, 1900)
	currentTime := now.Unix() + ntpEpochOffset

	return uint64(currentTime)
}

const (
	// https://datatracker.ietf.org/doc/html/rfc4566#section-6
	ModeRecvonly string = "recvonly"
	ModeSendrecv string = "sendrecv"
	ModeSendonly string = "sendonly"
)

func GenAudio(originIP net.IP, connectionIP net.IP, rtpPort int, mode string, fmts Formats) []byte {
	ntpTime := GetCurrentNTPTimestamp()

	sd := psdp.SessionDescription{
		SessionName: psdp.SessionName("Sip Go Media"),
		TimeDescriptions: []psdp.TimeDescription{
			{
				Timing: psdp.Timing{
					StartTime: 0,
					StopTime:  0,
				},
			},
		},
		Origin: psdp.Origin{
			Username:       "-",
			SessionID:      ntpTime,
			SessionVersion: ntpTime,

			AddressType:    "IP4",
			NetworkType:    "IN",
			UnicastAddress: originIP.String(),
		},
		ConnectionInformation: &psdp.ConnectionInformation{
			AddressType: "IP4",
			NetworkType: "IN",
			Address: &psdp.Address{
				Address: connectionIP.String(),
			},
		},
	}
	md := &psdp.MediaDescription{
		MediaName: psdp.MediaName{
			Media:   "audio",
			Port:    psdp.RangedPort{Value: rtpPort},
			Protos:  []string{"RTP", "AVP"},
			Formats: fmts,
		},
	}

	for _, f := range fmts {
		switch f {
		case FORMAT_TYPE_ULAW:
			md.Attributes = append(md.Attributes, psdp.NewAttribute("rtpmap", "0 PCMU/8000"))
			//formatsMap = append(formatsMap, "a=rtpmap:0 PCMU/8000")
		case FORMAT_TYPE_ALAW:
			md.Attributes = append(md.Attributes, psdp.NewAttribute("rtpmap", "8 PCMU/8000"))

			//formatsMap = append(formatsMap, "a=rtpmap:8 PCMA/8000")
		case FORMAT_TYPE_OPUS:
			md.Attributes = append(md.Attributes, psdp.NewAttribute("rtpmap", "96 opus/48000/2"))
			md.Attributes = append(md.Attributes, psdp.NewAttribute("fmtp", "96 useinbandfec=0"))

			//formatsMap = append(formatsMap, "a=rtpmap:96 opus/48000/2")
			//// Providing 0 when FEC cannot be used on the receiving side is RECOMMENDED.
			//// https://datatracker.ietf.org/doc/html/rfc7587
			//formatsMap = append(formatsMap, "a=fmtp:96 useinbandfec=0")
		case FORMAT_TYPE_TELEPHONE_EVENT:
			md.Attributes = append(md.Attributes, psdp.NewAttribute("rtpmap", "101 telephone-event/8000"))
			md.Attributes = append(md.Attributes, psdp.NewAttribute("fmtp", "101 0-16"))

			//formatsMap = append(formatsMap, "a=rtpmap:101 telephone-event/8000")
			//formatsMap = append(formatsMap, "a=fmtp:101 0-16")
		}
	}
	md.Attributes = append(md.Attributes, psdp.NewAttribute("ptime", "20"))
	md.Attributes = append(md.Attributes, psdp.NewAttribute("maxptime", "20"))
	md.Attributes = append(md.Attributes, psdp.NewAttribute(mode, ""))

	sd.MediaDescriptions = []*psdp.MediaDescription{md}

	bytes, _ := sd.Marshal()
	return bytes
}
