// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils/fileutils"
	"github.com/nicksnyder/go-i18n/i18n"
)

var locales map[string]string = make(map[string]string)

func TranslationsPreInit() error {

	i18nDirectory, found := fileutils.FindDir("plugins/" + manifest.Id + "/server/dist/i18n/")
	if !found {
		return fmt.Errorf("unable to find i18n directory")
	}

	files, _ := ioutil.ReadDir(i18nDirectory)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			filename := f.Name()
			locales[strings.Split(filename, ".")[0]] = filepath.Join(i18nDirectory, filename)

			if err := i18n.LoadTranslationFile(filepath.Join(i18nDirectory, filename)); err != nil {
				return err
			}
		}
	}

	return nil
}

func GetUserTranslations(locale string) i18n.TranslateFunc {
	if _, ok := locales[locale]; !ok {
		locale = model.DEFAULT_LOCALE
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

		t, _ := i18n.Tfunc(model.DEFAULT_LOCALE)
		return t(translationID, args...)
	}
}
