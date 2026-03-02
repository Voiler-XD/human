# human-plugin-graphql

A GraphQL schema and resolver generator plugin for the Human compiler.

## What It Generates

- `schema.graphql` — Complete GraphQL schema with types, inputs, enums, queries, and mutations
- `resolvers/index.ts` — Resolver index that aggregates all model resolvers
- `resolvers/<model>.ts` — CRUD resolvers per data model (Prisma-style)

## Category

backend

## Quick Start

```bash
# Build the plugin
make build

# Install into Human
make install

# Verify installation
human plugin list
```

## Type Mapping

| Human IR Type | GraphQL Type |
|---------------|-------------|
| `text`, `string` | `String` |
| `number`, `integer` | `Int` |
| `float`, `decimal` | `Float` |
| `boolean` | `Boolean` |
| `email` | `String` |
| `date`, `datetime` | `DateTime` |
| `enum` | Generated enum type |

## IR Fields Used

| Field | Usage |
|-------|-------|
| `app.Name` | Schema comments |
| `app.Data` | GraphQL types, input types, enums |
| `app.Data[].Fields` | Type fields with nullability from `required` |
| `app.Data[].Relations` | Relation fields (has_many → [Type!]!, belongs_to → Type) |
| `app.APIs` | Future: custom query/mutation mapping |
| `app.Auth` | Future: @auth directive generation |

## Generated Schema Example

For a .human file with a `Task` model:

```graphql
type Task {
  id: ID!
  title: String!
  status: TaskStatusEnum!
  createdAt: DateTime!
  updatedAt: DateTime!
}

enum TaskStatusEnum {
  PENDING
  DONE
}

input CreateTaskInput {
  title: String!
  status: TaskStatusEnum!
}

type Query {
  task(id: ID!): Task
  tasks(limit: Int, offset: Int): [Task!]!
}

type Mutation {
  createTask(input: CreateTaskInput!): Task!
  updateTask(id: ID!, input: UpdateTaskInput!): Task!
  deleteTask(id: ID!): Boolean!
}
```
