package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"croupier/pkg/auth"
	"croupier/pkg/game"

	"croupier/pkg/controllers"

	cards "croupier/pkg/cards"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/play?room=1&jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImRldmJ5dG9tQGdtYWlsLmNvbSIsImlhdCI6MTUxNjIzOTAyMn0.L3pBFKCdzXuV8AdzYFH73SVTZW3ZUx3KCV0N-chheYg")
}

func main() {

	haa := cards.NewDeck()
	shuffled := cards.Shuffle(haa)
	// for _, c := range shuffled {
	deal := shuffled[0 : len(shuffled)-50]
	fmt.Println(deal)

	fivefinalcards := deal

	fivefinalcards = append(fivefinalcards, shuffled[14])
	fivefinalcards = append(fivefinalcards, shuffled[21])
	fivefinalcards = append(fivefinalcards, shuffled[28])
	fmt.Println(fivefinalcards)
	a := game.CalculateFiveBestCards(fivefinalcards)
	fmt.Printf("%+v \n %+v \n", a, a.ToString())

	router := mux.NewRouter()
	router.Use(auth.JwtHandler)

	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/play", controllers.Play).Methods("GET")

	router.HandleFunc("/api/users/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/users/login", controllers.Authenticate).Methods("POST")

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":"+port, router))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
