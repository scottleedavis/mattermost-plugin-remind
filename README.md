# mattermost-plugin-remind

a plugin for mattermost that provides reminders from a slash command /remind

Inspired by the java integration [mattermost-remind](http://github.com/scottleedavis/mattermost-remind)

### setup
Mattermost config.json changes

Setup siteUrl (set to your url)
```
...
"SiteURL": "http://127.0.0.1",
...
```

Enable timezones (each user set timezone)
```
...
   "DisplaySettings": {
        "CustomUrlSchemes": [],
        "ExperimentalTimezone": true
    },
...
```

Enable plugin uploads
```
...
"PluginsSetting": {
  ...
  "EnableUploads": true,
  ...
}
```

```
docker stop containerid
docker start containerid
```

### Building 
```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/io.github.mattermost-plugin-remind.tar.gz
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration and [http](https://httpie.org/) to be installed:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065/
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```
In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).


