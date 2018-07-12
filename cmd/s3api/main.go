package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
)

/*
  fakeUsers = [{id: 1, firstName: 'Dhiraj', lastName: 'Ray', email: 'dhiraj@gmail.com'},
    {id: 1, firstName: 'Tom', lastName: 'Jac', email: 'Tom@gmail.com'},
    {id: 1, firstName: 'Hary', lastName: 'Pan', email: 'hary@gmail.com'},
    {id: 1, firstName: 'praks', lastName: 'pb', email: 'praks@gmail.com'},
  ];
*/

// User defines a okta user
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Diag struct {
	Headers map[string][]string `json:"headers"`
	Uri string
}

func main() {
	users := []User{{FirstName: "Dhiraj", LastName: "Ray", Email: "dhiraj@gmail.com"}}
	usersJSON, _ := json.Marshal(users)
	var RootCmd = &cobra.Command{
		Use:   "s3api",
		Short: "Simple HTTP server read and write to an s3 bucket",
		Run: func(cmd *cobra.Command, args []string) {
			println("Starting server on port 8888")
			http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
				fmt.Printf("Headers %v", r.Header)
				w.Header().Add("Content-Type", "application/json")
				w.Write(usersJSON)
			})
			http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
			})
			http.HandleFunc("/diag", func(w http.ResponseWriter, r *http.Request) {
				response := &Diag{
					Headers: r.Header,
					Uri: r.RequestURI,
				}
				json, _ := json.Marshal(response)
				w.Header().Add("Content-Type", "application/json")
				w.Write(json)
			})

			log.Fatal(http.ListenAndServe(":8888", nil))

		},
	}

	if err := RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
