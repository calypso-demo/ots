package otssmc

import (
	"github.com/calypso-demo/ots"
	"gopkg.in/dedis/cothority.v1/skipchain"
	"gopkg.in/dedis/crypto.v0/abstract"
	"gopkg.in/dedis/onet.v1"
	"gopkg.in/dedis/onet.v1/network"
)

type Client struct {
	*onet.Client
}

func NewClient() *Client {
	return &Client{Client: onet.NewClient(ServiceName)}
}

func (c *Client) OTSDecrypt(r *onet.Roster, writeSBF *skipchain.SkipBlockFix, readSBF *skipchain.SkipBlockFix, inclusionProof *skipchain.BlockLink, acPubKeys []abstract.Point, privKey abstract.Scalar) ([]*DecryptedShare, onet.ClientError) {
	data := &DecryptRequestData{
		WriteSBF:       writeSBF,
		ReadSBF:        readSBF,
		InclusionProof: inclusionProof,
		ACPublicKeys:   acPubKeys,
	}
	msg, err := network.Marshal(data)
	if err != nil {
		return nil, onet.NewClientErrorCode(ErrorParse, err.Error())
	}
	sig, err := ots.SignMessage(msg, privKey)
	if err != nil {
		return nil, onet.NewClientErrorCode(ErrorParse, err.Error())
	}
	rootIdx := 0
	req := &DecryptRequest{
		RootIndex: rootIdx,
		Roster:    r,
		Data:      data,
		Signature: &sig,
	}
	reply := &DecryptReply{}
	err = c.SendProtobuf(r.List[req.RootIndex], req, reply)
	if err != nil {
		return nil, onet.NewClientErrorCode(ErrorParse, err.Error())
	}
	return reply.DecShares, nil
}
