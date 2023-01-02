package worker

// This struct will contain the data of the task that we want to store in redis
type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}
