package dfget

import "time"

type Result struct {
	Error    error
	Duration time.Duration
}

type Scheduler interface {
	// GetPieces split the given task into pieces.
	GetPieces(task *Task) ([]*Piece, error)
	// FinishPiece reports the download result of a piece
	FinishPiece(p *Piece, result *Result) error
}
