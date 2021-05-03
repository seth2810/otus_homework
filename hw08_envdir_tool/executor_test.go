package main

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/rhysd/go-fakeio"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RunCmdTestSuite struct {
	suite.Suite
	out, err *fakeio.FakedIO
}

func (s *RunCmdTestSuite) BeforeTest(_, _ string) {
	s.err = fakeio.Stderr()
	s.out = fakeio.Stdout()
}

func (s *RunCmdTestSuite) TearDownTest() {
	s.err.Restore()
	s.out.Restore()
}

func (s *RunCmdTestSuite) TestCmdEmpty() {
	require.Equal(s.T(), RunCmd([]string{}, Environment{}), -1)
}

func (s *RunCmdTestSuite) TestCmdNotExist() {
	code := RunCmd([]string{"/dev/nil"}, Environment{})
	out, _ := s.out.String()

	require.Equal(s.T(), code, 1)
	require.Contains(s.T(), out, "no such file or directory")
}

func (s *RunCmdTestSuite) TestCmdNoPermission() {
	code := RunCmd([]string{"/etc/sudoers"}, Environment{})
	out, _ := s.out.String()

	require.Equal(s.T(), code, 1)
	require.Contains(s.T(), out, "permission denied")
}

func (s *RunCmdTestSuite) TestCmdExitError() {
	code := RunCmd([]string{"cat", "/etc/sudoers"}, Environment{})
	err, _ := s.err.String()

	require.Equal(s.T(), code, 1)
	require.Contains(s.T(), err, "/etc/sudoers: Permission denied")
}

func (s *RunCmdTestSuite) TestCmdExitOk() {
	code := RunCmd([]string{"echo", "ok"}, Environment{})
	out, _ := s.out.String()

	require.Equal(s.T(), code, 0)
	require.Equal(s.T(), out, "ok\n")
}

func (s *RunCmdTestSuite) TestPassOsEnv() {
	RunCmd([]string{"env"}, Environment{})

	out, _ := s.out.String()

	require.Equal(s.T(), out, fmt.Sprintln(strings.Join(os.Environ(), "\n")))
}

func (s *RunCmdTestSuite) TestUnsetEnvVar() {
	os.Setenv("UNSET", "VALUE")

	RunCmd([]string{"env"}, Environment{
		"UNSET": EnvValue{NeedRemove: true},
	})

	out, _ := s.out.String()

	require.NotContains(s.T(), out, "UNSET=VALUE")
	require.Empty(s.T(), os.Getenv("UNSET"))
}

func (s *RunCmdTestSuite) TestOverrideEnvVar() {
	os.Setenv("OVERRIDE", "OLD")

	RunCmd([]string{"env"}, Environment{
		"OVERRIDE": EnvValue{"NEW", false},
	})

	out, _ := s.out.String()

	require.Contains(s.T(), out, "OVERRIDE=NEW")
	require.Equal(s.T(), os.Getenv("OVERRIDE"), "NEW")
}

func TestRunCmd(t *testing.T) {
	suite.Run(t, new(RunCmdTestSuite))
}
