package UsersBackend

type ResponseBack struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}
