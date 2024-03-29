package main

import (
  "fmt"
  "html/template"
  "net/http"
  "strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from GWiki" as the response body.
// Change the signature of the home handler so it is defined as a method against
// *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

  // Check if the current request URL path exactly matches "/". If it doesn't, use
  // the http.NotFound() function to send a 404 response to the client.
  // Importantly, we then return from the handler. If we don't return the handler
  // would keep executing and also write the "Hello from SnippetBox" message.
  if r.URL.Path != "/" {
    app.notFound(w) // Use the notFound() helper
    return
  }

  // Initialize a slice containing the paths to the two files. Note that the
  // home.page.tmpl file must be the *first* file in the slice.
  files := []string{
    "./ui/html/home.page.tmpl",
    "./ui/html/base.layout.tmpl",
    "./ui/html/footer.partial.tmpl",
  }

  // Use the template.ParseFiles() function to read the files and store the
  // templates in a template set. Notice that we can pass the slice of file paths
  // as a variadic parameter?
  ts, err := template.ParseFiles(files...)
  if err != nil {
    // Because the home handler function is now a method against application
    // it can access its fields, including the error logger. We'll write the log
    // message to this instead of the standard logger.
    app.serverError(w, err) // Use the serverError() helper.
    return
  }

  // We then use the Execute() method on the template set to write the template
  // content as the response body. The last parameter to Execute() represents any
  // dynamic data that we want to pass in, which for now we'll leave as nil.
  err = ts.Execute(w, nil)
  if err != nil {
    // Also update the code here to use the error logger from the application
    // struct.
    app.serverError(w, err) // Use the serverError() helper.
  }
}

// Add a showPageSnippet handler function.
// Change the signature of the showSnippet handler so it is defined as a method
// against *application.
func (app *application) showPageSnippet(w http.ResponseWriter, r *http.Request) {
  // Extract the value of the id parameter from the query string and try to
  // convert it to an integer using the strconv.Atoi() function. If it can't
  // be converted to an integer, or the value is less than 1, we return a 404 page
  // not found response.
  id, err := strconv.Atoi(r.URL.Query().Get("id"))
  if err != nil || id < 1 {
    app.notFound(w) // Use the notFound() helper.
    return
  }

  // Use the fmt.Fprintf() function to interpolate the id value with our response
  // and write it to the http.ResponseWriter.
  fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

// Add a createPageSnippet handler function.
// Change the signature of the createSnippet handler so it is defined as a method
// against *application.
func (app *application) createPageSnippet(w http.ResponseWriter, r *http.Request) {
  // Use r.Method to check whether the request is using POST or not. Note that
  // http.MethodPost is a constant equal to the string "POST".
  if r.Method != http.MethodPost {
    // If it's not, use the w.WriteHeader() method to send a 405 status
    // code and the w.Write() method to write a "Method Not Allowed"
    // response body. We then return from the function so that the
    // subsequent code is not executed.
    // Use the Header().Set() method to add an 'Allow: POST' header to the
    // response header map. The first parameter is the header name, and
    // the second parameter is the header value.
    w.Header().Set("Allow", http.MethodPost)

    // Use the http.Error() function to send a 405 status code and "Method Not
    // Allowed" string as the response body.
    app.clientError(w, http.StatusMethodNotAllowed) // Use the clientError() helper.
    return
  }

  w.Write([]byte("Create a new snippet..."))
}

