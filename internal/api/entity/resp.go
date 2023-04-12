package entity

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data T      `json:"data,omitempty"`
}

func NewOkResp[T any](data T, msg *string) Resp[T] {
	r := Resp[T]{
		Code: 0,
		Data: data,
	}
	if msg != nil {
		r.Msg = *msg
	}
	return r
}

func NewErrResp[T any](data T, msg *string) Resp[T] {
	r := Resp[T]{
		Code: 1,
		Data: data,
	}
	if msg != nil {
		r.Msg = *msg
	}
	return r
}
