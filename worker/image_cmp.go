package worker

import (
	"context"

	"github.com/cybozu-go/neco"
)

func (o *operator) needContainerImageUpdate(ctx context.Context, name string) (bool, error) {
	tag, err := o.storage.GetContainerTag(ctx, o.mylrn, name)
	if err != nil {
		return false, err
	}

	img, err := neco.CurrentArtifacts.FindContainerImage(name)
	if err != nil {
		return false, err
	}

	return img.Tag != tag, nil
}
