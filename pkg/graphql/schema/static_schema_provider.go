package schema

import "net/http"

type StaticSchemaProvider struct{}

func (sp *StaticSchemaProvider) GetSchema(r *http.Request) *string {
	s := `
	directive @backend(product: String, collection: String, method: String, endpoint: String) on OBJECT | FIELD_DEFINITION
	directive @connection(primaryKey: String, foreignKey: String) on FIELD_DEFINITION

	type Query {
		book(id: ID!): Book
		books(filter: String, sort: String): [Book]
		booksCount(filter: String): KinveyCount @backend(product: kinvey, endpoint: "/appdata/:kid/books/_count")

		author(id: ID!): Author
		authors(filter: String, sort: String): [Author]
		authorsCount(filter: String): KinveyCount @backend(product: kinvey, endpoint: "/appdata/:kid/authors/_count")

		newsItem(id: ID!): NewsItem
		newsItems(filter: String, sort: String): [NewsItem]

		blog(id: ID!): Blog
		blogs(filter: String, sort: String): [Blog]

		blogPost(id: ID!): BlogPost
		blogPosts(filter: String, sort: String): [BlogPost]

		event(id: ID!): Event
		events(filter: String, sort: String): [Event]

		calendar(id: ID!): Calendar
		calendars(filter: String, sort: String): [Calendar]
	}

	type Mutation {
		addBook(data: BookInput!): Book @backend(product: kinvey, collection: books, method: POST)
		addBooks(data: [BookInput]!): KinveyCount @backend(product: kinvey, collection: books, method: POST)
		updateBook(id: ID!, data: BookInput!): Book @backend(product: kinvey, collection: books, method: PUT)
		removeBook(id: ID!): KinveyCount @backend(product: kinvey, collection: books, method: DELETE)
		removeBooks(filter: String): KinveyCount @backend(product: kinvey, collection: books, method: DELETE)

		addAuthor(data: AuthorInput!): Author @backend(product: kinvey, collection: authors, method: POST)
		addAuthors(data: [AuthorInput]!): KinveyCount @backend(product: kinvey, collection: authors, method: POST)
		updateAuthor(id: ID!, data: AuthorInput!): Author @backend(product: kinvey, collection: authors, method: PUT)
		removeAuthor(id: ID!): KinveyCount @backend(product: kinvey, collection: authors, method: DELETE)
		removeAuthors(filter: String): KinveyCount @backend(product: kinvey, collection: authors, method: DELETE)

		addNewsItem(data: NewsItemInput!): NewsItem @backend(product: sitefinity, collection: newsitems, method: POST)
		updateNewsItem(id: ID!, data: NewsItemInput!): Boolean @backend(product: sitefinity, collection: newsitems, method: PATCH)
		removeNewsItem(id: ID!): Boolean @backend(product: sitefinity, collection: newsitems, method: DELETE)

		addBlog(data: BlogInput!): Blog @backend(product: sitefinity, collection: blogs, method: POST)
		updateBlog(id: ID!, data: BlogInput!): Boolean @backend(product: sitefinity, collection: blogs, method: PATCH)
		removeBlog(id: ID!): Boolean @backend(product: sitefinity, collection: blogs, method: DELETE)

		addBlogPost(data: BlogPostInput!): BlogPost @backend(product: sitefinity, collection: blogposts, method: POST)
		updateBlogPost(id: ID!, data: BlogPostInput!): Boolean @backend(product: sitefinity, collection: blogposts, method: PATCH)
		removeBlogPost(id: ID!): Boolean @backend(product: sitefinity, collection: blogposts, method: DELETE)

		addEvent(data: EventInput!): Event @backend(product: sitefinity, collection: events, method: POST)
		updateEvent(id: ID!, data: EventInput!): Boolean @backend(product: sitefinity, collection: events, method: PATCH)
		removeEvent(id: ID!): Boolean @backend(product: sitefinity, collection: events, method: DELETE)

		addCalendar(data: CalendarInput!): Calendar @backend(product: sitefinity, collection: calendars, method: POST)
		updateCalendar(id: ID!, data: CalendarInput!): Boolean @backend(product: sitefinity, collection: calendars, method: PATCH)
		removeCalendar(id: ID!): Boolean @backend(product: sitefinity, collection: calendars, method: DELETE)
	}

	type Book @backend(product: kinvey, collection: books) {
		_id: ID
		title: String
		authorId: String
		author: Author @connection(primaryKey: authorId)
		contents: [String]
		pages: Int
	}

	input BookInput {
		title: String
		authorId: String
		contents: [String]
		pages: Int
	}

	type Author @backend(product: kinvey, collection: authors) {
		_id: ID
		name: String
	}

	input AuthorInput {
		name: String
	}

	type KinveyCount {
		count: Int
	}

	type NewsItem @backend(product: sitefinity, collection: newsitems)  {
		Id: ID
		Title: String
	}

	input NewsItemInput {
		Title: String
	}

	type Blog @backend(product: sitefinity, collection: blogs)  {
		Id: ID
		Title: String
		UrlName: String
	}

	input BlogInput {
		Title: String
	}

	type BlogPost @backend(product: sitefinity, collection: blogposts)  {
		Id: ID
		Title: String
		ParentId: String
		Parent: Blog @connection(primaryKey: ParentId)
	}

	input BlogPostInput {
		Title: String
		ParentId: String
	}

	type Event @backend(product: sitefinity, collection: events) {
		Id: ID
		Title: String
		ParentId: String
		Parent: Calendar @connection(primaryKey: ParentId)
	}

	input EventInput {
		Title: String
		ParentId: String
	}

	type Calendar @backend(product: sitefinity, collection: calendars)  {
		Id: ID
		Title: String
	}

	input CalendarInput {
		Title: String
	}
`
	return &s
}
