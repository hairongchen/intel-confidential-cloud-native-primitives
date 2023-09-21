package quote_sdk

const (
	UDS_PATH = "unix:/run/ccnp/uds/quote-server.sock"
	TYPE_TDX = "TDX"
	TYPE_TPM = "TPM"
)

func GetQuote(user_data string, nonce string) (string, error) {

}
