package gocron_server

import (
	"os/exec"
	"strings"
)

type UnitExecutable interface {
	Call() (string, error)
}

type UnitExecCmd struct {
	name string
	args []string
}

func (ue *UnitExecCmd) Call() (string, error) {
	cmd := exec.Command(ue.name, ue.args...)
	out, err := cmd.Output()
	strout := string(out[:])

	if err != nil {
		return strout, err
	}

	return strout, nil
}

func NewUnitExecCmd(command string) *UnitExecCmd {
	cmdFrags := strings.Split(command, " ")
	return &UnitExecCmd{name: cmdFrags[0], args: cmdFrags[1:]}
}

type UnitExecFn struct {
	fn JobFunc
}

func (ue *UnitExecFn) Call() (string, error) {
	return ue.fn()
}

func NewUnitExecFn(fn JobFunc) *UnitExecFn {
	return &UnitExecFn{fn}
}

type SchedulerUnit struct {
	Name string
	Exec UnitExecutable
	Cron string
}

func NewSchedulerUnit(name, cron string, exec UnitExecutable) *SchedulerUnit {
	return &SchedulerUnit{
		Name: name,
		Exec: exec,
		Cron: cron,
	}
}

func NewManualUnit(name string, exec UnitExecutable) *SchedulerUnit {
	return &SchedulerUnit{
		Name: name,
		Exec: exec,
		Cron: "",
	}
}
