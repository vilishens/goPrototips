package web

import "net/http"

func handleCfgRelayInterval(w http.ResponseWriter, point string, cfgCd int) {

	tmpl := "cfgrelayinterval"

	cfgData := []string{"Kuznec"}

	err := tmpls.ExecuteTemplate(w, tmpl, cfgData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
