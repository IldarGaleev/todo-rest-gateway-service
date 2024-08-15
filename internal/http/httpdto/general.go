package httpdto

type GeneralResponseStatus string

const (
	StatusOK    = GeneralResponseStatus("ok")
	StatusError = GeneralResponseStatus("error")
)

type GeneralResponse struct {
	Status GeneralResponseStatus `json:"status"`
}
