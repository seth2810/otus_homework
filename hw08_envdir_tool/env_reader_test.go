package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ReadDirTestSuite struct {
	suite.Suite
	dir string
}

func (s *ReadDirTestSuite) BeforeTest(suiteName, _ string) {
	s.dir, _ = ioutil.TempDir(os.TempDir(), suiteName)
}

func (s *ReadDirTestSuite) TearDownTest() {
	_ = os.Remove(s.dir)
}

func (s *ReadDirTestSuite) CreateEnvFile(dir, name string, content []byte) {
	ioutil.WriteFile(path.Join(dir, name), content, fs.ModePerm)
}

func (s *ReadDirTestSuite) TestDirNotExist() {
	_, err := ReadDir("/dev/nil")

	require.ErrorIs(s.T(), err, os.ErrNotExist)
}

func (s *ReadDirTestSuite) TestDirNoPermission() {
	_, err := ReadDir("/etc/sudoers")

	require.ErrorIs(s.T(), err, os.ErrPermission)
}

func (s *ReadDirTestSuite) TestDirEmpty() {
	env, err := ReadDir(s.dir)

	require.NoError(s.T(), err)
	require.Equal(s.T(), env, Environment{})
}

func (s *ReadDirTestSuite) TestFileNoPermission() {
	_, err := ReadDir("/etc")

	require.ErrorIs(s.T(), err, os.ErrPermission)
}

func (s *ReadDirTestSuite) TestRegularVariable() {
	s.CreateEnvFile(s.dir, "S", []byte("T"))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{
		"S": EnvValue{"T", false},
	})
}

func (s *ReadDirTestSuite) TestTrimSpacesAndEOL() {
	spacing := strings.Repeat(" ", 2)

	s.CreateEnvFile(s.dir, "S1", []byte(fmt.Sprintf("T1\t%s\t", spacing)))
	s.CreateEnvFile(s.dir, "S2", []byte(fmt.Sprintf("T2%s\t%s", spacing, spacing)))
	s.CreateEnvFile(s.dir, "S3", []byte(fmt.Sprintf("\t%s\tT3", spacing)))
	s.CreateEnvFile(s.dir, "S4", []byte(fmt.Sprintf("%s\t%sT4", spacing, spacing)))
	s.CreateEnvFile(s.dir, "S5", []byte(fmt.Sprintf("T\t%s\t5", spacing)))
	s.CreateEnvFile(s.dir, "S6", []byte(fmt.Sprintf("T%s\t%s6", spacing, spacing)))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{
		"S1": EnvValue{"T1", false},
		"S2": EnvValue{"T2", false},
		"S3": EnvValue{fmt.Sprintf("\t%s\tT3", spacing), false},
		"S4": EnvValue{fmt.Sprintf("%s\t%sT4", spacing, spacing), false},
		"S5": EnvValue{fmt.Sprintf("T\t%s\t5", spacing), false},
		"S6": EnvValue{fmt.Sprintf("T%s\t%s6", spacing, spacing), false},
	})
}

func (s *ReadDirTestSuite) TestMultiline() {
	s.CreateEnvFile(s.dir, "S", []byte("T1\nT2"))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{
		"S": EnvValue{"T1", false},
	})
}

func (s *ReadDirTestSuite) TestTerminalZeros() {
	s.CreateEnvFile(s.dir, "S", []byte("T1\x00T2"))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{
		"S": EnvValue{"T1\nT2", false},
	})
}

func (s *ReadDirTestSuite) TestEmptyFile() {
	s.CreateEnvFile(s.dir, "S", []byte{})

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{
		"S": EnvValue{NeedRemove: true},
	})
}

func (s *ReadDirTestSuite) TestNameWithEqualSign() {
	s.CreateEnvFile(s.dir, "S=", []byte("T"))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{})
}

func (s *ReadDirTestSuite) TestIgnoreChildDir() {
	dir, _ := ioutil.TempDir(s.dir, "env")

	s.CreateEnvFile(dir, "S", []byte("T"))

	env, _ := ReadDir(s.dir)

	require.Equal(s.T(), env, Environment{})
}

func TestReadDir(t *testing.T) {
	suite.Run(t, new(ReadDirTestSuite))
}
