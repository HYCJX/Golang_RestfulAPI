# Golang_RestfulAPI

##### To run this golang restful api, use the following command:
##### ```go run main.go $string_of_url_to_crew_json_file```
##### This api server runs on PORT 8080.
<br/>

##### This api handles the following requests:
##### *GET pilot with pilot ID.*
##### *GET all pilots.*
##### *GET available pilots with a query.*
##### *GET flight with pilot ID.*
##### *GET all flights.*
##### *POST flight handled according to pilot availability.*
<br/>

##### The codebase is divided into 3 packages besides the main package:
##### *utils package which provides utilities*
##### *controller package which handles http requests*
##### *model package which connects controller package to database with logic*
