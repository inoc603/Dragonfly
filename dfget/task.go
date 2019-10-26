package dfget

import (
	"fmt"
	"time"

	"github.com/inoc603/Dragonfly/dfget/config"
	"github.com/inoc603/Dragonfly/pkg/rate"
)

type Task struct {
	// URL is the download url
	URL string `json:"url"`

	// LocalLimit rate limit about a single download task, format: G(B)/g/M(B)/m/K(B)/k/B
	// pure number will also be parsed as Byte.
	LocalLimit rate.Rate `json:"localLimit,omitempty"`

	// Minimal rate about a single download task, format: G(B)/g/M(B)/m/K(B)/k/B
	// pure number will also be parsed as Byte.
	MinRate rate.Rate `json:"minRate,omitempty"`

	// Timeout download timeout(second).
	Timeout time.Duration `json:"timeout,omitempty"`

	// Md5 expected file md5.
	Md5 string `json:"md5,omitempty"`

	// Identifier identify download task, it is available merely when md5 param not exist.
	Identifier string `json:"identifier,omitempty"`

	// CallSystem system name that executes dfget.
	CallSystem string `json:"callSystem,omitempty"`

	// Pattern download pattern, must be 'p2p' or 'cdn' or 'source',
	// default:`p2p`.
	Pattern string `json:"pattern,omitempty"`

	// Insecure indicates whether skip secure verify when supernode interact with the source.
	Insecure bool `json:"insecure,omitempty"`

	// CA certificate to verify when supernode interact with the source.
	Cacerts []string `json:"cacert,omitempty"`

	// Filter filter some query params of url, use char '&' to separate different params.
	// eg: -f 'key&sign' will filter 'key' and 'sign' query param.
	// in this way, different urls correspond one same download task that can use p2p mode.
	Filter []string `json:"filter,omitempty"`

	// Header of http request.
	// eg: --header='Accept: *' --header='Host: abc'.
	Header []string `json:"header,omitempty"`

	// Notbs indicates whether to not back source to download when p2p fails.
	Notbs bool `json:"notbs,omitempty"`

	// ClientQueueSize is the size of client queue
	// which controls the number of pieces that can be processed simultaneously.
	// It is only useful when the pattern not equals "source".
	// The default value is 6.
	ClientQueueSize int `json:"clientQueueSize,omitempty"`

	// Start time.
	StartTime time.Time `json:"-"`

	// Sign the value is 'Pid + float64(time.Now().UnixNano())/float64(time.Second) format: "%d-%.3f"'.
	// It is unique for downloading task, and is used for debugging.
	Sign string `json:"-"`

	// RV stores the variables that are initialized and used at downloading task executing.
	RV config.RuntimeVariable `json:"-"`

	// The reason of backing to source.
	BackSourceReason int `json:"-"`
}

type Range struct {
	Start uint
	End   uint
}

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.Start, r.End)
}

type Piece struct {
	URL string

	// Index is the position of the piece in all pieces from the task,
	// starting from 0.
	Index int

	// Range shows what part of the file this piece represents.
	Range Range

	MD5 string
}

func (p *Piece) Size() uint {
	return p.Range.End - p.Range.Start + 1
}
