package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"./resource"
	"./staticPage"
	"github.com/gorilla/sessions"
	"github.com/gorilla/context"
	"github.com/gorilla/schema"
	"github.com/danzar/fcwServer/common"
	"net/http/httputil"
	"net/url"
	"strings"
	"strconv"
)


var debug = false
var sessionName = "fcw-server"
var camLogin = "member"
var camPassword = "fcw4225"

var remoteAddress = true
var ipAddress = "24.196.196.242:"
//var ipAddress = "192.168.1.{{port}}:"


//94 seems to not responding, could be a firewill error
var ports = []string{"98", "97", "93","94"}
var secretKey = "12345POIUYT"



var store = sessions.NewCookieStore([]byte(secretKey))

var staticPages = template.New("templates")

type FormLogin struct {
	Login  string
	Password string
}


type defaultContext struct {
	Ports []string
}




func main() {


	common.LogDebugData("Main Loading", debug)
	staticPages = staticPage.PopulateStaticPages(getThemeName())

	gorillaRoute := mux.NewRouter()

	//Handle all calls to page calls
	gorillaRoute.HandleFunc("/",serverHandler)
	gorillaRoute.HandleFunc("/{page_alias}",serverHandler)
	gorillaRoute.HandleFunc("/{page_alias}/",serverHandler)
	//This was a test to see how/what is done with the second passed param
//	gorillaRoute.HandleFunc("/{page_alias}/{data}",serverHandler)


	//loop though the ports and create the reverse proxy and handlers for them
	for port := range ports{
		remote, _ := url.Parse("")
		if remoteAddress{
			remote.Host = ipAddress + ports[port]
		}else{
			remote.Host = strings.Replace(ipAddress, "{{port}}",ports[port],-1)  + ports[port]
		}

		remote.Scheme = "http"
		q := remote.Query()
		q.Set("user", camLogin)
		q.Set("pwd", camPassword)
		remote.RawQuery = q.Encode()

		//Create the proxy and send it the handler
		proxy := httputil.NewSingleHostReverseProxy(remote)

		//This has to be http or will not work correct
		http.HandleFunc("/cam"+ strconv.Itoa(port), camHandler(proxy))
	}


	//Handle the assest types
	http.HandleFunc("/img/", resource.ServerResourceFiles)
	http.HandleFunc("/css/", resource.ServerResourceFiles)
	http.HandleFunc("/js/", resource.ServerResourceFiles)


	http.Handle("/", gorillaRoute)
	err := http.ListenAndServe(":8080",context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		common.LogData(err.Error());
	}

}

func serverHandler(w http.ResponseWriter, r *http.Request)  {

        logedIn := false


	//Get Session from store if there is one, if not it will return a new one.
	session , err := getSession(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Deal with Form Items.
	err = r.ParseForm()

	if err != nil {
		// Handle error
	}

	formLogin := new(FormLogin)

	decoder := schema.NewDecoder()
	err = decoder.Decode(formLogin, r.PostForm)

	if err != nil {
		common.LogData("Error logging Form data")
	}

	//if login data is correct then store session.
	if formLogin.Login == camLogin && formLogin.Password == camPassword{
		//We can set the user to logged in on this session once they log in.
		setLoggedIn(session)
		session.Save(r,w)
		logedIn = true
	}



	//get page alias
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]

	//If empty set it to home
	if page_alias == ""{
		page_alias = "home"
	}


	if !logedIn && !getIsLoggedIn(session){
		page_alias = "login"
		common.LogDebugData("Session LoggedIn was false",debug)
	}




	common.LogDebugData("RequestPage: " + page_alias,debug)


	//Get the page from the static pages we have loaded.
	staticPage := staticPages.Lookup(page_alias + ".html")
	if staticPage == nil{
		common.LogDebugData("Page was nil:"+ page_alias,debug)
		staticPage = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}


	//context :=  defaultContext{Ports:ports}

	//Execute the template/Page
	staticPage.Execute(w,nil)
}


func camHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		//Check session to make sure they are logged in to view the cam's
		//This will also prevent someone from just looking stright at the cam link.
		session, _ := getSession(r,sessionName)
		if getIsLoggedIn(session){
			url , _ := url.Parse("")
			// if the Url path has "cam" in it then add the videostream attubute
			if (strings.Contains(r.URL.String(), "cam")){
				// for all the Cam types we are adding the needed item to the request.
				url, _ = url.Parse("videostream.cgi")
			}

			r.URL = url

			p.ServeHTTP(w, r)
			//p.FlushInterval, _ = time.ParseDuration("1h")
		}
		//Send them to the unAuth page then to log in
		http.Redirect(w, r, "/home",403)
		common.LogDebugData("Redirected Login from Cam",debug)

	}
}


func getThemeName() string{
	return "bs3"
}




func getSession( r *http.Request, sessionName string) (*sessions.Session, error) {

	session, err := store.Get(r, sessionName)
	if err != nil {
		return  session, err
	}
	session.Options = &sessions.Options{
		MaxAge:   30,
	}

	return session , nil
}

func getIsLoggedIn(s *sessions.Session) bool {
	if !s.IsNew && s.Values["logged-in"] == true{
		return true
	}
	return false
}

func setLoggedIn(s *sessions.Session)  {
	s.Values["logged-in"] = true
}



//Old items
/*
func buildContextType() defaultContext{
	//--- TODO Update this to a more dynamic passed object. ---
	//Making a static object to send to the pages.
	context := defaultContext{}

	//token , timeNow := createToken(secretKey)

	//context.Token = token
	//	context.TokenTime = timeNow
	context.SiteTitle = "FCW Cam's"
	context.Title = "FCW Cam's"
	context.Message ="The Cam's will be back soon"
	context.ErrorMsgs = ""
	context.SuccessMegs = ""
	//TODO -----------------------------------------------------
	return context
}
*/
