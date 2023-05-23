# Synonyms application

Go application that exposes REST API for storing and retrieving synonyms.
Uses simple in-memory map structure for storing synonyms. 

## Task

Implement a set of basic HTTP endpoint handlers.
The complete system should act as a thesaurus - it should enable users to store and fetch sets of synonyms. 
Please also consider testability, readability, algorithmic complexity, and maintainability!

Your implementation should fulfil the following requirements:
- Endpoint to add new sets of synonyms. For example, adding a pair of synonyms such as “begin” and “start”.
- Endpoint to search for synonyms. In the above example, searching for either “begin” or “start” should return the respective synonym (symmetrical relationship).
- A word may have multiple synonyms, and all should be returned at a user request.
- The solution needs to support concurrent requests in a thread-safe way.
- Make the solution with simple data structures in memory - no persistence/database needed.
- Bonus*: Implement a transitive rule. For example, if “A” is added as a synonym for “B”, and “B” is added as a synonym for “C”, then searching the word “C” should return both “B” and “A”.

## Requirements
- Go 1.20 or Docker 

## Usage

### Run API server

#### 1. With local Go installed
- Run on default port (8080): `make run`

- Run on specific port (9090): `make run-port port=9090`

#### 2. With Docker
- Build image: `make docker-build`
- Run container: `make docker-run`

### Testing

#### 1. Add synonyms pairs

- Add A and B pair:

`curl -X POST -H "Content-Type: application/json" -d '{ "synonyms": ["A", "B"] }' http://localhost:8080/api/v1/synonym`

- Add B and C pair:

`curl -X POST -H "Content-Type: application/json" -d '{ "synonyms": ["B", "C"] }' http://localhost:8080/api/v1/synonym`

#### 2. Retrieve synonyms

- Retrieve synonyms for A: 

`curl http://localhost:8080/api/v1/synonym?word=A`

- Retrieve synonyms for C:

`curl http://localhost:8080/api/v1/synonym?word=C`