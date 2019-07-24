package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	vomni "vk/omnibus"

	"github.com/gorilla/mux"
)

func handleCfgRelayInterval(w http.ResponseWriter, point string, cfgCd int) {

	tmpl := "cfgrelayinterval"

	cfgData := []string{"Kuznec"}

	err := tmpls.ExecuteTemplate(w, tmpl, cfgData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlePointCfg(w http.ResponseWriter, r *http.Request) {

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

func handlePointListData(w http.ResponseWriter, r *http.Request) {
	data := allPointData()

	newData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(newData)
}
