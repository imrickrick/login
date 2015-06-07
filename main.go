package main

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"github.com/gorilla/mux"
	//"html"
	//"log"
	//"appengine/user"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func init() {

	router := mux.NewRouter().StrictSlash(true)

	router.Handle("/static/CSS/", http.StripPrefix("/static/CSS/", http.FileServer(http.Dir("static/CSS"))))

	router.HandleFunc("/api/v1/farm/{farmID}", createHandler)
	router.HandleFunc("/api/v1/farm/", farmIndexHandler)
	router.HandleFunc("/api/v1/recipe/{recipeID}", recipeValHandler)
	router.HandleFunc("/api/v1/recipe/", recipeIndexHandler)
	router.HandleFunc("/api/v1/", apiIndexHandler)
	router.HandleFunc("/master/", masterListHandler)
	router.HandleFunc("/createrecipe/", createRecipeHandler)
	router.HandleFunc("/main", mainHandler)
	router.HandleFunc("/", rootHandler)

	http.Handle("/", router)
}

type users struct {
	UserName       string
	Password       string
	DateRegistered time.Time
}

type variables struct {
	ID           int       //`json:"id"`
	VariableName string    //`json:"variable name"`
	Description  string    //`json:"description"`
	Unit         string    //`json:"unit"`
	Validation   string    //`json:"validation"`
	DateAdded    time.Time //`json:"date added"`
}

type newVariables struct {
	ID           int    //`json:"id"`
	VariableName string //`json:"variable name"`
	Unit         string //`json:"unit"`
}

type values struct {
	TimeCreated time.Time `json:"time created"`
	TimeStored  time.Time `json:"time stored"`
	DeviceUse   string    `json:"device use"`
	VariableID  string    `json:"variable id"`
	FarmID      string    `json:"farm id"`
	Token       string    `json:"token"`
	Value       string    `json:"value"`
}

type recipes struct {
	Recipeid    string    //`json:"id"`
	Name        string    //`json:"name"`
	Description string    //`json:"description"`
	RecipeVars  string    //`json:"recipe variables"`
	DateCreated time.Time //`json:"date created"`
}
type RecipeV map[string]string

type recipeVars struct {
	ID       int    //`json:"id"`
	Name     string //`json:"name"`
	DataType string //`json:"data type"`
}

func farmIndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		c := appengine.NewContext(r)
		q := datastore.NewQuery("tblDATA")

		var value []values
		if _, err := q.GetAll(c, &value); err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		response, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}
		fmt.Fprintln(w, string(response))

	} else {
		fmt.Fprintln(w, "ACCESS DENIED!")
		return
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)
		farmID := vars["farmID"]

		// c := appengine.NewContext(r)
		// q := datastore.NewQuery("tblDATA").Filter("FarmID =", farmID) //.Order("-DateAdded")
		// for t := q.Run(c); ; {
		// 	var x values
		// 	_, err := t.Next(&x)
		// 	if err == datastore.Done {
		// 		break
		// 	}
		// 	if err != nil {
		// 		c.Infof("Error in querying datastore: %v", err)
		// 		errorHandler(w, r, http.StatusInternalServerError, "")
		// 	}

		// 	response, err := json.MarshalIndent(x, "", "  ")
		// 	if err != nil {
		// 		errorHandler(w, r, http.StatusInternalServerError, "")
		// 	}

		// 	fmt.Fprintln(w, string(response))
		// } //end for going through q.Run
		//test commit
		c := appengine.NewContext(r)
		q := datastore.NewQuery("tblDATA").Filter("FarmID =", farmID) //.Order("-DateAdded")
		var value []values
		if _, err := q.GetAll(c, &value); err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		response, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}
		fmt.Fprintln(w, string(response))

	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		vars := mux.Vars(r)
		farmID := vars["farmID"]

		decoder := json.NewDecoder(r.Body)

		var v values
		err := decoder.Decode(&v)
		if err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		vals := values{
			TimeCreated: time.Now(),
			TimeStored:  time.Now(),
			VariableID:  farmID,
			DeviceUse:   "M",
			FarmID:      farmID,
			Token:       r.Header.Get("X-Farm-Token"),
			Value:       v.Value,
		}

		c := appengine.NewContext(r)

		key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "tblDATA", nil), &vals)
		if err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusNotFound, "")
			return
		}
		fmt.Print(key)
		response, err := json.MarshalIndent(vals, "", "  ")
		if err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
		}
		fmt.Fprintln(w, string(response))

	}
}

func createRecipeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createrecipe/" {
		errorHandler(w, r, http.StatusNotFound, "")
		//return
	}
	if r.Method == "GET" {

		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/createRecipe.gtpl",
		))

		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	} else {

		c := appengine.NewContext(r)
		q, err := datastore.NewQuery("tblRecipes").KeysOnly().Count(c) //.Filter("Name =", recipeID) //.Order("-DateAdded")
		if err != nil {
			fmt.Fprintf(w, `count err: %s`, err)
			return
		}
		//i := strconv.Itoa(q)
		recipe := recipes{
			Recipeid:    "recipe" + strconv.Itoa(q),
			Name:        r.FormValue("recipename"),
			Description: r.FormValue("description"),
			DateCreated: time.Now(),
			//RecipeVars:  r.FormValue("var1") + "," + r.FormValue("var2") + "," + r.FormValue("var3") + "," + r.FormValue("var4") + "," + r.FormValue("var5"),

			//RecipeVars:  recipeVars,
			//RecipeVars: "{\"id\":\"" + r.FormValue("var1") + "\"}, {\"id\":\"" + r.FormValue("var2") + "\"}, {\"id\":\"" + r.FormValue("var3") + "\"}, {\"id\":\"" + r.FormValue("var4") + "\"}, {\"id\":\"" + r.FormValue("var5") + "\"}",
			//RecipeVars: "\"" + r.FormValue("var1") + "\", \"" + r.FormValue("var2") + "\", \"" + r.FormValue("var3") + "\", \"" + r.FormValue("var4") + "\", \"" + r.FormValue("var5") + "\"",
			//RecipeVars: r.FormValue("var1"),
		}

		key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "tblRecipes", nil), &recipe)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}
		fmt.Print(key)
		http.Redirect(w, r, "/createrecipe/", http.StatusFound)
	}
}

func recipeIndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/v1/recipe/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	if r.Method == "GET" {

		c := appengine.NewContext(r)
		q := datastore.NewQuery("tblRecipes").Order("Recipeid")

		var recipe []recipes
		if _, err := q.GetAll(c, &recipe); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		response, err := json.MarshalIndent(recipe, "", "  ")
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}
		fmt.Fprintln(w, string(response))

	} else {
		fmt.Fprint(w, "{\n  'method':'post',\n  'details':'access denied'\n}\n")
		return
	}
}

func recipeValHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	if r.Method == "GET" {
		vars := mux.Vars(r)
		recipeID := vars["recipeID"]
		q := datastore.NewQuery("tblRecipes").Filter("Recipeid =", recipeID)
		var recipe []recipes
		if _, err := q.GetAll(c, &recipe); err != nil {
			panic(err)
			//errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		// response, err := json.MarshalIndent(recipe, "", "  ")
		// if err != nil {
		// 	errorHandler(w, r, http.StatusInternalServerError, "")
		// 	return
		// }
		// fmt.Fprintln(w, string(response))

		for _, rec := range recipe {
			// fmt.Fprintln(w, rec.Recipeid)
			// fmt.Fprintln(w, rec.Name)
			// fmt.Fprintln(w, rec.Description)
			// fmt.Fprintln(w, rec.DateCreated)
			// fmt.Fprintln(w, rec.RecipeVars)

			data := rec.RecipeVars
			res := strings.Split(data, ",")
			//fmt.Fprintln(w, res)
			//dc := strconv.Itoa(rec.DateCreated)
			fmt.Fprintln(w, "{\"recipe id\": \""+rec.Recipeid+"\",")
			fmt.Fprintln(w, "\"name\": \""+rec.Name+"\",")
			fmt.Fprintln(w, "\"description\": \""+rec.Description+"\",")
			//fmt.Fprintln(w, "\"date created\" :\""+dc+"\"")
			fmt.Fprintln(w, "\"recipe vars\": [")

			for count, i := range res {
				var s string = i
				id, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					panic(err)
				}
				q := datastore.NewQuery("tblVariables").Filter("ID =", id)

				var vars []variables
				if _, err := q.GetAll(c, &vars); err != nil {
					errorHandler(w, r, http.StatusInternalServerError, "")
					return
				}
				for _, vars2 := range vars {
					fmt.Fprint(w, "{ \"id\" :\"", vars2.ID, "\",")
					fmt.Fprint(w, "\"name\" :\"", vars2.VariableName, "\",")
					fmt.Fprint(w, "\"validation\" :\"", vars2.Validation, "\",")
					fmt.Fprint(w, "\"unit\" :\"", vars2.Unit, "\"}")
					if count+1 != len(res) {
						fmt.Fprint(w, ",\n")
					}
				}

			}
			fmt.Fprintln(w, "]}")

		}

	} else {
		fmt.Fprint(w, "{\n  'method':'post',\n  'details':'access denied'\n}\n")
		return
	}
}

func apiIndexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/api/v1/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	if r.Method == "GET" {
		fmt.Fprintln(w, "BUKID UTILITY API VERSION 1.0")
		fmt.Fprintln(w, time.Now())

	} else {
		fmt.Fprint(w, "{\n  'method':'post',\n  'details':'access denied'\n}\n")
		return
	}

}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	us := users{
		UserName:       r.FormValue("username"),
		Password:       r.FormValue("password"),
		DateRegistered: time.Now(),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "tblUsers", nil), &us)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, "")
		return
	}

	// var e2 users
	// if err = datastore.Get(c, key, &e2); err != nil {
	// 	errorHandler(w, r, http.StatusInternalServerError, "")
	// 	return
	// }

	//fmt.Fprintf(w, "Stored and retrieved the Employee named %q", us.UserName)

	if r.URL.Path != "/main" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}
	page := template.Must(template.ParseFiles(
		"static/_base.gtpl",
		"static/main.gtpl",
	))

	if err := page.Execute(w, us); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}
func masterListHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/master/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}
	c := appengine.NewContext(r)

	if r.Method == "GET" {
		q := datastore.NewQuery("tblVariables") //.Filter("VariableName =", "PH")

		var vars []variables
		if _, err := q.GetAll(c, &vars); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/variables.gtpl",
		))
		if err := page.Execute(w, vars); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	} else {
		c := appengine.NewContext(r)
		q, err := datastore.NewQuery("tblVariables").KeysOnly().Count(c) //.Filter("Name =", recipeID) //.Order("-DateAdded")
		if err != nil {
			fmt.Fprintf(w, `count err: %s`, err)
			return
		}

		vars := variables{
			ID:           q + 1,
			VariableName: r.FormValue("variableName"),
			Description:  r.FormValue("variableDescription"),
			Unit:         r.FormValue("variableUnit"),
			Validation:   r.FormValue("variableValidation"),
			DateAdded:    time.Now(),
		}

		key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "tblVariables", nil), &vars)
		if err != nil {
			panic(key)
			errorHandler(w, r, http.StatusInternalServerError, "")
			return
		}

		http.Redirect(w, r, "/master/", http.StatusFound)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "")
		return
	}

	page := template.Must(template.ParseFiles(
		"static/_base.gtpl",
		"static/index.gtpl",
	))
	if err := page.Execute(w, nil); err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

}

func errorHandler(w http.ResponseWriter, r *http.Request, status int, err string) {
	w.WriteHeader(status)
	switch status {

	case http.StatusNotFound:
		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/404.gtpl",
		))
		if err := page.Execute(w, nil); err != nil {
			errorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}

	case http.StatusInternalServerError:
		page := template.Must(template.ParseFiles(
			"static/_base.gtpl",
			"static/500.gtpl",
		))
		if err := page.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
