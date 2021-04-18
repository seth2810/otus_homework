package main

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CopyTestSuite struct {
	suite.Suite
	fromPath, toDir, toPath string
}

func (s *CopyTestSuite) SetupSuite() {
	s.toDir = os.TempDir()
	s.fromPath = "testdata/input.txt"
}

func (s *CopyTestSuite) BeforeTest(suiteName, testName string) {
	s.toPath = path.Join(s.toDir, fmt.Sprintf("%s-%s", suiteName, testName))
}

func (s *CopyTestSuite) TearDownTest() {
	_ = os.Remove(s.toPath)
}

func (s *CopyTestSuite) TearDownSuite() {
	_ = os.RemoveAll(s.toDir)
}

func (s *CopyTestSuite) CheckNotExist(path string) {
	_, err := os.Stat(path)
	require.ErrorIs(s.T(), err, os.ErrNotExist)
}

func (s *CopyTestSuite) TestFromPathRequired() {
	require.ErrorIs(s.T(), Copy("", s.toPath, 0, 0), ErrFromPathIsRequired)
}

func (s *CopyTestSuite) TestToPathRequired() {
	require.ErrorIs(s.T(), Copy(s.fromPath, "", 0, 0), ErrToPathIsRequired)
}

func (s *CopyTestSuite) TestFilesWithoutSize() {
	require.ErrorIs(s.T(), Copy("/dev/urandom", s.toPath, 0, 0), ErrUnsupportedFile)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestFromPathNotExist() {
	require.ErrorIs(s.T(), Copy("/dev/nil", s.toPath, 0, 0), os.ErrNotExist)
}

func (s *CopyTestSuite) TestFromPathNoPermission() {
	require.ErrorIs(s.T(), Copy("/etc/sudoers", s.toPath, 0, 0), os.ErrPermission)
}

func (s *CopyTestSuite) TestToPathNoPermission() {
	require.ErrorIs(s.T(), Copy(s.fromPath, "/etc/hosts", 0, 0), os.ErrPermission)
}

func (s *CopyTestSuite) TestDir() {
	require.ErrorIs(s.T(), Copy("/dev", s.toPath, 0, 0), ErrUnsupportedFile)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestNegativeOffset() {
	require.ErrorIs(s.T(), Copy(s.fromPath, s.toPath, -1, 0), ErrOffsetNegative)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestNegativeLimit() {
	require.ErrorIs(s.T(), Copy(s.fromPath, s.toPath, 0, -1), ErrLimitNegative)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestOffsetOutOfSize() {
	stat, _ := os.Stat(s.fromPath)

	require.ErrorIs(s.T(), Copy(s.fromPath, s.toPath, stat.Size()+1, 0), ErrOffsetExceedsFileSize)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestSamePath() {
	require.ErrorIs(s.T(), Copy(s.fromPath, s.fromPath, 0, 0), ErrFromToPathsSame)

	s.CheckNotExist(s.toPath)
}

func (s *CopyTestSuite) TestLimitGraterThanFileSize() {
	fromStat, _ := os.Stat(s.fromPath)

	require.NoError(s.T(), Copy(s.fromPath, s.toPath, 0, fromStat.Size()+1))

	stat, err := os.Stat(s.toPath)

	require.NoError(s.T(), err)
	require.Equal(s.T(), stat.Size(), fromStat.Size())
}

func (s *CopyTestSuite) TestOffsetSameAsFileSize() {
	fromStat, _ := os.Stat(s.fromPath)

	require.NoError(s.T(), Copy(s.fromPath, s.toPath, fromStat.Size(), 0))

	stat, err := os.Stat(s.toPath)

	require.NoError(s.T(), err)
	require.Equal(s.T(), stat.Size(), int64(0))
}

func (s *CopyTestSuite) TestDefaultLimit() {
	fromStat, _ := os.Stat(s.fromPath)

	require.NoError(s.T(), Copy(s.fromPath, s.toPath, 0, 0))

	stat, err := os.Stat(s.toPath)

	require.NoError(s.T(), err)
	require.Equal(s.T(), stat.Size(), fromStat.Size())
}

func (s *CopyTestSuite) TestCopyHead() {
	var copySize int64 = 10

	require.NoError(s.T(), Copy(s.fromPath, s.toPath, 0, copySize))

	stat, _ := os.Stat(s.toPath)

	require.Equal(s.T(), stat.Size(), copySize)
}

func (s *CopyTestSuite) TestCopyTail() {
	var copySize int64 = 10
	fromStat, _ := os.Stat(s.fromPath)

	require.NoError(s.T(), Copy(s.fromPath, s.toPath, fromStat.Size()-copySize, 0))

	stat, _ := os.Stat(s.toPath)

	require.Equal(s.T(), stat.Size(), copySize)
}

func TestCopy(t *testing.T) {
	suite.Run(t, new(CopyTestSuite))
}
