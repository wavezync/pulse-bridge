package register

import "github.com/rs/zerolog/log"

func createResponse(resultChan ResultChanStruct) MonitorResponse {

	if resultChan.err != nil {
		log.Error().
			Err(resultChan.err).
			Str("monitor", resultChan.mntr.Name).
			Msg("Monitor check failed")
	} else {
		log.Info().
			Str("monitor", resultChan.mntr.Name).
			Msg("Monitor check successful")
	}

	return MonitorResponse{}
}
