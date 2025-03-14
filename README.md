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

When the start.gg GraphQL schema updates,
you need to fetch the updated `schema.graphql` file and run

```bash
go get github.com/Khan/genqlient
npm install -g graphqurl
```

## Updating Schemas

Dev API

```bash
gq https://api.start.gg/gql/alpha --introspect > schema.graphql \
  -H 'Authorization: Bearer <API_KEY>'
go run github.com/Khan/genqlient
```

Prod API

```bash
gq https://www.start.gg/api/-/gql --introspect \
-H 'x-web-source: gg-web-gql-client, gg-web-rest'
go run github.com/Khan/genqlient
```

## Installation

```bash
go build && go install
```
