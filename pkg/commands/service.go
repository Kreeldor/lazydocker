package commands

import (
	"os/exec"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/fatih/color"
	"github.com/jesseduffield/lazydocker/pkg/utils"
	"github.com/sirupsen/logrus"
)

// Service : A docker Service
type Service struct {
	Name      string
	ID        string
	OSCommand *OSCommand
	Log       *logrus.Entry
	Container *Container
}

// GetDisplayStrings returns the dispaly string of Container
func (s *Service) GetDisplayStrings(isFocused bool) []string {

	if s.Container == nil {
		return []string{utils.ColoredString("none", color.FgBlack), utils.ColoredString(s.Name, color.FgWhite), ""}
	}

	cont := s.Container
	return []string{utils.ColoredString(cont.Container.State, cont.GetColor()), utils.ColoredString(s.Name, color.FgWhite), cont.GetDisplayCPUPerc()}
}

// Remove removes the service's containers
func (s *Service) Remove(options types.ContainerRemoveOptions) error {
	return s.Container.Remove(options)
}

// Stop stops the service's containers
func (s *Service) Stop() error {
	templateString := s.OSCommand.Config.UserConfig.CommandTemplates.StopService
	command := utils.ApplyTemplate(templateString, s)
	return s.OSCommand.RunCommand(command)
}

// Restart restarts the service
func (s *Service) Restart() error {
	templateString := s.OSCommand.Config.UserConfig.CommandTemplates.RestartService
	command := utils.ApplyTemplate(templateString, s)
	return s.OSCommand.RunCommand(command)
}

// Attach attaches to the service
func (s *Service) Attach() (*exec.Cmd, error) {
	return s.Container.Attach()
}

// Top returns process information
func (s *Service) Top() (types.ContainerProcessList, error) {
	return s.Container.Top()
}

// ViewLogs attaches to a subprocess viewing the service's logs
func (s *Service) ViewLogs() (*exec.Cmd, error) {
	templateString := s.OSCommand.Config.UserConfig.CommandTemplates.ViewServiceLogs
	command := utils.ApplyTemplate(templateString, s)

	cmd := s.OSCommand.ExecutableFromString(command)
	// so long as this is commented in, the child process does not receive the interrupt
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	return cmd, nil
}