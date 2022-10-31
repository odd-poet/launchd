package launchd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Bootstrap a new service (`launchctl bootstrap` née `launchctl load`)
func (s *Service) Bootstrap() ([]byte, error) {
	path, err := s.DefinitionPath()
	if err != nil {
		return nil, err
	}
	return s.launchctl("bootstrap", domain, path)
}

// Bootout uninstalls a service (`launchctl bootout` née `launchctl unload`)
// Optionally deletes the service definition file (plist) as well.
func (s *Service) Bootout(removePlist bool) error {
	err := s.launchctlUnaryCommand("bootout")
	if err != nil {
		return err
	}
	if removePlist {
		path, err := s.DefinitionPath()
		if err != nil {
			return fmt.Errorf("could not find definition path: %w", err)
		}
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("could not remove definition: %w", err)
		}
	}
	return nil
}

// Start a service (`launchctl start`)
func (s *Service) Start() error {
	return s.launchctlUnaryCommand("start")
}

// Stop a service (`launchctl stop`)
func (s *Service) Stop() error {
	return s.launchctlUnaryCommand("stop")
}

// Print service state (`launchctl print`)
func (s *Service) Print() ([]byte, error) {
	return s.launchctl("print", s.UserSpecifier())
}

// Run a launchctl(1) subcommand against a service with the only argument being
// the service name. Returns error if not successful.
func (s *Service) launchctlUnaryCommand(command string) error {
	_, err := s.launchctl(command, s.UserSpecifier())
	return err
}

// Run a launchctl(1) subcommand for the service and return the output or an error
func (s *Service) launchctl(args ...string) ([]byte, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("launchctl", args...)
	cmd.Stdin = strings.NewReader("")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return nil, fmt.Errorf("(%w) running `launchctl %v` for %s:\n%s", err, args, s.Name, stderr.String())
	}
	return stdout.Bytes(), nil
}
