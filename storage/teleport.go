package storage

import (
	"context"
	"encoding/json"
)

// GetTeleportAuthServers returns auth servers.
func (s Storage) GetTeleportAuthServers(ctx context.Context) ([]string, error) {
	var servers []string
	serversJSON, err := s.get(ctx, KeyTeleportAuthServers)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(serversJSON), &servers)
	if err != nil {
		return nil, err
	}
	return servers, nil
}
