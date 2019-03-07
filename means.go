package updater

import (
	"context"

	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

type MeansType string

// Means manages methods related to specified update means.
// For example, fetches the latest tag, update binary, or
// check whether the software is installed by this.
type Means interface {
	LatestTag(context.Context) (*version.Version, error)
	Update(context.Context, *version.Version) error
	Installed(context.Context) bool

	CommandText(*version.Version) string

	Type() MeansType
}

type MeansBuilder func() (Means, error)

func SelectAvailableMeansFrom(ctx context.Context, ma ...MeansBuilder) (Means, error) {
	for i := range ma {
		m, err := ma[i]()
		// if the means unavailable, ignore it
		// but other errors found, abort selection and return its err
		if err == ErrUnavailable {
			continue
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to instantiate means")
		}

		// found
		if m.Installed(ctx) {
			return m, nil
		}
	}
	// maybe manually installed (like go get)
	return nil, ErrUnavailable
}

func NewMeans(b MeansBuilder) (Means, error) {
	return b()
}
