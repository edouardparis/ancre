package attestation


type Attestation interface {
	Match(int) bool
	Input() []byte
	Data() map[string]interface{}
	Exec([]byte) []byte
}
