package ots

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/edouardparis/ancre/operation"
)

const hexDump = `004f70656e54696d657374616d7073000050726f6f6600bf89e2e884e89294010803ba204e50d126e4674c005e04d82e84c21366780af1f43bd54a37816b6ab34003f1c8010100000001e482f9d32ecc3ba657b69d898010857b54457a90497982ff56f97c4ec58e6f98010000006b483045022100b253add1d1cf90844338a475a04ff13fc9e7bd242b07762dea07f5608b2de367022000b268ca9c3342b3769cdd062891317cdcef87aac310b6855e9d93898ebbe8ec0121020d8e4d107d2b339b0050efdd4b4a09245aa056048f125396374ea6a2ab0709c6ffffffff026533e605000000001976a9140bf057d40fbba6744862515f5b55a2310de5772f88aca0860100000000001976a914f00688ac000000000808f120a987f716c533913c314c78e35d35884cac943fa42cac49d2b2c69f4003f85f880808f120dec55b3487e1e3f722a49b55a7783215862785f4a3acb392846019f71dc64a9d0808f120b2ca18f485e080478e025dab3d464b416c0e1ecb6629c9aefce8c8214d0424320808f02011b0e90661196ff4b0813c3eda141bab5e91604837bdf7a0c9df37db0e3a11980808f020c34bc1a4a1093ffd148c016b1e664742914e939efabe4d3d356515914b26d9e20808f020c3e6e7c38c69f6af24c2be34ebac48257ede61ec0a21b9535e4443277be306460808f1200798bf8606e00024e5d5d54bf0c960f629dfb9dad69157455b6f2652c0e8de810808f0203f9ada6d60baa244006bb0aad51448ad2fafb9d4b6487a0999cff26b91f0f5360808f120c703019e959a8dd3faef7489bb328ba485574758e7091f01464eb65872c975c80808f020cbfefff513ff84b915e3fed6f9d799676630f8364ea2a6c7557fad94a5b5d7880808f1200be23709859913babd4460bbddf8ed213e7c8773a4b1face30f8acfdf093b7050808000588960d73d7190103f7ef15`

func TestTimestampFileFromReader(t *testing.T) {
	is := assert.New(t)

	b, err := hex.DecodeString(hexDump)
	is.NoError(err)

	r := bytes.NewReader(b)
	file, err := TimestampFileFromReader(r)
	is.NoError(err)

	is.True(file.DigestType.Match(operation.Sha256))

	firstStep := file.Timestamp.FirstStep
	is.True(firstStep.Match(operation.Ripemd160))

	is.True(file.Timestamp.FirstStep.HasNext())
}

func TestTimestampFileToWriter(t *testing.T) {
	is := assert.New(t)

	expected, err := hex.DecodeString(hexDump)
	is.NoError(err)

	r := bytes.NewReader(expected)
	file, err := TimestampFileFromReader(r)
	is.NoError(err)

	var buf bytes.Buffer
	is.NoError(TimestampFileToWriter(file, &buf))

	actual := buf.Bytes()
	for i := range expected {
		is.Equal(expected[i], actual[i])
	}
}
