package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTranslationsPreInit(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "TestTranslationsPreInit")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	assetsPath := filepath.Join(tmpDir, "assets")
	err = os.Mkdir(assetsPath, 0777)
	require.NoError(t, err)

	i18nPath := filepath.Join(tmpDir, "assets", "i18n")

	t.Run("failure to get bundle path", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetBundlePath", mock.Anything).Return(tmpDir, nil)

		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api
		err := p.TranslationsPreInit()
		require.EqualError(t, err, fmt.Sprintf("unable to read i18n directory: open %s: no such file or directory", i18nPath))
	})

	t.Run("failure to read i18n directory", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetBundlePath", mock.Anything).Return(tmpDir, nil)

		defer api.AssertExpectations(t)

		file, err := os.Create(i18nPath)
		require.NoError(t, err)
		file.Close()
		defer os.Remove(file.Name())

		p := &Plugin{}
		p.API = api
		err = p.TranslationsPreInit()
		require.True(t, strings.Contains(err.Error(), "unable to read i18n directory"))

	})

	t.Run("no i18n files", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetBundlePath", mock.Anything).Return(tmpDir, nil)

		defer api.AssertExpectations(t)

		err := os.Mkdir(i18nPath, 0777)
		require.NoError(t, err)
		defer os.Remove(i18nPath)

		p := &Plugin{}
		p.API = api
		err = p.TranslationsPreInit()
		require.NoError(t, err)
	})

	t.Run("various i18n files", func(t *testing.T) {
		api := &plugintest.API{}
		api.On("GetBundlePath", mock.Anything).Return(tmpDir, nil)
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
			t.Helper()
			t.Log(args...)
		})
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
			t.Helper()
			t.Log(args...)
		})

		defer api.AssertExpectations(t)

		err := os.Mkdir(i18nPath, 0777)
		require.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(i18nPath, "not-i18n.txt"), []byte{}, 0777)
		require.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(i18nPath, "invalid.json"), []byte{}, 0777)
		require.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(i18nPath, "en.json"), []byte(`[{"id":"id","translation":"translation"}]`), 0777)
		require.NoError(t, err)

		err = ioutil.WriteFile(filepath.Join(i18nPath, "es.json"), []byte(`[{"id":"id","translation":"translation2"}]`), 0777)
		require.NoError(t, err)

		p := NewPlugin()
		p.API = api
		err = p.TranslationsPreInit()
		require.NoError(t, err)

		require.Equal(t, map[string]string{
			"en": filepath.Join(i18nPath, "en.json"),
			"es": filepath.Join(i18nPath, "es.json"),
		}, p.locales)
	})
}
