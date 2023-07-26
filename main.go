package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type User struct {
	ID     int    `json:"id"`
	Nom    string `json:"nom"`
	Surname string `json:"surname"`
	Number int `json:"number"`
}

type Tontine struct {
	ID  int    `json:"id"`
	Nom string `json:"nom"`
}

var db *sql.DB

func main() {
	// Paramètres de connexion à la base de données MySQL
	dbUser := "root"
	dbPass := ""
	dbHost := "localhost"
	dbPort := "3306"
	dbName := "go"

	// Chaîne de connexion à la base de données
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// Connexion à la base de données
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connexion à la base de données MySQL réussie !")

	// Initialisation du routeur
	router := mux.NewRouter()

	// Routes de l'API
	router.HandleFunc("/table/create/{name}", createTable).Methods("POST")
	router.HandleFunc("/table/add/{tablename}", addElement).Methods("POST")
	router.HandleFunc("/table/{name}", getTableContent).Methods("GET")
	router.HandleFunc("/tables", getAllTables).Methods("GET")
	router.HandleFunc("/user/info/{name}", getUserInfo).Methods("GET")

	log.Println("Démarrage du serveur sur le port 3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func createTable(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tablename := params["name"]

	// Vérification de la connexion à la base de données
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Table créée au nom de:", tablename)

	createTableQuery := fmt.Sprintf("INSERT INTO tontine1 (nom) VALUES ('%s');", tablename)
	_, err = db.Exec(createTableQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	fmt.Println("Réussi")

	w.WriteHeader(http.StatusCreated)
}

func addElement(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    tablename := params["tablename"]

    // Lire les données du corps de la requête
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Vérification de la connexion à la base de données
    err = db.Ping()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Insertion des données dans la table
    insertQuery := fmt.Sprintf("INSERT INTO %s (nom, surname, number) VALUES (?, ?, ?)", tablename)
    _, err = db.Exec(insertQuery, user.Nom, user.Surname, user.Number)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Println("Réussi")

    w.WriteHeader(http.StatusCreated)
}

func getTableContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tablename := params["name"]

	// Vérification de la connexion à la base de données
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Récupération de toutes les données de la table
	selectQuery := fmt.Sprintf("SELECT U.id, U.nom, U.surname, U.number FROM user U JOIN faire F ON U.id = F.id JOIN tontine1 T ON F.id_tontine = T.id_tontine WHERE T.nom = '%s';", tablename)
	rows, err := db.Query(selectQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID,&user.Nom, &user.Surname, &user.Number)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Envoi des données en tant que réponse JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getAllTables(w http.ResponseWriter, r *http.Request) {
	// Vérification de la connexion à la base de données
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	showTablesQuery := "SHOW TABLES"
	rows, err := db.Query(showTablesQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		err := rows.Scan(&table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tables = append(tables, table)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(tables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// Vérification de la connexion à la base de données
	err := db.Ping()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	selectQuery := fmt.Sprintf("SELECT T.nom FROM tontine1 T JOIN faire F ON T.id_tontine = F.id_tontine JOIN user U ON F.id = U.id WHERE U.nom = '%s';", name)
	rows, err := db.Query(selectQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tontines []Tontine
	for rows.Next() {
		var tontine Tontine
		err := rows.Scan(&tontine.Nom)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tontines = append(tontines, tontine)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tontines)
}
