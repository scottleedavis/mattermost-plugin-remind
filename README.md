# Mattermost Plugin Remind 

[![Build Status](https://img.shields.io/circleci/project/github/scottleedavis/mattermost-plugin-remind/master.svg)](https://circleci.com/gh/scottleedavis/mattermost-plugin-remind)  [![codecov](https://codecov.io/gh/scottleedavis/mattermost-plugin-remind/branch/master/graph/badge.svg)](https://codecov.io/gh/scottleedavis/mattermost-plugin-remind)  [![Go Report Card](https://goreportcard.com/badge/github.com/scottleedavis/mattermost-plugin-remind)](https://goreportcard.com/report/github.com/scottleedavis/mattermost-plugin-remind)  [![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3119/badge)](https://bestpractices.coreinfrastructure.org/projects/3119)  [![Releases](https://img.shields.io/github/release/scottleedavis/mattermost-plugin-remind.svg)](https://github.com/scottleedavis/mattermost-plugin-remind/releases/latest)

_**A bot that schedules reminders for [Mattermost](https://mattermost.com/)**_

<img src="remind.png">

### Installation

_requires Mattermost 5.14 or greater._

1) Go to the [releases page](https://github.com/scottleedavis/mattermost-plugin-remind/releases) of this GitHub repository and download the latest release for your Mattermost server.
2) Upload this file in the Mattermost System Console > Plugins > Management page to install the plugin. To learn more about how to upload a plugin, see the documentation.
3) For a better cross timezone experience, enable Experimental timezone support.  `System Console -> Experimental Features -> Timezone  = true`


### Usage

See the full list of [Usage Examples](https://github.com/scottleedavis/mattermost-plugin-remind/wiki/Usage) in the [wiki](https://github.com/scottleedavis/mattermost-plugin-remind/wiki) 
* `/remind` - opens up an [interactive dialog](https://docs.mattermost.com/developer/interactive-dialogs.html) to schedule a reminder
* `/remind help` - displays help examples
* `/remind list` - displays a list of reminders
* `/remind [who] [what] [when]`
  * `/remind [who] [what] in [# (seconds|minutes|hours|days|weeks|months|years)]`
  * `/remind [who] [what] at [(noon|midnight|one..twelve|00:00am/pm|0000)] (every) [day|date]`
  * `/remind [who] [what] (on) [(monday-sunday|month&day|m/d/y|d.m.y)] (at) [time]`
  * `/remind [who] [what] every (other) [monday,...,sunday|weekdays|month&day|m/d|d.m] (at) [time]`
* `/remind [who] [when] [what]`

### Build

#### Requirements
* [Go 1.12](https://golang.org/)
* [Dep](https://github.com/golang/dep)

```
make
```

This will produce a single plugin file (with support for multiple architectures) for upload to your Mattermost server:

```
dist/com.github.scottleedavis.mattermost-plugin-remind.tar.gz
```

There is a build target to automate deploying and enabling the plugin to your server, but it requires configuration as below:
```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
make deploy
```
In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

