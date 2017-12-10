# Convoy

### Convoy is a Kubernetes controller that allows your team to receive notifications about specific events in your cluster.

Convoy is a Kubernetes controller that can notify you about specific cluster events.
For example if you want to watch all `Pod` events in a certain namespace then Convoy can dispatch
these events to your team via a Slack channel.

![Convoy](https://github.com/owainlewis/convoy/blob/master/convoy.png?raw=true)

## Roadmap

* Ability to filter by source component ("Notify me of all events generated by the deployment controller in the default namespace")
* Ability to filter by event type ("Notify me of all Pod events in the cluster")
* Alerts

## Configuration

Event type: i.e Pod | Secret
