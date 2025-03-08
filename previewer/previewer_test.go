package previewer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/dir"
	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/ls"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	testHomePath     = "/home/test"
	testDownloadPath = testHomePath + "/Downloads"
	testCodePath     = testHomePath + "/Code/JSXQL"
)

type PreviewerTestSuite struct {
	suite.Suite
	mockLister  *lister.MockLister
	mockTmux    *tmux.MockTmux
	mockIcon    *icon.MockIcon
	mockDir     *dir.MockDir
	mockHome    *home.MockHome
	mockLs      *ls.MockLs
	mockShell   *shell.MockShell
	mockConnfig *model.Config
	previewer   Previewer
}

func (suite *PreviewerTestSuite) SetupTest() {
	suite.initializeMocks()
	suite.initializePreviewer()
}

func (suite *PreviewerTestSuite) TearDownTest() {
	suite.mockLister.AssertExpectations(suite.T())
	suite.mockTmux.AssertExpectations(suite.T())
	suite.mockIcon.AssertExpectations(suite.T())
	suite.mockDir.AssertExpectations(suite.T())
	suite.mockHome.AssertExpectations(suite.T())
	suite.mockLs.AssertExpectations(suite.T())
	suite.mockShell.AssertExpectations(suite.T())
}

func (suite *PreviewerTestSuite) initializeMocks() {
	suite.mockLister = new(lister.MockLister)
	suite.mockTmux = new(tmux.MockTmux)
	suite.mockIcon = new(icon.MockIcon)
	suite.mockDir = new(dir.MockDir)
	suite.mockHome = new(home.MockHome)
	suite.mockLs = new(ls.MockLs)
	suite.mockShell = new(shell.MockShell)
}

func (suite *PreviewerTestSuite) initializePreviewer() {
	suite.previewer = NewPreviewer(
		suite.mockLister,
		suite.mockTmux,
		suite.mockIcon,
		suite.mockDir,
		suite.mockHome,
		suite.mockLs,
		model.Config{},
		suite.mockShell,
	)
}

func TestPreviewerTestSuite(t *testing.T) {
	suite.Run(t, new(PreviewerTestSuite))
}

func (suite *PreviewerTestSuite) TestPreview_TmuxStrategy() {
	testCase := struct {
		inputName      string
		trimmedName    string
		expectedOutput string
	}{
		inputName:      " test-session",
		trimmedName:    "test-session",
		expectedOutput: "Fake tmux ansi output",
	}

	suite.setupTmuxMocks(testCase.inputName, testCase.trimmedName, testCase.expectedOutput)

	output, err := suite.previewer.Preview(testCase.inputName)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testCase.expectedOutput, output)
}

func (suite *PreviewerTestSuite) TestPreview_DefaultConfigStrategy() {
	testCase := struct {
		inputName      string
		trimmedName    string
		expectedPath   string
		expectedOutput string
	}{
		inputName:    " Downloads 📥",
		trimmedName:  "Downloads 📥",
		expectedPath: testDownloadPath,
		expectedOutput: `.rw-r--r-- 761k test  8 apr 17:56 export.csv
.rw-r--r--  93k test 17 feb 16:42 IMG-20240217-WA0002.jpg
.rw-r--r--  63k test  8 apr 15:55 'La stella.epub'
drwxrwxr-x    - test 16 dic 18:37 'Learning Go.pdf'`,
	}

	suite.setupDefaultConfigMocks(testCase.inputName, testCase.trimmedName, testCase.expectedPath, testCase.expectedOutput)

	output, err := suite.previewer.Preview(testCase.inputName)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testCase.expectedOutput, output)
}

func (suite *PreviewerTestSuite) TestPreview_ConfigStrategy() {
	testCase := struct {
		inputName           string
		trimmedName         string
		previewCommand      string
		previewCommandParts []string
		expectedPath        string
		expectedOutput      string
	}{
		inputName:           "📁 JSXQL",
		trimmedName:         "JSXQL",
		previewCommand:      "ls -la",
		previewCommandParts: []string{"ls", "-la"},
		expectedPath:        testCodePath,
		expectedOutput: `.rw-rw-r--   299 test 15 dic 16:17 -- global.d.ts
.rw-rw-r--   251 test 15 dic 16:17 -- index.tsx
.rw-rw-r-- 1,7Ki test 15 dic 16:17 -- jsxql.ts`,
	}

	suite.setupConfigMocks(testCase.inputName, testCase.trimmedName, testCase.previewCommand, testCase.previewCommandParts, testCase.expectedPath, testCase.expectedOutput)

	output, err := suite.previewer.Preview(testCase.inputName)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testCase.expectedOutput, output)
}

