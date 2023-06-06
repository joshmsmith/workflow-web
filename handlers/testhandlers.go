package handlers

import (
  "log"
  "net/http"
  "webapp/utils"
)

type TestData struct {
    Id      int
    Number  int
    Name    string
    Float   float64
}

/* Bootstrap test handler */
func Bootstrap (w http.ResponseWriter, r *http.Request) {

  testdata := returnTestData()
  log.Println("Bootstrap testhandler: called")
  utils.RawRender(w, "templates/Bootstrap.html", testdata)
}

/* Initialise some test data */
func returnTestData () (tdata TestData) {
  return TestData {
      Id:     12,
      Number: 237782,
      Name:   "Test Name Here",
      Float:  float64(2789.23),
    }
}
