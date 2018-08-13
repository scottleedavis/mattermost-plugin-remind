# mattermost-plugin-remind
a plugin for mattermost that provides reminders from a slash command /remind

### setup
config.json changes

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
