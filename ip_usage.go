package main

import (
	"net"
	"net/http"
	"time"

	"github.com/kaspanet/faucet/database"
	"github.com/kaspanet/kasparov/httpserverutils"
	"github.com/pkg/errors"
)

const minRequestInterval = time.Hour * 24

type ipUse struct {
	IP      string
	LastUse time.Time
}

func ipFromRequest(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func validateIPUsage(r *http.Request) error {
	db, err := database.DB()
	if err != nil {
		return err
	}
	now := time.Now()
	timeBeforeMinRequestInterval := now.Add(-minRequestInterval)
	ip, err := ipFromRequest(r)
	if err != nil {
		return err
	}
	count, err := db.Model(&ipUse{}).
		Where("ip = ?", ip).
		Where("last_use BETWEEN ? AND ?", timeBeforeMinRequestInterval, now).
		Count()

	if count != 0 {
		return httpserverutils.NewHandlerError(http.StatusForbidden, errors.New("A user is allowed to to have one request from the faucet every 24 hours"))
	}
	return nil
}

func updateIPUsage(r *http.Request) error {
	db, err := database.DB()
	if err != nil {
		return err
	}

	ip, err := ipFromRequest(r)
	if err != nil {
		return err
	}

	ipUse := &ipUse{IP: ip}
	_, err = db.Model(ipUse).
		OnConflict("(ip) DO UPDATE").
		Set("last_use = ?", time.Now()).
		Insert()

	if err != nil {
		return err
	}
	return nil
}
