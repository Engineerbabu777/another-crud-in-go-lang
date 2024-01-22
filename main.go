
package main;


import (
	"encoding/json"
	"net/http"
	"log"
	
	"github.com/gorilla/mux"
); 


func main(){
	router := mux.NewRouter();

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
         w.Header().Set("Content-Type","application/json");
	     w.WriteHeader(http.StatusOK);
	     json.NewEncoder(w).Encode(map[string]string{"message":"Hello, World!"});


	}).Methods("GET");

	log.Fatal(http.ListenAndServe(":8000", router));
}