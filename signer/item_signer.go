package signer

import (
	"crypto/sha256"
	"errors"
	"io"

	"github.com/everFinance/goar/utils"
	"github.com/everFinance/goether"
	Data "github.com/liteseed/goar/tx"
	"github.com/liteseed/goar/types"
)

type ItemSigner struct {
	signType   int
	signer     interface{}
	owner      string // only rsa has owner
	signerAddr string
}

func NewItemSigner(signer interface{}) (*ItemSigner, error) {
	signType, signerAddr, owner, err := reflectSigner(signer)
	if err != nil {
		return nil, err
	}
	return &ItemSigner{
		signType:   signType,
		signer:     signer,
		owner:      owner,
		signerAddr: signerAddr,
	}, nil
}

func (i *ItemSigner) CreateAndSignItem(data []byte, target string, anchor string, tags []types.Tag) (types.DataItem, error) {
	bundleItem, err := Data.NewBundleItem(i.owner, i.signType, target, anchor, data, tags)
	if err != nil {
		return types.DataItem{}, err
	}
	// sign
	if err := SignBundleItem(i.signType, i.signer, bundleItem); err != nil {
		return types.DataItem{}, err
	}
	// get itemBinary
	itemBinary, err := Data.GenerateItemBinary(bundleItem)
	if err != nil {
		return types.DataItem{}, err
	}
	bundleItem.RawData = itemBinary
	return *bundleItem, nil
}

func (i *ItemSigner) CreateAndSignNestedItem(target string, anchor string, tags []types.Tag, items ...types.DataItem) (types.DataItem, error) {
	bundleTags := []types.Tag{
		{Name: "Bundle-Format", Value: "binary"},
		{Name: "Bundle-Version", Value: "2.0.0"},
	}
	tags = append(tags, bundleTags...)

	bundle, err := Data.NewBundle(items...)
	if err != nil {
		return types.DataItem{}, err
	}
	return i.CreateAndSignItem(bundle.BundleBinary, target, anchor, tags)
}

func (i *ItemSigner) CreateAndSignItemStream(data io.Reader, target string, anchor string, tags []types.Tag) (types.DataItem, error) {
	bundleItem, err := Data.NewBundleItemStream(i.owner, i.signType, target, anchor, data, tags)
	if err != nil {
		return types.DataItem{}, err
	}
	// sign
	if err := SignBundleItem(i.signType, i.signer, bundleItem); err != nil {
		return types.DataItem{}, err
	}
	if _, err := bundleItem.DataReader.Seek(0, 0); err != nil {
		return types.DataItem{}, err
	}
	return *bundleItem, nil
}

func reflectSigner(signer interface{}) (signType int, signerAddr, owner string, err error) {
	if s, ok := signer.(*Signer); ok {
		signType = types.ArweaveSignType
		signerAddr = s.Address
		owner = s.Owner()
		return
	}
	if s, ok := signer.(*goether.Signer); ok {
		signType = types.EthereumSignType
		signerAddr = s.Address.String()
		owner = utils.Base64Encode(s.GetPublicKey())
		return
	}
	err = errors.New("not support this signer")
	return
}

func SignBundleItem(signatureType int, signer interface{}, item *types.DataItem) error {
	signMsg, err := Data.BundleItemSignData(*item)
	if err != nil {
		return err
	}
	var sigData []byte
	switch signatureType {
	case types.ArweaveSignType:
		arSigner, ok := signer.(*Signer)
		if !ok {
			return errors.New("signer must be goar signer")
		}
		sigData, err = utils.Sign(signMsg, arSigner.PrvKey)
		if err != nil {
			return err
		}

	case types.EthereumSignType:
		ethSigner, ok := signer.(*goether.Signer)
		if !ok {
			return errors.New("signer not goether signer")
		}
		sigData, err = ethSigner.SignMsg(signMsg)
		if err != nil {
			return err
		}
	default:
		// todo come soon supprot ed25519
		return errors.New("not supprot this signType")
	}
	id := sha256.Sum256(sigData)
	item.Id = utils.Base64Encode(id[:])
	item.Signature = utils.Base64Encode(sigData)
	return nil
}
