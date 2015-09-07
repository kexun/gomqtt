// Copyright (c) 2014 The gomqtt Authors. All rights reserved.
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

package message

import (
	"encoding/binary"
	"fmt"
)

// A single subscription in a SubscribeMessage.
type Subscription struct {
	Topic []byte
	QoS   byte
}

// The SUBSCRIBE Packet is sent from the Client to the Server to create one or more
// Subscriptions. Each Subscription registers a Client’s interest in one or more
// Topics. The Server sends PUBLISH Packets to the Client in order to forward
// Application Messages that were published to Topics that match these Subscriptions.
// The SUBSCRIBE Packet also specifies (for each Subscription) the maximum QoS with
// which the Server can send Application Messages to the Client.
type SubscribeMessage struct {
	header

	Subscriptions []Subscription
}

var _ Message = (*SubscribeMessage)(nil)

// NewSubscribeMessage creates a new SUBSCRIBE message.
func NewSubscribeMessage() *SubscribeMessage {
	msg := &SubscribeMessage{}
	msg.Type = SUBSCRIBE
	return msg
}

func (this SubscribeMessage) String() string {
	msgstr := fmt.Sprintf("%s, Packet ID=%d", this.header, this.PacketId)

	for i, t := range this.Subscriptions {
		msgstr = fmt.Sprintf("%s, Topic[%d]=%q/%d", msgstr, i, string(t.Topic), t.QoS)
	}

	return msgstr
}

func (this *SubscribeMessage) Len() int {
	ml := this.msglen()
	return this.header.len(ml) + ml
}

func (this *SubscribeMessage) Decode(src []byte) (int, error) {
	total := 0

	hl, _, rl, err := this.header.decode(src[total:])
	total += hl
	if err != nil {
		return total, err
	}

	this.PacketId = binary.BigEndian.Uint16(src[total:])
	total += 2

	remlen := int(rl) - (total - hl)
	for remlen > 0 {
		t, n, err := readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		this.Subscriptions = append(this.Subscriptions, Subscription{t, src[total]})
		total++

		remlen = remlen - n - 1
	}

	if len(this.Subscriptions) == 0 {
		return 0, fmt.Errorf(this.Name() + "/Decode: Empty subscription list")
	}

	return total, nil
}

func (this *SubscribeMessage) Encode(dst []byte) (int, error) {
	l := this.Len()

	if len(dst) < l {
		return 0, fmt.Errorf(this.Name()+"/Encode: Insufficient buffer size. Expecting %d, got %d.", l, len(dst))
	}

	total := 0

	n, err := this.header.encode(dst[total:], 0, this.msglen())
	total += n
	if err != nil {
		return total, err
	}

	binary.BigEndian.PutUint16(dst[total:], this.PacketId)
	total += 2

	for _, t := range this.Subscriptions {
		n, err := writeLPBytes(dst[total:], t.Topic)
		total += n
		if err != nil {
			return total, err
		}

		dst[total] = t.QoS
		total++
	}

	return total, nil
}

func (this *SubscribeMessage) msglen() int {
	// packet ID
	total := 2

	for _, t := range this.Subscriptions {
		total += 2 + len(t.Topic) + 1
	}

	return total
}
