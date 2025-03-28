# Octagon Go CLI

Scripts and tools for automating the running of The Octagon,
a Smash Ultimate tournament hosted every Tuesday
at 7pm at the Octopus Bar in Seattle.
<https://start.gg/octagon>

## Start.gg API

Calls the start.gg GraphQL API

[docs](https://developer.start.gg/docs/intro)
[schema](https://smashgg-schema.netlify.app/reference/query.doc.html)

## Development

Run the CLI with the following command

```bash
go run . <args>
```

If you want to see debug logs,
add a DEBUG environment variable in the config or inline

```bash
DEBUG=1 go run . <args>
```

## Contribution

If you want to contribute to the project,
feel free to check out the GitHub issues and pick one up.

### Testing

Whenever possible, add test cases for features

### Conventional Commit Messages

Please try to use [conventional commit messages](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional)
whenever possible.

Examples

```text
feat: 'conflict remove' subcommand
fix: increase pageSize in getParticipants to prevent missing players
perf: cached getTournament responses
refactor: clean up bracket creation logic
docs: improve contribution instructions
chore: remove unnused commented code
style: whitespace changes to conflict list initialization
```

## Configuration

When running the `octagon` command it will check for a configuration file in `~/.config/octagon/octagonrc`.
If it is unable to find a config file there, it will look for a local .env file.

### Config fields

```dotenv
STARTGG_API_KEY=
FIREBASE_API_KEY=
FIREBASE_DATABASE_URL=
```

## Installation & Usage

This will add the `octagon` command to your GOPATH,
allowing the command to be run from anywhere.

```bash
go build && go install
```

If you are using a local .env configuration instead of ~/.config/octagon,
the command will fail outside of the source code directory.

View the command usage

```bash
octagon -h
```

## Updating Schemas

When the start.gg GraphQL schema updates,
you need to fetch the updated `schema.graphql` file and run

```bash
go get github.com/Khan/genqlient
npm install -g graphqurl
```

Dev API

```bash
cd startgg
gq https://api.start.gg/gql/alpha --introspect > schema.graphql \
  -H 'Authorization: Bearer <API_KEY>'
go run github.com/Khan/genqlient
```

Prod API

```bash
cd startgg/prod
gq https://www.start.gg/api/-/gql --introspect \
-H 'x-web-source: gg-web-gql-client, gg-web-rest'
go run github.com/Khan/genqlient
```
