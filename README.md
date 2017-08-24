# db-migration-auditer
Recently we had a postgres stored procedure regression where a function was updated in migration 10, then updated again in migration 11 (using the definition from 10), then updated a third time in migration 12 (using the definition from migration 10 again, doh!)

This Golang tester should help catch these stored procedure regressions prior to committing.

Run:
 * Update constants.go to have the correct values
 * run `go test` in this repo

Assumptions:
 * There is a consistent start and end for all stored procedure definitions (see constants.go)
 * Migration files are ordered alphanumerically where the first one is the oldest migration
 * There might be others...

 TODO:
  * create a constants.go file for easier future re-use
