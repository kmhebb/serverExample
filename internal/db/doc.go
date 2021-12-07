// Package db is our data layer and as with the rest of stack functions here
// should generally follow a pattern, though the exact structure isn't as well-
// defined as the service layer due to it's place at the bottom of the stack.
// In general however, our database functions all take a cloud.Context as the
// first parameter and a pg.Tx as the second. By only operating using the pg.Tx
// type, we allow the service layer to define transactional boundaries without
// needing to know the implementation details of how transactions are handled.
// By providing the full cloud.Context, this layer has access to all of the
// request-specific context required to perform the database operation such as
// scoping database queries.
package db
