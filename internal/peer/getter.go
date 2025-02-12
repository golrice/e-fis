package peer

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	pb "github.com/golrice/e-fis/internal/protocal"
)

type HttpGetter struct {
	BaseURL string
}

func (h *HttpGetter) Get(in *pb.Request, out *pb.Response) error {
	url := fmt.Sprintf("%v%v/%v", h.BaseURL, url.QueryEscape(in.NodeName), url.QueryEscape(in.Key))

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server return: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	out.Value = bytes

	return nil
}

// make sure httpgetter is peergetter
var _ PeerGetter = (*HttpGetter)(nil)