func (suite *PreviewerTestSuite) TestPreview_DirectoryStrategy() {
	testCase := struct {
		inputName      string
		trimmedName    string
		expectedPath   string
		expectedOutput string
	}{
		inputName:    " ~/Code/JSXQL/",
		trimmedName:  "~/Code/JSXQL/",
		expectedPath: testCodePath,
		expectedOutput: `.rw-rw-r--   299 test 15 dic 16:17 -- global.d.ts
.rw-rw-r--   251 test 15 dic 16:17 -- index.tsx
.rw-rw-r-- 1,7Ki test 15 dic 16:17 -- jsxql.ts`,
	}

	suite.setupDirectoryMocks(testCase.inputName, testCase.trimmedName, testCase.expectedPath, testCase.expectedOutput)

	output, err := suite.previewer.Preview(testCase.inputName)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), testCase.expectedOutput, output)
}

func (suite *PreviewerTestSuite) TestPreview_NoMatch() {
	testCase := struct {
		inputName   string
		trimmedName string
	}{
		inputName:   "nonexistent",
		trimmedName: "nonexistent",
	}

	suite.setupNoMatchMocks(testCase.inputName, testCase.trimmedName)

	output, err := suite.previewer.Preview(testCase.inputName)

	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), output)
}

func (suite *PreviewerTestSuite) setupTmuxMocks(inputName, trimmedName, expectedOutput string) {
	suite.mockIcon.On("RemoveIcon", inputName).Return(trimmedName)
	suite.mockLister.On("FindTmuxSession", trimmedName).Return(model.SeshSession{
		Name: trimmedName,
		Path: testHomePath + "/c/" + trimmedName,
	}, true)
	suite.mockTmux.On("CapturePane", trimmedName).Return(expectedOutput, nil)
}

func (suite *PreviewerTestSuite) setupDefaultConfigMocks(inputName, trimmedName, expectedPath, expectedOutput string) {
	suite.mockIcon.On("RemoveIcon", inputName).Return(trimmedName)
	suite.mockLister.On("FindTmuxSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockLister.On("FindConfigSession", trimmedName).Return(model.SeshSession{
		Name: trimmedName,
		Path: expectedPath,
	}, true)
	suite.mockLs.On("ListDirectory", expectedPath).Return(expectedOutput, nil)
}

func (suite *PreviewerTestSuite) setupConfigMocks(inputName string, trimmedName string, previewCommand string, previewCommandParts []string, expectedPath string, expectedOutput string) {
	suite.mockIcon.On("RemoveIcon", inputName).Return(trimmedName)
	suite.mockLister.On("FindTmuxSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockLister.On("FindConfigSession", trimmedName).Return(model.SeshSession{
		Name:           trimmedName,
		Path:           expectedPath,
		PreviewCommand: previewCommand,
	}, true)
	suite.mockShell.On("PrepareCmd", previewCommand, map[string]string{"{}": expectedPath}).Return(previewCommandParts, nil)
	suite.mockShell.On("Cmd", "ls", "-la").Return(expectedOutput, nil)
}

func (suite *PreviewerTestSuite) setupDirectoryMocks(inputName, trimmedName, expectedPath, expectedOutput string) {
	suite.mockIcon.On("RemoveIcon", inputName).Return(trimmedName)
	suite.mockLister.On("FindTmuxSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockLister.On("FindConfigSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockHome.On("ExpandHome", trimmedName).Return(expectedPath, nil)
	suite.mockDir.On("Dir", expectedPath).Return(true, expectedPath)
	suite.mockLs.On("ListDirectory", expectedPath).Return(expectedOutput, nil)
}

func (suite *PreviewerTestSuite) setupNoMatchMocks(inputName, trimmedName string) {
	suite.mockIcon.On("RemoveIcon", inputName).Return(trimmedName)
	suite.mockLister.On("FindTmuxSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockLister.On("FindConfigSession", trimmedName).Return(model.SeshSession{}, false)
	suite.mockHome.On("ExpandHome", trimmedName).Return("", nil)
	suite.mockDir.On("Dir", "").Return(false, "")
}
