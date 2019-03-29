package web

import (
	"fmt"
	"net/http"

	//	vscanip "vk/code/net/scanpoints"

	//	xrun "vk/run/a_runningpoints"

	"github.com/gorilla/mux"
)

var rtr = mux.NewRouter()

func setMux() {

	fmt.Println("Jāatjauno pageStatic!!!!!")

	rtr.HandleFunc("/about", pageAbout) //
	//	rtr.HandleFunc("/login", pageLogin) //
	rtr.HandleFunc("/", pageHome)
	rtr.HandleFunc("/home", pageHome)

	//	rtr.HandleFunc("/pointlist", tmplPointList)
	//	rtr.HandleFunc("/pointlist/data", tmplPointListData)

	//	rtr.HandleFunc("/point/{point}/{todo}", pointToDo)
	//	rtr.HandleFunc("/point/handlecfg/{point}/{todo}", handleCfg)
	//	rtr.HandleFunc("/station/{todo}", handleStation)

	//	http.HandleFunc(vomni.WebPrefix, StaticFile) // usually read files required for templates css, js, ...

	http.Handle("/", rtr)
}

func pageAbout(w http.ResponseWriter, r *http.Request) {
	pageStatic("about", w, r)
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	pageStatic("home", w, r)
}

func pageStatic(tmpl string, w http.ResponseWriter, r *http.Request) {

	//var data interface{}

	err := error(nil)
	//err := tmpls.ExecuteTemplate(w, tmpl, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//----------------------------------------------------------------------------->

/*
func tmplPointList(w http.ResponseWriter, r *http.Request) {
	this_tmpl := "pointlist"

	fmt.Println("Kiriloff ", this_tmpl, " polina")

	err := tmpls.ExecuteTemplate(w, this_tmpl, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func tmplPointListData(w http.ResponseWriter, r *http.Request) {
	//this_tmpl := "pointlist"

	fmt.Println("Maxim ", " polina")

	//	err := tmpls.ExecuteTemplate(w, this_tmpl, r)

	data := pointList()

	a, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(a)
}

func pageLogin(w http.ResponseWriter, r *http.Request) {
	pageStatic("login", w, r)
}

func pageHome(w http.ResponseWriter, r *http.Request) {
	pageStatic("home", w, r)
}


func pageStatic(tmpl string, w http.ResponseWriter, r *http.Request) {

	var data interface{}

	err := tmpls.ExecuteTemplate(w, tmpl, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func pointToDo(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	todo := vars["todo"]
	point := vars["point"]

	var err error
	var data interface{}

	switch todo {
	case "showcfg":
		tmplStr := "pointcfg"
		data = pointCfg(point)

		refl := reflect.ValueOf(data)

		zType := refl.FieldByName("Type")

		switch zType.Int() {
		case vomni.PointTypeRelayOnOffInterval:
			tmplStr = "cfgrelayonoffinterval"
		default:
			tmplStr = "pointcfg"
		}

		err = tmpls.ExecuteTemplate(w, tmplStr, point)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "getpointcfg":
		data := pointCfg(point)

		a, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(a)

	default:
		http.NotFound(w, r)
	}
}

func handleCfg(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	todo := strings.ToUpper(vars["todo"])
	point := vars["point"]
	var data interface{}

	switch todo {
	case "LOADCFG", "SAVECFG":
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			panic(err.Error())
		}
	case "FREEZE", "UNFREEZE", "LOADDEFAULTCFG", "LOADSAVEDCFG":
	default:
		log.Fatal(fmt.Sprintf("===> Don't know what to do with \"%s\"", todo))
	}

	responseOK(w)
	xrun.ReceivedWebMsg(point, todo, data)
}

func handleStation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todo := strings.ToUpper(vars["todo"])

	switch todo {
	case "SCANIP":

		chDone := make(chan bool)
		chErr := make(chan error)

		go vscanip.ScanPoints(chDone, chErr)

		responseOK(w)

		select {
		case <-chDone:
		case err := <-chErr:
			vomni.LogErr.Println(vutils.ErrFuncLine(err))
		}
	case "RESTART":
		vomni.RootDone <- vomni.DoneRestart
		responseOK(w)

	case "EXIT":
		vomni.RootDone <- vomni.DoneStop

	default:
		log.Fatal(fmt.Sprintf("===> Don't know what to do with \"%s\"", todo))
	}

	responseOK(w)
}

func responseOK(w http.ResponseWriter) {
	type resp struct {
		RC string
	}

	a, err := json.Marshal(resp{RC: "OK"})
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(a)
}

func StaticFile(w http.ResponseWriter, req *http.Request) {
	staticFile := req.URL.Path[len(vomni.WebPrefix):]

	if len(staticFile) != 0 {
		f, err := http.Dir(vomni.WebStaticPath).Open(staticFile)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, staticFile, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}
*/
