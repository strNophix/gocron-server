package main

import (
	gocron_server "github.com/strnophix/gocron-server/pkg"
)

type ServerConfig struct {
	Host string `toml:"host"`
}

type UnitConfig struct {
	Name    string `toml:"name"`
	Command string `toml:"command"`
	Cron    string `toml:"cron"`
}

func (u *UnitConfig) ToSchedulerUnit() *gocron_server.SchedulerUnit {
	cmd := gocron_server.NewUnitExecCmd(u.Command)
	unit := gocron_server.NewSchedulerUnit(u.Name, u.Cron, cmd)
	return unit
}

type Config struct {
	Server ServerConfig `toml:"server"`
	Unit   []UnitConfig `toml:"unit"`
}
