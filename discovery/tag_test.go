package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func prepareTestcase(a *assert.Assertions, id TagID, dataLen int, data []byte) *Tag {
	a.Equal(dataLen+3, len(data))
	tag, err := ParseTag(id, uint16(dataLen), data[3:])

	a.Nil(err)
	a.NotNil(tag)
	a.Equal(id, tag.ID)
	return tag
}

func TestParseUptimeTag(t *testing.T) {
	assert := assert.New(t)

	tag := prepareTestcase(assert, tagUptime, 4, []byte{
		0x0a, 0x00, 0x04, // header+len
		0x01, 0x02, 0x03, 0x04, // uptime
	})

	val, ok := tag.value.(uint32)
	assert.True(ok)
	assert.Equal(uint32(16909060), val)
	assert.Equal("uptime", tag.Name())
}

func TestParseIPInfoTag(t *testing.T) {
	assert := assert.New(t)

	tag := prepareTestcase(assert, tagIPInfo, 6+4, []byte{
		0x02, 0x00, 0x0a, // header+len
		0x04, 0x18, 0xd6, 0x83, 0xf8, 0xec, // mac
		172, 16, 0, 1, // ip
	})

	val, ok := tag.value.(*ipInfo)
	assert.True(ok)
	assert.Equal("04:18:d6:83:f8:ec", val.MacAddress.String())
	assert.Equal("172.16.0.1", val.IPAddress.String())
}

func TestParseUnknownTag(t *testing.T) {
	assert := assert.New(t)

	tag := prepareTestcase(assert, TagID(0x42), 3, []byte{
		0x02, 0x00, 0x01, // header+len
		0xc0, 0xff, 0xee, // 0xC0FFEE
	})

	val, ok := tag.value.(string)
	assert.True(ok)
	assert.Equal("len:3<c0ffee>", val)
	assert.Equal("unknown", tag.Name())
	assert.Equal("unknown (0x42)", tag.Description())
}
