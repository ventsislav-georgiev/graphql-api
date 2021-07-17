## Description
GraphQL API for interacting with different backends.

> Proof of concept implementation for Kinvey and Sitefinity.

## Example query with OData filter and sort

<img width="574" alt="image" src="https://user-images.githubusercontent.com/5616486/126671803-e8d5e1e8-226b-4e7a-a8d8-55be7b9f043e.png">

## Sample queries

(Kinvey) Books with filter and sort and second query for total books:
```vim
curl -s -XPOST -d '{"query": "{ books(filter: \"pages gt 400 and contains(title, JavaScript)\" sort: \"pages desc\") { _id, title, pages, author { name }, contents } booksCount { count } }"}' localhost:8080/graphql | jq 
```

(Kinvey) Book by id:
```vim
curl -s -XPOST -d '{"query": "{ book(id: \"60f8a214ff6d6a0013b2e700\") { title, author { name } } }"}' localhost:8080/graphql | jq
```

(Kinvey) Add book:
```vim
curl -s -XPOST -d '{"query": "mutation { addBook(data: { title: \"Newly added book\" }) { _id } } "}' localhost:8080/graphql
```

(Kinvey) Remove book by id:
```vim
curl -s -XPOST -d '{"query": "mutation { removeBook(id: \"\") { count } }"}' localhost:8080/graphql | jq
```

(Sitefinity) Event with calendar expand:
```vim
curl -s -XPOST -d '{"query": "{ events(filter: \"startswith(Title, '\''XC'\'')\") { Title, Parent { Title } } }"}' localhost:8080/graphql | jq
```

Introspection:
```vim
curl -s -XPOST -d @testdata/curl-introspection.json localhost:8080/graphql | jq
```
