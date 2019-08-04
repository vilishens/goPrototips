package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	vomni "vk/omnibus"
	vpointrun "vk/run/pointrun"

	"github.com/gorilla/mux"
)

func handleCfgRelayIntervalSSS(w http.ResponseWriter, point string, cfgCd int) {

	tmpl := "cfgrelayinterval"

	cfgData := []string{"Kuznec"}

	err := tmpls.ExecuteTemplate(w, tmpl, cfgData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPointCfg(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	point := vars["point"]
	cfg := vars["cfg"]

	//	var data interface{}

	fmt.Println("SVIRIDOFF-MIKAHIL_ABALMASOFF", cfg)

	/*
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
	*/

	switch cfg {
	case strconv.Itoa(vomni.CfgTypeRelayInterval):
	default:
		err := fmt.Errorf("Don't have code to handle configuration %q of the point %q", cfg, point)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := pointData(point)
	newData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(newData)
}

func handlePointCfg(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	todo := vars["todo"]
	point := vars["point"]
	cfg := vars["cfg"]

	fmt.Println("viorika-VISKOPOLEANU", cfg, "Point", point, "TODO", todo)

	switch strings.ToUpper(todo) {
	case "GET":
		getPointCfg(w, r)
		return

	case "LOADINP", "LOADDEFAULT", "LOADSAVED":
		cfgCd, _ := strconv.Atoi(cfg)
		send2Point(w, r, point, vomni.PointCmdLoadCfgIntoPoint|cfgCd)
		//		loadCfgData
		return

	case "SAVECFG":
		cfgCd, _ := strconv.Atoi(cfg)
		send2Point(w, r, point, vomni.PointCmdSaveCfg|cfgCd)
		return

	default:
		log.Fatal("Don't know what to do with %q", todo)

	}
}

/*
		//rescanPoint(point)
		//tmplStr := "pointcfg"

		//cfg, _ := strconv.Atoi(subtype)

		fmt.Printf("Kods %q subtype %q\n", strconv.Itoa(vomni.CfgTypeRelayInterval), point)

		//data := pointData(point)

		err = tmpls.ExecuteTemplate(w, tmplStr, point)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	case "LOADCFG", "SAVECFG":
		/*
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal(body, &data)
			if err != nil {
				panic(err.Error())
			}
		* /
	case "FREEZE", "UNFREEZE", "LOADDEFAULTCFG", "LOADSAVEDCFG":
	default:
		log.Fatal(fmt.Sprintf("===> Don't know what to do with the point %q configuration %q )", point, cfg))
	}

	responseOK(w)
}
*/

//####################################

func loadCfgData(w http.ResponseWriter, r *http.Request, point string, cfg int) {

	var data interface{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	vpointrun.WebSent(vomni.PointCmdLoadCfgIntoPoint|cfg, point, data)
}

func send2Point(w http.ResponseWriter, r *http.Request, point string, cmd int) {

	var data interface{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	vpointrun.WebSent(cmd, point, data)
}

//#####################################

func handlePointListData(w http.ResponseWriter, r *http.Request) {
	data := allPointData()

	newData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(newData)
}
