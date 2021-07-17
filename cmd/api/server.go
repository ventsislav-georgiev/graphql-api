package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ventsislav-georgiev/graphql-api/pkg/graphql"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers"
	"github.com/ventsislav-georgiev/graphql-api/pkg/helpers/constants"
)

func main() {
	env := checkEnvVars()

	srvAddress := ":8080"
	if env == constants.DevEnvName {
		srvAddress = "localhost" + srvAddress
	}

	http.Handle("/graphql", graphql.NewGraphQLHandler())
	log.Fatal(http.ListenAndServe(srvAddress, nil))
}

func checkEnvVars() string {
	env := flag.String("env", constants.DevEnvName, "env: dev (for local development)")
	kinveyHost := flag.String("kinveyHost", os.Getenv("kinveyHost"), "Host of the Kinvey backend")
	kinveyID := flag.String("kinveyID", os.Getenv("kinveyID"), "KID of the Kinvey App Environment")
	kinveyMasterSecret := flag.String("kinveyMasterSecret", os.Getenv("kinveyMasterSecret"), "MasterSecret of the Kinvey App Environment")
	sitefinityHost := flag.String("sitefinityHost", os.Getenv("sitefinityHost"), "Host of the Sitefinity backend")
	sitefinityToken := flag.String("sitefinityToken", os.Getenv("sitefinityToken"), "Auth token for Sitefinity")
	flag.Parse()

	if helpers.IsNullOrEmpty(kinveyHost) {
		log.Fatal("kinveyHost required")
	}
	if helpers.IsNullOrEmpty(kinveyID) {
		log.Fatal("kinveyID required")
	}
	if helpers.IsNullOrEmpty(kinveyMasterSecret) {
		log.Fatal("kinveyMasterSecret required")
	}
	if helpers.IsNullOrEmpty(sitefinityHost) {
		log.Fatal("sitefinityHost required")
	}
	if helpers.IsNullOrEmpty(sitefinityToken) {
		log.Fatal("sitefinityToken required")
	}

	fmt.Printf("kinveyHost=%s\nkinveyID=%s\nkinveyMasterSecret=%s\nsitefinityHost=%s\nsitefinityToken=%s\n",
		*kinveyHost, *kinveyID, *kinveyMasterSecret, *sitefinityHost, *sitefinityToken)

	os.Setenv("env", *env)
	os.Setenv("kinveyHost", *kinveyHost)
	os.Setenv("kinveyID", *kinveyID)
	os.Setenv("kinveyMasterSecret", *kinveyMasterSecret)
	os.Setenv("sitefinityHost", *sitefinityHost)
	os.Setenv("sitefinityToken", *sitefinityToken)

	return *env
}
