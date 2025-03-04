# Octagon Go CLI

Scripts and tools for automating the running of The Octagon, a Smash Ultimate tournament hosted every Tuesday at 7pm at the Octopus Bar in Seattle. <https://start.gg/octagon>

## Scripts

WIP: Last minute attendee adder, seeder, and conflict resolver

## Start.gg API

Calls the start.gg GraphQL api

[docs](https://developer.start.gg/docs/intro)
[schema](https://smashgg-schema.netlify.app/reference/query.doc.html)

## Development

When the start.gg GraphQL schema updates, you need to fetch the updated schema.graphql file and run

```
go get github.com/Khan/genqlient
npm install -g graphqurl
```

```bash
gq https://api.start.gg/gql/alpha --introspect > schema.graphql -H 'Authorization: Bearer <API_KEY>'
go run github.com/Khan/genqlient
```
