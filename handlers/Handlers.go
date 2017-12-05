package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/enricod/qaria-be/db"
	"github.com/enricod/qaria-model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func stazioneFind(stazioni []qariamodel.Stazione, stazId int) *qariamodel.Stazione {
	for _, s := range stazioni {
		if s.StazioneId == stazId {
			return &s
		}
	}
	return nil
}

func StazioniIndex(w http.ResponseWriter, r *http.Request) {
	stazioni := qariamodel.ElencoStazioni()
	stazioniResp := qariamodel.StazioniResp{stazioni}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(stazioniResp); err != nil {
		panic(err)
	}
}

func Misure(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	stazioneId, _ := strconv.Atoi(vars["StazioneId"])
	inquinante := vars["Inquinante"]

	stazionePtr := stazioneFind(qariamodel.ElencoStazioni(), stazioneId)
	if stazionePtr != nil {
		misure, err := dbLeggiMisure(stazionePtr, inquinante)
		if err != nil {
			log.Fatal(err)
		}

		var slice []qariamodel.Misura
		for _, mPtr := range misure {
			slice = append(slice, *mPtr)
		}
		misureResp := qariamodel.MisureResp{slice}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err2 := json.NewEncoder(w).Encode(misureResp); err2 != nil {
			panic(err2)
		}

	}

}

/**
 * leggi dal database le misure di una stazione
 */
func dbLeggiMisure(staz *qariamodel.Stazione, inq string) ([]*qariamodel.Misura, error) {

	stazId := strconv.Itoa(staz.StazioneId)

	log.Printf("DB - caricamento misure per stazione %v e inquinante %v", stazId, inq)
	rows, err := db.Db.Query("SELECT id, dataStr, valore, inquinante, stazioneId FROM misura "+
		" WHERE stazioneId=? AND inquinante=? "+
		" ORDER by dataStr DESC LIMIT 0,20",
		stazId, inq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]*qariamodel.Misura, 0)
	for rows.Next() {
		bk := new(qariamodel.Misura)
		err := rows.Scan(&bk.Id, &bk.DataMisura, &bk.Valore, &bk.Inquinante, &bk.StazioneId)
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}
