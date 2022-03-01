package api

import (
	"fmt"
)

// type KVVersion struct {
// 	Version        int
// 	CreatedTime    time.Time
// 	DeletionTime   time.Time
// 	Destroyed      bool
// 	CustomMetadata map[string]interface{}
// }

// type KVMetadata struct {
// 	CasRequired    bool
// 	CurrentTime    time.Time
// 	CustomMetadata map[string]interface{}

// 	Versions map[int]KVVersion
// }

type KV struct {
	c          *Client
	MountPoint string
}

type KVSecret struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
	Raw      *Secret
}

func (c *Client) KV() *KV {
	return c.KVWithMountPoint("secret")
}

func (c *Client) KVWithMountPoint(mountPoint string) *KV {
	// todo: validate mountpoint
	return &KV{
		c:          c,
		MountPoint: mountPoint,
	}
}

func (c *KV) Put(path string, data map[string]interface{}) (*KVSecret, error) {
	// todo: this is POC only, we should probably use lower level func here
	secret, err := c.c.Logical().Write(
		fmt.Sprintf("%s/data/%s", c.MountPoint, path),
		map[string]interface{}{
			"data": data,
		},
	)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	return &KVSecret{
		Data:     nil,
		Metadata: secret.Data,
		Raw:      secret,
	}, nil
}

func (c *KV) Get(path string) (*KVSecret, error) {
	// todo: this is POC only, we should probably use lower level func here
	secret, err := c.c.Logical().Read(fmt.Sprintf("%s/data/%s", c.MountPoint, path))
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, nil
	}

	return extractDataAndMetadata(secret)
}

func extractDataAndMetadata(secret *Secret) (*KVSecret, error) {
	data, ok := secret.Data["data"]
	if !ok {
		return nil, fmt.Errorf("missing expected 'data' element")
	}

	dataTyped, ok := data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for 'data' element: %T (%#v)", data, data)
	}

	metadata, ok := secret.Data["metadata"]
	if !ok {
		return nil, fmt.Errorf("missing expected 'metadata' element")
	}

	metadataTyped, ok := metadata.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for 'metadata' element: %T (%#v)", metadata, metadata)
	}

	return &KVSecret{
		Data:     dataTyped,
		Metadata: metadataTyped,
		Raw:      secret,
	}, nil
}
