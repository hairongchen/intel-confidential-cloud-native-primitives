package quote

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/quote/proto"
)

const (
	UDS_PATH = "unix:/run/ccnp/uds/quote-server.sock"
	TYPE_TDX = "TDX"
	TYPE_TPM = "TPM"
)

func GetQuote(user_data string, nonce string) (string, error) {

	conn, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGetQuoteClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.GetQuote(ctx, &pb.GetQuoteRequest{UserData: "", Nonce: ""})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Response quote type: %s", r.QuoteType)
	log.Printf("Response quote: %s", r.Quote)

	return r.QuoteType, nil
}
