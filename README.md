# Exercism Events

Creates custom events for Exercism and defines processors/consumers for those events that can act on them.

## Download

See the [Releases](https://github.com/andrerfcsantos/exercism-events/releases) page of the repository.

Alternatively, if you have Go installed, you can run:

```console
go install github.com/andrerfcsantos/exercism-events@latest
```

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
-sources    Comma-separated list of event sources that will be activated.
            (default: notifications,mentoring)
-consumers  Comma-separated list of event consumers that will be activated.
            (default: desktopnotifier)
-pushovertracks
```

### Examples

- Tracks new mentoring requests or the tracks the user is mentoring and exercism notifications:

```console
$ exercism-events
```

- Tracks new mentoring requests for just `go` and `python`:

```console
$ exercism-events -tracks=go,python
```

- Receive desktop notifications just for Exercism notifications:

```console
$ exercism-events -sources=notifications -consumers=desktopnotifier
```

- Save mentoring requests to a database (leaving desktop notifications disabled):

```console
$ exercism-events -sources=mentoring -consumers=database
```

See below for more details on how to use the consumer `database`.


## Available sources

Sources generate events. These are the available sources:

* **Mentoring requests** - Emit events when a mentoring request enters or leaves the queue. Pass the key `mentoring` to the `-sources` flag to activate this source. The flag `-tracks` affects the tracks for which mentoring requests will be monitored.

* **Exercism Notifications** - Emit events when you receive a notification in exercism. Pass the key `notifications` to the `-sources` flag to activate this source.

## Available consumers

Consumers take the events produced by the sources and perform some kind of action on them. These are the available processors:

* **Desktop Notifier** - Displays desktop notifications for the events produced by the sources. Pass the key `desktopnotifier` to the `-consumers` flag to activate this consumer.

* **Database** - Saves events to a database. Pass the key `database` to the `-consumers` flag to activate this consumer.

For the database consumer to work, you must set these environment variables:

- `EXERCISM_EVENTS_DB_HOST` - database host (e.g `mydb.mydomain.com`)
- `EXERCISM_EVENTS_DB_USER` - database user (e.g `postgres`)
- `EXERCISM_EVENTS_DB_PASSWORD` - password for the `EXERCISM_EVENTS_DB_USER` user (e.g `admin1234`)
- `EXERCISM_EVENTS_DB_NAME` - database name to use (eg. `exercism-events`)
- `EXERCISM_EVENTS_DB_PORT` - database port (e.g `5432`)

Currently only postgres databases are supported.

In addition to setting the environment variables, [liquibase](https://www.liquibase.org/) migrations in the file `liquibase/master.xml` must also be run before running the program. Please refer to the liquibase documentation on how to do this.

* **pushover** - Sends notifications to [pushover](https://pushover.net/). Pass the key `pushover` to the `-consumers` flag to activate this consumer.

For this consumer to work, you must set these environment variables:

- `PUSHOVER_TOKEN` - Your App Token. You must create an application in pushover to get this token.
- `PUSHOVER_USER` - Your pushover User Key.

