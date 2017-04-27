package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"github.com/enricod/qaria-be/db"
	"github.com/enricod/qaria-be/model"
)


func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func stazioniElenco() []model.Stazione {
	rezzato := model.Stazione{StazioneId:661,
		Nome:"Rezzato",
		Inquinanti:"PM10,NO2,CO",
		Url:"http://www2.arpalombardia.it/sites/QAria/_layouts/15/QAria/DettaglioStazione.aspx?IdStaz=661"}

	milano := model.Stazione{StazioneId:539,
		Nome:"Milano Liguria",
		Inquinanti:"NO2,CO",
		Url:"http://www2.arpalombardia.it/sites/qaria/_layouts/15/qaria/DettaglioStazione.aspx?zona=MI&comune=451&IdStaz=539&isPDV=True"}

	stazioni := []model.Stazione{rezzato, milano }
	return stazioni
}

func stazioneFind(stazioni []model.Stazione, stazId int) *model.Stazione {
	for _, s := range stazioni {
		if s.StazioneId == stazId {
			return &s
		}
	}
	return nil
}

func StazioniIndex(w http.ResponseWriter, r *http.Request) {
	stazioni := stazioniElenco()
	stazioniResp := model.StazioniResp{stazioni}

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

	log.Printf("caricamento misure per stazione %v e inquinante %v", stazioneId, inquinante)
	stazionePtr := stazioneFind(stazioniElenco(), stazioneId)
	if stazionePtr != nil {
		misure, err := dbLeggiMisure(stazionePtr, inquinante)
		if (err != nil) {
			log.Fatal(err)
		}

		var slice []model.Misura
		for _, mPtr := range misure {
			slice = append(slice, *mPtr)
		}
		misureResp := model.MisureResp{slice}

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
func dbLeggiMisure(staz *model.Stazione, inq string) ([]*model.Misura, error) {

	stazId := strconv.Itoa(staz.StazioneId)
	rows, err := db.Db.Query("SELECT id, dataStr, valore, inquinante, stazioneId FROM misura " +
		" WHERE stazioneId=? AND inquinante=? ORDER by dataStr DESC",
		stazId, inq)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]*model.Misura, 0)
	for rows.Next() {
		bk := new(model.Misura)
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

	/*
	m1 := Misura{
		DataMisura:"20170502",
		Inquinante:inq,
		Valore:2.0,
		StazioneId:staz.StazioneId}

	m2 := Misura{
		DataMisura:"20170502",
		Inquinante:inq,
		Valore:2.0,
		StazioneId:staz.StazioneId}
	result := []Misura{m1, m2 }
	return result
	*/
}

