package router

import (
    "net/http"

    "github.com/gorilla/mux"
    "github.com/enricod/qaria-be/handlers"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(route.HandlerFunc)
    }

    return router
}

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        handlers.Index,
    },
    Route{
        "StazioniIndex",
        "GET",
        "/api/stazioni",
        handlers.StazioniIndex,
    },
    Route{
        "Misure",
        "GET",
        "/api/misure/{StazioneId}/{Inquinante}",
        handlers.Misure,
    },
}