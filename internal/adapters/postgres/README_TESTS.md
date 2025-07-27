# Repository Tests

This directory contains unit and integration tests for PostgreSQL repository implementations.

## Test Types

### Unit Tests
- **Files**: `*_test.go` (without build tags)
- **Purpose**: Test business logic, validation, and SQL query structure
- **Dependencies**: None (uses sqlmock for database mocking)
- **Run with**: `go test ./internal/adapters/postgres/...`

### Integration Tests
- **Files**: `integration_test.go` (with `// +build integration` tag)
- **Purpose**: Test actual database interactions with PostgreSQL
- **Dependencies**: Requires running PostgreSQL database
- **Run with**: `go test -tags=integration ./internal/adapters/postgres/...`

## Running Tests

### Unit Tests Only
```bash
go test ./internal/adapters/postgres/...
```

### Integration Tests Only
```bash
go test -tags=integration ./internal/adapters/postgres/...
```

### All Tests
```bash
go test -tags=integration ./internal/adapters/postgres/...
```

## Test Database Setup

For integration tests, you need a PostgreSQL database running with:

- **Host**: localhost
- **Port**: 5432
- **Database**: smart_goal_calendar_test
- **Username**: postgres
- **Password**: postgres

### Using Docker
```bash
docker run -d \
  --name postgres-test \
  -e POSTGRES_DB=smart_goal_calendar_test \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15
```

### Manual Setup
```sql
CREATE DATABASE smart_goal_calendar_test;
CREATE USER postgres WITH PASSWORD 'postgres';
GRANT ALL PRIVILEGES ON DATABASE smart_goal_calendar_test TO postgres;
```

## Test Coverage

### UserRepository
- ✅ JSON marshaling/unmarshaling for Profile and Settings
- ✅ Basic validation logic
- ✅ SQL query structure validation
- ✅ Integration test for CRUD operations

### GoalRepository  
- ✅ Goal entity validation
- ✅ Progress validation (0-100 range)
- ✅ Category and Status constants
- ✅ Deadline handling (nullable)
- ✅ Integration test for CRUD operations

### EventRepository
- ✅ Event entity validation
- ✅ Time range validation (end >= start)
- ✅ Status constants
- ✅ External integration fields (Google Calendar)
- ✅ Goal linking (optional association)
- ✅ Timezone handling

### MoodRepository
- ✅ Mood level validation (1-5 range)
- ✅ Mood level string and emoji representations
- ✅ Tag management (add/remove/has)
- ✅ Date comparison logic
- ✅ Integration test for CRUD operations

## Test Helpers

### `test_helper.go`
Provides utilities for integration tests:

- `SetupTestDB()` - Creates test database connection and runs migrations
- `CreateTestUser()` - Creates a test user for other entity tests
- `CreateTestGoal()` - Creates a test goal
- `CreateTestEvent()` - Creates a test event
- `CreateTestMood()` - Creates a test mood
- `cleanupTestData()` - Removes all test data after tests

## Best Practices

1. **Unit Tests**: Focus on business logic, validation, and edge cases
2. **Integration Tests**: Test actual database behavior and constraints
3. **Test Isolation**: Each test should be independent and clean up after itself
4. **Meaningful Names**: Test names should clearly describe what is being tested
5. **Arrange-Act-Assert**: Structure tests with clear setup, execution, and verification phases

## CI/CD Integration

The test suite is designed to work in CI/CD environments:

- Unit tests run without external dependencies
- Integration tests can be skipped in environments without database access
- Test database is automatically set up and cleaned up
- All tests are deterministic and repeatable

## Future Improvements

- [ ] Add performance benchmarks
- [ ] Add fuzz testing for validation logic
- [ ] Add transaction testing
- [ ] Add concurrent access testing
- [ ] Add database migration testing
- [ ] Add error scenario testing (network failures, etc.)