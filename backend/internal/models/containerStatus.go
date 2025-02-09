package models

import "time"

type ContainerStatus struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`
	PingTime  time.Time `json:"ping_time"`
	LastCheck time.Time `json:"last_check"`
}

type RawContainerStatus struct {
	IP       string `json:"ip"`
	PingTime string `json:"ping_time"`
}

func (r *RawContainerStatus) ToContainerStatus() (*ContainerStatus, error) {
	pingTime, err := time.Parse(time.RFC3339, r.PingTime)
	if err != nil {
		return nil, err
	}

	return &ContainerStatus{
		IP:        r.IP,
		PingTime:  pingTime,
		LastCheck: time.Now(),
	}, nil
}
