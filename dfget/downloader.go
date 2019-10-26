package dfget

import (
	"context"
	"io"
)

type Downloader interface {
	Download(ctx context.Context, piece *Piece) (io.ReadCloser, error)
}
