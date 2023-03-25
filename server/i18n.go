// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/pkg/errors"
)

func (p *Plugin) TranslationsPreInit() error {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return errors.Wrap(err, "unable to find i18n directory")
	}

	i18nDirectory := path.Join(bundlePath, "assets", "i18n")
	files, err := ioutil.ReadDir(i18nDirectory)
	if err != nil {
		return errors.Wrap(err, "unable to read i18n directory")
	}
	for _, f := range files {
		filename := f.Name()
		if filepath.Ext(filename) == ".json" {
			if err := i18n.LoadTranslationFile(filepath.Join(i18nDirectory, filename)); err != nil {
				p.API.LogError("Failed to load translation file", "filename", filename, "err", err.Error())
				continue
			}

			p.API.LogDebug("Loaded translation file", "filename", filename)
			p.locales[strings.TrimSuffix(filename, filepath.Ext(filename))] = filepath.Join(i18nDirectory, filename)
		}
	}

	return nil
}

func (p *Plugin) GetUserTranslations(locale string) i18n.TranslateFunc {
	if _, ok := p.locales[locale]; !ok {
		locale = model.DefaultLocale
	}

	translations := TfuncWithFallback(locale)
	return translations
}

func TfuncWithFallback(pref string) i18n.TranslateFunc {
	t, _ := i18n.Tfunc(pref)
	return func(translationID string, args ...interface{}) string {
		if translated := t(translationID, args...); translated != translationID {
			return translated
		}

		t, _ := i18n.Tfunc(model.DefaultLocale)
		return t(translationID, args...)
	}
}
