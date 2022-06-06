# Exercism Events

Creates custom events for Exercism and defines processors/consumers for those events that can act on them.

## Download

See the [Releases](https://github.com/andrerfcsantos/exercism-events/releases) page of the repository. 

## Usage

### Authentication

For the program to work, it needs an Exercism API Token. You can see your API token for Exercism [here](https://exercism.org/settings/api_cli).

The token can be given to the program by setting the `EXERCISM_TOKEN` environment variable.
If this variable is not set, the program will try to retrieve it from the config dir of [Exercism CLI](https://github.com/exercism/cli).

### Running

```
exercism-events [flags]

Available Flags:

-tracks     Comma-separated list of tracks for which to track mentoring requests.
            (default: tracks the user is mentor of)
```

### Examples

- Tracks new mentoring requests for the tracks the user is mentoring:

```console
$ exercism-events
```

- Tracks new mentoring requests for just `go` and `python`:

```console
$ exercism-events -tracks=go,python
```

## Available sources

Sources generate events. These are the available source:

* **Mentoring requests** - Emit events when a mentoring request enters or leaves the queue

**Note:** More sources are planned and it'll be possible to activate/deactivate individual sources. Right now the "mentoring requests" source is the only source available and is always enabled.
## Available consumers

Consumers take the events produced by the sources and perform some kind of action on them. These are the available processors:

* **Desktop Notifier** - Displays desktop notifications for the events produced by the sources.

**Note:** More consumers are planned and it'll be possible to activate/deactivate individual consumers. Right now the desktop notifier is the only consumer and is always enabled.

