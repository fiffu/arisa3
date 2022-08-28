## arisa3

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/fiffu/arisa3/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/fiffu/arisa3/tree/main)

Arisa is a [Discord](https://discordapp.com/) bot written in Go, using the
[discordgo](https://github.com/bwmarrin/discordgo) library.

arisa3 is a continuation of predecessors [Arisa](https://arisa-chan.herokuapp.com) and
[arisa2](https://github.com/fiffu/arisa2).

The main motivation for this rewrite is to adopt a relatively new feature by Discord, called
[Application Commands](https://discord.com/developers/docs/interactions/application-commands).
This project rewrite also aims to provide better stability through unit testing and
leveraging Go's language features, such as static typing.

## Architecture

#### Business logic is contained in cogs

arisa3 inherits the concept of [cogs](https://discordpy.readthedocs.io/en/latest/ext/commands/cogs.html)
from arisa2. Each cog independently implements the business logic for a particular group
of features - they are structurally comparable to Usecases (per *Clean Architecture*) or
Domains (per *Domain-Driven Design*).

The most basic cog structure directly maps commands to module-level functions. For more
complicated commands that require backing services such as databases, the cog initializes
a Domain (containing business logic) that connects to a Repository interface (the database)
or other external resources (e.g. HTTP APIs).

#### Config

This project follows the [Twelve Factor guidelines for configs](https://12factor.net/config).
During startup, the app expects to read a path to a config file. If any config fields are
[tagged](https://www.digitalocean.com/community/tutorials/how-to-use-struct-tags-in-go)
with `envvar`, those fields are overwritten by the corresponding environment variables.

Configs exist at both app-level (e.g. for database connections) and cog-level. Cog-level
configs are loaded during dependency injection.

#### Dependency injection

Every cog must implement the `OnStartup()` hook. This hook is triggered during app startup,
*before* logging on to Discord's API gateway. This hook is intended to initiate dependency
injection.

Most cogs implement the `IBootable` and `IRepository` interfaces, which makes them eligible
to use the helper method `engine.Bootstrap()` to automatically ingest cog-level configs,
setup event handlers, and run database migrations.

#### Stack

arisa3 talks directly to PostgreSQL without any ORM layer. Some degree of flexibility is
available thanks to Go's `database/sql` abstraction and 12-Factor configurability, but
the migrations have to be translated from the PostgreSQL dialect.

## Contributing

#### Setup

1. Set up Go locally, following these [instructions](https://go.dev/doc/install).
2. Run `make install-dev` - downloads build deps and installs dev tooling (e.g. pre-commit)
   to your machine.
3. Run `make test` - execute unit tests.
4. Run `cp config.sample.yml config.yml` - this creates a config file based on the sample.
   The default config doesn't come with secrets (such as database DSN and API keys). You
   will have to provide those values.
5. Run `go build -o arisa3 && ./arisa3 -config-file ./config.yml` - this builds and runs
   the binary. On Windows, you might have to use `arisa3.exe` instead.

#### Repo conventions

1. Try to open an issue for each pull request.
2. Use `main` as the source for branching, and the merge target for pull requests.
