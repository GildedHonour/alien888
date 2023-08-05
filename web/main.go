package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	helper "alien888/helper"
)

const (
	defaultConfigFileName = "config.yaml"
	sqlLiteDbFileName     = "alien888.db"
)

var (
	dbConn *sql.DB
	cfg    *Config
	// cookieStore = sessions.NewCookieStore([]byte(cfg.SessionKey))
	cookieStore *sessions.CookieStore
)

func main() {
	var err error

	err = godotenv.Load()
	helper.CheckErr(err)
	cookieStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/gallery", galleryHandler)
	router.HandleFunc("/wardrobe", wardrobeHandler)
	router.HandleFunc("/wallet", walletHandler)
	router.HandleFunc("/current_wallet", currentWalletHandler)
	router.HandleFunc("/current_custom_avatar", currentCustomAvatarHandler)
	router.HandleFunc("/whitelist", whiteListHandler)
	router.HandleFunc("/erc1155-tokens/{id}.json", tokensERC1155MetadataHandler)
	router.HandleFunc("/erc721-tokens/{id}.json", tokensERC721MetadataHandler)
	//TODO
	router.HandleFunc("/erc721-tokens/{id}", tokensERC721MetadataHandler)

	cfg, err = loadConfig(defaultConfigFileName)
	helper.CheckErr(err)

	dbConn, err = sql.Open("sqlite3", sqlLiteDbFileName)
	helper.CheckErr(err)
	defer dbConn.Close()

	fmt.Printf("server on the port %d is running...\n", cfg.Port)
	portStr := ":" + strconv.Itoa(cfg.Port)

	// fSrv := http.FileServer(http.Dir("./assets/"))
	// http.Handle("/assets/", http.StripPrefix("/assets/", fSrv))
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// err = http.ListenAndServe(portStr, nil)
	err = http.ListenAndServe(portStr, router)

	helper.CheckErr(err)
}
