package main

import (
	"time"

	"github.com/coderero/toll_calculator/types"
	"github.com/sirupsen/logrus"
)

type LoggingMiddleware struct {
	next DataProducer
}

func (lm *LoggingMiddleware) Produce(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obu_id": data.OBUID,
			"lat":    data.Lat,
			"long":   data.Long,
			"took":   time.Since(start),
		}).Info("Produced data")
	}(time.Now())
	return lm.next.Produce(data)
}
