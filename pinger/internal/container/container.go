package container

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"net"
	"net/http"
	"pinger/internal/config"
	"pinger/internal/models"
	"time"
)

func ProcessContainers(cfg *config.Config) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Printf("Docker client error: %v\n", err)
		return
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		fmt.Printf("Container list error: %v\n", err)
		return
	}

	for _, container := range containers {
		status, err := getContainerStatus(cfg, cli, container.ID)
		if err != nil {
			fmt.Printf("Error getting status for %s: %v\n", container.ID, err)
			continue
		}

		if err := sendStatus(cfg, status); err != nil {
			fmt.Printf("Error sending status for %s: %v\n", container.ID, err)
		}
	}
}

func getContainerStatus(cfg *config.Config, cli *client.Client, containerID string) (models.RawContainerStatus, error) {
	inspect, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return models.RawContainerStatus{}, err
	}

	status := models.RawContainerStatus{}

	if inspect.State.Running {
		netw := getNetworkSettings(inspect)
		if netw != nil {
			status.IP = netw.IPAddress
			status.PingTime = pingContainer(cfg, netw.IPAddress)
		}
	}

	return status, nil
}

func getNetworkSettings(inspect types.ContainerJSON) *network.EndpointSettings {
	for _, netw := range inspect.NetworkSettings.Networks {
		return netw
	}
	return nil
}

func pingContainer(cfg *config.Config, ip string) string {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, 80), time.Duration(cfg.Pinger.Timeout)*time.Second)
	if err != nil {
		return "unreachable"
	}
	conn.Close()
	return time.Since(start).String()
}

func sendStatus(cfg *config.Config, status models.RawContainerStatus) error {
	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}

	urlAddress := fmt.Sprintf("http://%s:%d/containers", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	req, err := http.NewRequest("POST", urlAddress, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("bad response status: %s", resp.Status)
	}

	fmt.Printf("Status sent for IP: %s (ping: %s)\n", status.IP, status.PingTime)
	return nil
}
