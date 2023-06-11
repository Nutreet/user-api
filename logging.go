package main

import (
	"context"
	"time"

	"firebase.google.com/go/v4/auth"
	proto "github.com/nutreet/common/gen/user"
	log "github.com/sirupsen/logrus"
)

type logger struct {
	next UserService
}

func NewLogger(s UserService) UserService {
	return &logger{
		next: s,
	}
}

func getLogFields(begin time.Time, err error, baseFields log.Fields) log.Fields {
	fields := log.Fields{
		"service": "User",
		"took":    time.Since(begin),
	}

	if err != nil {
		fields["error"] = err
	} else {
		for key, value := range baseFields {
			fields[key] = value
		}
	}

	return fields
}

func (l *logger) Register(ctx context.Context, data *proto.RegisterRequest) (u *auth.UserRecord, err error) {
	defer func(begin time.Time) {
		fields := log.Fields{
			"user": u,
		}
		entry := log.WithFields(getLogFields(begin, err, fields))

		if err != nil {
			entry.Error("Register")
		} else {
			entry.Info("Register")
		}

	}(time.Now())

	return l.next.Register(ctx, data)
}

func (l *logger) GetAutenticatedUser(ctx context.Context, token string) (uid string, err error) {
	defer func(begin time.Time) {
		fields := log.Fields{
			"uid": uid,
		}
		entry := log.WithFields(getLogFields(begin, err, fields))

		if err != nil {
			entry.Error("GetAutenticatedUser")
		} else {
			entry.Info("GetAutenticatedUser")
		}

	}(time.Now())

	return l.next.GetAutenticatedUser(ctx, token)
}
