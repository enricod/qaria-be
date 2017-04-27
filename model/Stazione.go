package model

type Stazione struct {
    StazioneId      int    `json:"StazioneId"`
    Nome            string      `json:"Nome"`
    Url       string `json:"Url"`
    Inquinanti string `json:"Inquinanti"`
}

type Stazioni []Stazione

type StazioniResp struct {
    Stazioni []Stazione
}

type Misura struct {
    Id int
    DataMisura string
    Inquinante string
    StazioneId int
    Valore float64
}

type MisureResp struct {
    Misure []Misura
}
