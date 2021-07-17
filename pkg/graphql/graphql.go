package graphql

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/graph-gophers/graphql-go"
	"github.com/ventsislav-georgiev/graphql-api/pkg/backends"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/resolvers"
	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql/schema"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
)

func NewGraphQLHandler() http.Handler {
	graphqlSchemaProvider := &schema.StaticSchemaProvider{}
	graphqlHandler := &graphqlHandler{
		SchemaProvider: graphqlSchemaProvider,
	}

	return graphqlHandler
}

type graphqlHandler struct {
	SchemaProvider schema.SchemaProvider
}

func (h *graphqlHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
		}
	}()

	batch := batchParams{}
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := h.SchemaProvider.GetSchema(r)
	if s == nil {
		http.NotFound(w, r)
		return
	}

	defaultResolversProvider := &resolvers.DefaultResolversProvider{
		KinveyResolver: &resolvers.KinveyResolver{
			Backend: backends.KinveyBackendProvider{
				Host:         os.Getenv("kinveyHost"),
				KinveyID:     os.Getenv("kinveyID"),
				MasterSecret: os.Getenv("kinveyMasterSecret"),
				HttpClient:   &backends.DefaultHttpClient{},
			},
			DataLoaders: &sync.Map{},
		},
		SitefinityResolver: &resolvers.SitefinityResolver{
			Backend: backends.SitefinityBackendProvider{
				Host:       os.Getenv("sitefinityHost"),
				Token:      os.Getenv("sitefinityToken"),
				HttpClient: &backends.DefaultHttpClient{},
			},
			DataLoaders: &sync.Map{},
		},
	}

	schema, err := graphql.ParseSchema(*s, &emptyStruct{},
		graphql.UseDefaultResolvers(defaultResolversProvider),
		graphql.UseDynamicResolvers())

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responses := make([]*graphql.Response, 0)
	fmt.Printf("****************************** Query ******************************\nBatch: %v\n", len(batch))

	for _, params := range batch {
		fmt.Printf("%sVariables: %s\n", params.Query, params.Variables)

		errors := schema.ValidateWithVariables(params.Query, params.Variables)
		if errors != nil {
			errorMessages := make([]string, 0)
			for _, err := range errors {
				msg := err.Error()
				if !helpers.IsEmpty(msg) {
					errorMessages = append(errorMessages, msg)
				}
			}

			var errMsg string
			jsonBytes, err := json.Marshal(errorMessages)
			if err == nil {
				errMsg = string(jsonBytes)
			} else {
				errMsg = strings.Join(errorMessages, "\n")
			}

			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}

		response := schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
		responses = append(responses, response)

		responseJSON, err := json.Marshal(response)
		if err == nil {
			fmt.Printf("Response: %s\n", responseJSON)
		} else {
			fmt.Printf("ResponseErr: %s\n", err.Error())
		}
	}

	responseJSON, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
