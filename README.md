# Convoy

### Convoy is a Kubernetes controller that allows your team to receive notifications about specific events in your cluster.

Convoy is a Kubernetes controller that can notify you about specific cluster events.
For example if you want to watch all `Pod` events in a certain namespace then Convoy can dispatch
these events to your team via a Slack channel.



## Features

* Post Kubernetes events into Slack (or other)
* Filter events (type, namespace etc)

TODO take action on specific events (automation)

![Convoy](https://github.com/owainlewis/convoy/blob/master/convoy.png?raw=true)
