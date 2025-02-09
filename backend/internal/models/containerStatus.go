package models

import "time"

type ContainerStatus struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`
	PingTime  string    `json:"ping_time"`
	LastCheck time.Time `json:"last_check"`
}

type RawContainerStatus struct {
	IP       string `json:"ip"`
	PingTime string `json:"ping_time"`
}

func (r *RawContainerStatus) ToContainerStatus() (*ContainerStatus, error) {

	return &ContainerStatus{
		IP:        r.IP,
		PingTime:  r.PingTime,
		LastCheck: time.Now(),
	}, nil
}
