package breamiobeta

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
)

const form = `
<html>
	<head>
		<title>Bream IO Beta Sign-up</title>
		<link rel="stylesheet" type="text/css" href="/static/style.css">
	</head>
	<body>
		<form action="/register" method="post">
			<div class="field">
				<label>Full name:</label>
				<div>
				<input type="text" name="fullname" size="40" placeholder="Bob Allison"></input>
				</div>
			</div>
			<div class="field">
				<label>E-mail:</label>
				<div>
				<input type="email" name="email" size="40" placeholder="awesome_bob94@gmail.com"></input>
				</div>
			</div>
			<div class="field">
				<label>Channel:</label>
				<div>
				<input type="text" name="channel" size="40" placeholder="http://twitch.tv/bob_io" ></input>
				</div>
			</div>
			<div class="field"><input type="submit" value="Sign-up for beta" ></div>
		</form>
	</body>
</html>
`

type Registration struct {
	Name    string
	Email   string
	Channel string
}

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/register", register)
	http.Handle("/thanks", ThankYouPage{})
}

// registrationkey returns the key used for all guestbook entries.
func registrationKey(c appengine.Context) *datastore.Key {
	// The string "default_guestbook" here could be varied to have multiple guestbooks.
	return datastore.NewKey(c, "Registrations", "beta", 0, nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, form)
}

func register(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	reg := Registration{
		Name:    r.FormValue("fullname"),
		Email:   r.FormValue("email"),
		Channel: r.FormValue("channel"),
	}
	fmt.Println(reg)
	key := datastore.NewIncompleteKey(c, "Registration", registrationKey(c))
	_, err := datastore.Put(c, key, &reg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ThankYouPage{reg}.ServeHTTP(w, r)
}

type ThankYouPage struct {
	Registration
}

func (r Registration) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Thank you %s for your registration! We will get back to you when we are ready to start the beta.", r.Name)
}
