package quote

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"log"
	"strings"
	"time"

	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/quote/proto"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	UDS_PATH = "unix:/run/ccnp/uds/quote-server.sock"
	TYPE_TDX = "TDX"
	TYPE_TPM = "TPM"
)

func GetQuote(user_data string, nonce string) (interface{}, error) {

	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[quote SDK] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewGetQuoteClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetQuote(ctx, &pb.GetQuoteRequest{UserData: user_data, Nonce: nonce})
	if err != nil {
		log.Fatalf("[quote SDK] fail to get quote: %v", err)
	}

	quote, err := base64.StdEncoding.DecodeString(strings.Trim(response.Quote, "\""))
	if err != nil {
		log.Fatalf("[quote SDK] decode quote error: %v", err)
	}

	switch response.QuoteType {
	case TYPE_TDX:
		return parseTDXQuote(quote)
	case TYPE_TPM:
		return parseTPMQuote(quote)
	default:
		log.Fatalf("[quote SDK] unknown TEE enviroment!")
	}

	return nil, pkgerrors.New("[quote SDK] unknown TEE enviroment!")
}

type TDXQuote struct {
	Quote           []uint8
	Version         uint16
	Tdreport        [584]uint8
	Tee_type        uint32
	Tee_tcb_svn     [16]uint8
	Mrseam          [48]uint8
	Mrsignerseam    [48]uint8
	Seamattributes  [8]uint8
	Tdattributes    [8]uint8
	Xfam            [8]uint8
	Mrtd            [48]uint8
	Mrconfigid      [48]uint8
	Mrowner         [48]uint8
	Mrownerconfig   [48]uint8
	Rtmrs           [192]uint8
	Reportdata      [64]uint8
	Signature       [64]uint8
	Attestation_key [64]uint8
	Cert_data       []uint8
}

type SGX_Quote_Header struct {
	Version      uint16 ///< 0:  The version this quote structure.
	Att_key_type uint16 ///< 2:  sgx_attestation_algorithm_id_t.  Describes the type of signature in the signature_data[] field.
	Tee_type     uint32 ///< 4:  Type of Trusted Execution Environment for which the Quote has been generated.
	///      Supported values: 0 (SGX), 0x81(TDX)
	Reserved  uint32    ///< 8:  Reserved field.
	Vendor_id [16]uint8 ///< 12: Unique identifier of QE Vendor.
	User_data [20]uint8 ///< 28: Custom attestation key owner data.
}

type TDReport struct {
	Tee_tcb_svn    [16]uint8
	Mrseam         [48]uint8
	Mrsignerseam   [48]uint8
	Seamattributes [8]uint8
	Tdattributes   [8]uint8
	Xfam           [8]uint8
	Mrtd           [48]uint8
	Mrconfigid     [48]uint8
	Mrowner        [48]uint8
	Mrownerconfig  [48]uint8
	Rtmrs          [192]uint8
	Reportdata     [64]uint8
}

const (
	Quote_Header_Index                    = 0   // 48 bytes quote header, start from index 0 of quote string
	Quote_TDReport_Index                  = 48  // 584 bytes tdreport, start from index 48 of quote string
	Quote_Auth_data_size_Index            = 632 // 4 bytes auth size, start from index 632 of quote string
	Quote_Auth_data_content_Index         = 636 // auth_size byets of auth_data, start from index 636 of quote string
	Quote_Auth_data_signature_Index       = 700
	Quote_Auth_data_attestation_key_Index = 764
	Quote_Auth_data_cert_data__Index      = 770
)

func parseTDXQuote(quote []byte) (interface{}, error) {

	// https://github.com/intel/SGXDataCenterAttestationPrimitives/blob/6882afad8644c27db162b40994402c8ad2a7fb32/QuoteGeneration/quote_wrapper/common/inc/sgx_quote_4.h#L141

	var header = SGX_Quote_Header{}
	var err = binary.Read(bytes.NewReader(quote[Quote_Header_Index:Quote_TDReport_Index]), binary.LittleEndian, &header)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote header: %v", err)
	}

	var tdreport = TDReport{}
	err = binary.Read(bytes.NewReader(quote[Quote_TDReport_Index:Quote_Auth_data_size_Index]), binary.LittleEndian, &tdreport)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote tdreport: %v", err)
	}

	var auth_size uint32 = 0
	err = binary.Read(bytes.NewReader(quote[Quote_Auth_data_size_Index:Quote_Auth_data_content_Index]), binary.LittleEndian, &auth_size)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote auth data size: %v", err)
	}

	var signature = [64]uint8{}
	err = binary.Read(bytes.NewReader(quote[Quote_Auth_data_content_Index:Quote_Auth_data_signature_Index]), binary.LittleEndian, &signature)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote signature: %v", err)
	}

	var attestation_key = [64]uint8{}
	err = binary.Read(bytes.NewReader(quote[Quote_Auth_data_signature_Index:Quote_Auth_data_attestation_key_Index]), binary.LittleEndian, &attestation_key)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote attestation_key: %v", err)
	}

	var cert_data = make([]uint8, auth_size-128-6)
	err = binary.Read(bytes.NewReader(quote[Quote_Auth_data_cert_data__Index:Quote_Auth_data_cert_data__Index+auth_size-6-128]), binary.LittleEndian, &cert_data)
	if err != nil {
		log.Fatalf("[parseTDXQuote] fail to parse quote cert_data: %v", err)
	}

	var tdquote = TDXQuote{}
	var quote_len = len(quote)
	tdquote.Quote = make([]byte, quote_len)
	tdquote.Quote = quote
	tdquote.Version = header.Version
	copy(tdquote.Tdreport[:], quote[Quote_TDReport_Index:Quote_Auth_data_size_Index])
	tdquote.Tee_type = header.Tee_type
	tdquote.Tee_tcb_svn = tdreport.Tee_tcb_svn
	tdquote.Mrseam = tdreport.Mrseam
	tdquote.Mrsignerseam = tdreport.Mrsignerseam
	tdquote.Seamattributes = tdreport.Seamattributes
	tdquote.Tdattributes = tdreport.Tdattributes
	tdquote.Xfam = tdreport.Xfam
	tdquote.Mrtd = tdreport.Mrtd
	tdquote.Mrconfigid = tdreport.Mrconfigid
	tdquote.Mrowner = tdreport.Mrowner
	tdquote.Mrownerconfig = tdreport.Mrownerconfig
	tdquote.Rtmrs = tdreport.Rtmrs
	tdquote.Reportdata = tdreport.Reportdata
	tdquote.Signature = signature
	tdquote.Attestation_key = attestation_key
	tdquote.Cert_data = make([]byte, auth_size)
	tdquote.Cert_data = cert_data

	return tdquote, nil
}

type TPMQuote struct{}

func parseTPMQuote(quote []byte) (interface{}, error) {
	// TODO: add vTPM support later
	return nil, pkgerrors.New("TPM support to be implemented later.")
}
