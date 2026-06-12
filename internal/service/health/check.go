package health

import (
	"context"

	healthDTO "github.com/huypham67/user-service/internal/dto/health"
	"github.com/rs/zerolog/log"
)

// GetStatus checks the health status of the application by pinging the database and returns a HealthCheckResponse.
func (s *service) GetStatus(ctx context.Context) healthDTO.HealthCheckResponse {
	if err := s.pinger.Ping(ctx); err != nil {
		log.Error().
			Err(err).
			Str("service", s.serviceName).
			Msg("Database connection failed - health check")

		return healthDTO.HealthCheckResponse{
			Message:     failedStatusMessage,
			ServiceName: s.serviceName,
			InstanceID:  s.instanceID,
		}
	}

	return healthDTO.HealthCheckResponse{
		Message:     statusMessage,
		ServiceName: s.serviceName,
		InstanceID:  s.instanceID,
	}
}
