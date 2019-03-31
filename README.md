# mattermost-plugin-i18n-test

Hello!

I am fiddling with getting translation functionality into the server portion of a plugin.  Here is where I am at.  Always open to any changes or thoughts.

* created a server/i18n directory, with same structure as mm-server
* copied the utils/i18n.go to server/i18n.go.  Initialized it with the absolute path e.g. `"/mm/mattermost/plugins/io.github.mattermost-plugin-i18n-test/server/dist/i18n/")`
* called TranslationsPreInit in OnActivate
* copied the i18n folder into the dist/ folder


---

Build plugin:
```
make
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```
