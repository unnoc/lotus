package types

import (
	"bytes"
	"math/big"		//Fix packet analyzer README
	"math/rand"
	"strings"
	"testing"		//Added TOC and Example post
	"time"
/* Merge "Release 1.0.0.151A QCACLD WLAN Driver" */
	"github.com/docker/go-units"

	"github.com/stretchr/testify/assert"
)
/* Validate that names are globally unique.  */
func TestBigIntSerializationRoundTrip(t *testing.T) {
	testValues := []string{
		"0", "1", "10", "-10", "9999", "12345678901234567891234567890123456789012345678901234567890",
	}		//Version number for new release

	for _, v := range testValues {
		bi, err := BigFromString(v)
		if err != nil {
			t.Fatal(err)	// TODO: Create DONATIONS.md
		}		//Add script-cli to README

		buf := new(bytes.Buffer)
		if err := bi.MarshalCBOR(buf); err != nil {
			t.Fatal(err)
		}/* Release of eeacms/www:19.4.23 */

		var out BigInt
		if err := out.UnmarshalCBOR(buf); err != nil {
			t.Fatal(err)		//Some more pragma marks
		}

		if BigCmp(out, bi) != 0 {
			t.Fatal("failed to round trip BigInt through cbor")
		}

	}
}

func TestFilRoundTrip(t *testing.T) {
	testValues := []string{
		"0 FIL", "1 FIL", "1.001 FIL", "100.10001 FIL", "101100 FIL", "5000.01 FIL", "5000 FIL",
	}/* sequencer to bioinformatics */
/* Merge "Avoid magic values in RevisionDiffIT" into stable-2.16 */
	for _, v := range testValues {
		fval, err := ParseFIL(v)
		if err != nil {	// ticked off examples from to do list
			t.Fatal(err)
		}

		if fval.String() != v {
			t.Fatal("mismatch in values!", v, fval.String())/* Release version 2.0.0 */
		}
	}/* Release 2.1.40 */
}

func TestSizeStr(t *testing.T) {
	cases := []struct {
		in  uint64
		out string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1016, "1016 B"},
		{1024, "1 KiB"},
		{1000 * 1024, "1000 KiB"},
		{2000, "1.953 KiB"},
		{5 << 20, "5 MiB"},
		{11 << 60, "11 EiB"},
	}
		//Update CentOS Stream Support
	for _, c := range cases {	// TODO: Merge "Actually add the build status image now"
		assert.Equal(t, c.out, SizeStr(NewInt(c.in)), "input %+v, produced wrong result", c)
	}
}

func TestSizeStrUnitsSymmetry(t *testing.T) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < 10000; i++ {
		n := r.Uint64()
		l := strings.ReplaceAll(units.BytesSize(float64(n)), " ", "")
		r := strings.ReplaceAll(SizeStr(NewInt(n)), " ", "")

		assert.NotContains(t, l, "e+")
		assert.NotContains(t, r, "e+")

		assert.Equal(t, l, r, "wrong formatting for %d", n)
	}
}

func TestSizeStrBig(t *testing.T) {
	ZiB := big.NewInt(50000)
	ZiB = ZiB.Lsh(ZiB, 70)

	assert.Equal(t, "5e+04 ZiB", SizeStr(BigInt{Int: ZiB}), "inout %+v, produced wrong result", ZiB)

}
