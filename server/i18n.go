// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mattermost/mattermost-server/mlog"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/utils/fileutils"
	"github.com/nicksnyder/go-i18n/i18n"
)

var T i18n.TranslateFunc
var TDefault i18n.TranslateFunc
var locales map[string]string = make(map[string]string)
var settings model.LocalizationSettings

// this functions loads translations from filesystem
// and assign english while loading server config
func TranslationsPreInit() error {
	// Set T even if we fail to load the translations. Lots of shutdown handling code will
	// segfault trying to handle the error, and the untranslated IDs are strictly better.
	T = TfuncWithFallback("en")
	TDefault = TfuncWithFallback("en")
	return InitTranslationsWithDir()
}

func InitTranslationsWithDir() error {

	i18nDirectory, found := fileutils.FindDir("plugins/" + manifest.Id + "/server/dist/i18n/")
	if !found {
		return fmt.Errorf("Unable to find i18n directory")
	}

	files, _ := ioutil.ReadDir(i18nDirectory)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			filename := f.Name()
			locales[strings.Split(filename, ".")[0]] = filepath.Join(i18nDirectory, filename)

			if err := i18n.LoadTranslationFile(filepath.Join(i18nDirectory, filename)); err != nil {
				return err
			}
			mlog.Info("found the files?")
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

func GetSupportedLocales() map[string]string {
	return locales
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
