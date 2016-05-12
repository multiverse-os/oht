package network

import "testing"
import "strconv"

func TestValidate(t *testing.T) {
	msg := new(Message)

	// Empty Message
	err := msg.validate()
	if err == nil {
		t.Log("Failed at Empty Message Invalidation.")
		t.Fail()
	}

	// Valid Message
	msg.magic = knownMagic
	msg.payload = []byte{'a', 'b', 'c', 'd'}
	msg.length = 4
	msg.command = "version"
	msg.checksum = 0xd8022f20

	err = msg.validate()
	if err != nil {
		t.Log("Failed at Validating Message: ", err)
		t.Fail()
	}

	// Invalid Magic Number
	msg.magic += 1
	err = msg.validate()
	if err != MessageError(EMAGIC) {
		t.Log("Failed at invalidating magic number: ", err)
		t.Fail()
	}
	msg.magic -= 1

	// Invalid Checksum
	msg.checksum += 1
	err = msg.validate()
	if err != MessageError(ECHECK) {
		t.Log("Failed at invalidating checksum: ", err)
		t.Fail()
	}
	msg.checksum -= 1

	// Invalid Length
	msg.length += 1
	err = msg.validate()
	if err != MessageError(EPALEN) {
		t.Log("Failed at invalidating length: ", err)
		t.Fail()
	}
}

type Payload []byte

func (p Payload) Serialize() []byte {
	return []byte(p)
}

func TestMakeMessage(t *testing.T) {

	str := "version"

	msg := MakeMessage(str, Payload([]byte{'a', 'b', 'c', 'd'}), nil)

	// Basic Validation Test
	err := msg.validate()
	if err != nil {
		t.Log("Created invalid message: ", err)
		t.Fail()
	}

	// Check that values are correct
	str = "notversion"
	if msg.command != strconv.Quote("version") {
		t.Log("Command string doesn't make defensive copy: ", msg.command)
		t.Fail()
	}

	if msg.length != 4 {
		t.Log("Incorrect payload length.")
		t.Fail()
	}

}

func TestSerialize(t *testing.T) {
	msg := new(Message)
	msg2 := new(Message)
	var serial []byte
	test := true

	// Valid Message
	msg.magic = knownMagic
	msg.payload = []byte{'a', 'b', 'c', 'd'}
	msg.length = 4
	msg.command = "version"
	msg.checksum = 0xd8022f20

	serial = msg.Serialize()

	// Check final length
	if len(serial) != 28 { // 24 Header plus 4 Payload
		t.Log("Bad serialized slice length: ", len(serial))
		t.Fail()
	}

	// Check Magic Number
	test = (serial[0] == 0xe9) && (serial[1] == 0xbe) && (serial[2] == 0xb4) && (serial[3] == 0xd9)

	if !test {
		t.Log("Incorrect magic number")
		t.Fail()
		test = true
	}

	// Check command
	cmd := serial[4:16]
	cmd2 := make([]byte, 12, 12)
	copy(cmd2, msg.command)

	for i := 0; i < 12; i++ {
		if cmd[i] != cmd2[i] {
			t.Log("Incorrect command")
			t.Fail()
		}
	}

	// Check length
	for i := 16; i < 19; i++ {
		if serial[i] != 0 {
			t.Log("Bad length, too big")
			t.Fail()
		}
	}
	if serial[19] != 4 {
		t.Log("Bad length: ", serial[19])
		t.Fail()
	}

	// Check checksum
	test = (serial[20] == 0xd8) && (serial[21] == 0x02) && (serial[22] == 0x2f) && (serial[23] == 0x20)

	if !test {
		t.Log("Incorrect checksum")
		t.Fail()
		test = true
	}

	// Check payload
	test = (serial[24] == 'a') && (serial[25] == 'b') && (serial[26] == 'c') && (serial[27] == 'd')

	if !test {
		t.Log("Incorrect payload")
		t.Fail()
		test = true
	}

	// Check makeHeader
	msg2.makeHeader(serial[:24])

	if msg.magic != msg2.magic || msg.length != msg2.length || msg.checksum != msg2.checksum {
		t.Log("makeHeader() failed to return same message")
		t.Fail()
	}
}
