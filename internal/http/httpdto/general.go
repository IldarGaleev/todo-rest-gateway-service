package httpdto

type GeneralResponseStatus string //@Name GeneralResponseStatus

const (
	StatusOK    = GeneralResponseStatus("ok")
	StatusError = GeneralResponseStatus("error")
)

type GeneralResponse struct {
	Status GeneralResponseStatus `json:"status"`
} //@Name GeneralResponse
