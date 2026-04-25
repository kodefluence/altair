package usecase

// This file exists to surface internal test seams to *_test packages outside
// `usecase`. Keeping the seams in their own file (rather than runner.go)
// makes the production surface obvious in code review — anything declared
// here exists *only* to support tests.

// MigratorForTest re-exports the unexported migrator interface so callers
// in other packages (e.g. controller/command tests) can build a fake
// migrator without making the interface part of the production API.
type MigratorForTest = migrator

// MigratorFactoryForTest re-exports the unexported migratorFactory so
// callers in other packages can construct a stub factory and inject it
// via SetMigratorFactoryForTest.
type MigratorFactoryForTest = migratorFactory

// SetMigratorFactoryForTest swaps the Runner's migrator factory. Used by
// command tests that need to drive Runner without a real MySQL.
func SetMigratorFactoryForTest(r *Runner, f MigratorFactoryForTest) {
	r.migratorFactory = f
}
