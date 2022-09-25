package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	route.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	route.PathPrefix("/icon/").Handler(http.StripPrefix("/icon/", http.FileServer(http.Dir("./icon"))))

	route.HandleFunc("/home", home).Methods("GET")
	route.HandleFunc("/add_myproject", addMyProject).Methods("GET")
	route.HandleFunc("/contact_me", contactMe).Methods("GET")
	route.HandleFunc("/detail_project/{index}", detailProject).Methods("GET")
	route.HandleFunc("/add_myproject", ambilData).Methods("POST")
	route.HandleFunc("/delete_project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/halaman_edit/{index}", halamanEdit).Methods("GET")
	route.HandleFunc("/submit_halaman_edit/{indexedit}", submitHalamanEdit).Methods("POST")

	fmt.Println("server running on port 80")
	http.ListenAndServe("localhost:80", route)

}

type Project struct {
	NamaProject string
	Description string
	Durasi      string
	NodeJS      string
	ReactJS     string
	JavaScript  string
	SocketIO    string
}

var dataProject = []Project{
	{
		NamaProject: "Test Project",
		Description: "IniDescripsinya",
	},
}

func ambilData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	var projectName = r.PostForm.Get("project-name")
	var description = r.PostForm.Get("description")
	var startdate = r.PostForm.Get("start-date")
	var enddate = r.PostForm.Get("end-date")
	var nodejs = r.PostForm.Get("tech")
	var reactjs = r.PostForm.Get("tech2")
	var javascript = r.PostForm.Get("tech3")
	var socketio = r.PostForm.Get("tech4")

	Format := "2006-01-02"
	var sdate, _ = time.Parse(Format, startdate)
	var edate, _ = time.Parse(Format, enddate)
	durasiDalamJam := edate.Sub(sdate).Hours()

	durasiDalamHari := durasiDalamJam / 24
	durasiDalamBulan := durasiDalamHari / 30
	durasiDalamTahun := durasiDalamBulan / 12

	var durasi string
	var hari, _ float64 = math.Modf(durasiDalamHari)
	var bulan, _ float64 = math.Modf(durasiDalamBulan)
	var tahun, _ float64 = math.Modf(durasiDalamTahun)

	if tahun > 0 {
		durasi = "durasi: " + strconv.FormatFloat(tahun, 'f', 0, 64) + " Tahun"
	} else if bulan > 0 {
		durasi = "durasi: " + strconv.FormatFloat(bulan, 'f', 0, 64) + " Bulan"
	} else if hari > 0 {
		durasi = "durasi: " + strconv.FormatFloat(hari, 'f', 0, 64) + " Hari"
	} else if durasiDalamJam > 0 {
		durasi = "durasi: " + strconv.FormatFloat(durasiDalamJam, 'f', 0, 64) + " Jam"
	} else {
		durasi = "durasi: 0 Hari"
	}

	var newProject = Project{
		NamaProject: projectName,
		Durasi:      durasi,
		Description: description,
		NodeJS:      nodejs,
		ReactJS:     reactjs,
		JavaScript:  javascript,
		SocketIO:    socketio,
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/index.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	response := map[string]interface{}{
		"Projects": dataProject,
	}

	tmpl.Execute(w, response)
}

func addMyProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/add_myproject.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func contactMe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/contact_me.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func detailProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/detail_project.html")

	if err != nil {
		w.Write([]byte("message : " + err.Error()))
		return
	}

	var ProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				NamaProject: data.NamaProject,
				Description: data.Description,
				Durasi:      data.Durasi,
			}
		}
	}

	data := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpl.Execute(w, data)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	// fmt.Println(index)

	dataProject = append(dataProject[:index], dataProject[index+1:]...)
	// fmt.Println(dataBlog)

	http.Redirect(w, r, "/home", http.StatusFound)
}

func halamanEdit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var tmpl, err = template.ParseFiles("html/halaman_edit.html")

	if err != nil {
		w.Write([]byte("message :" + err.Error()))
		return
	}
	var ProjectDetail = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	for i, data := range dataProject {
		if index == i {
			ProjectDetail = Project{
				NamaProject: data.NamaProject,
				Description: data.Description,
				Durasi:      data.Durasi,
			}

		}
	}
	data := map[string]interface{}{
		"EditProject": ProjectDetail,
	}
	tmpl.Execute(w, data)
}

func submitHalamanEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	indexedit, _ := strconv.Atoi(mux.Vars(r)["indexedit"])

	var projectName = r.PostForm.Get("project-name")
	var description = r.PostForm.Get("description")
	var startdate = r.PostForm.Get("start-date")
	var enddate = r.PostForm.Get("end-date")
	var nodejs = r.PostForm.Get("tech")
	var reactjs = r.PostForm.Get("tech2")
	var javascript = r.PostForm.Get("tech3")
	var socketio = r.PostForm.Get("tech4")

	Format := "2006-01-02"
	var sdate, _ = time.Parse(Format, startdate)
	var edate, _ = time.Parse(Format, enddate)
	durasiDalamJam := edate.Sub(sdate).Hours()

	durasiDalamHari := durasiDalamJam / 24
	durasiDalamBulan := durasiDalamHari / 30
	durasiDalamTahun := durasiDalamBulan / 12

	var durasi string
	var hari, _ float64 = math.Modf(durasiDalamHari)
	var bulan, _ float64 = math.Modf(durasiDalamBulan)
	var tahun, _ float64 = math.Modf(durasiDalamTahun)

	if tahun > 0 {
		durasi = "durasi: " + strconv.FormatFloat(tahun, 'f', 0, 64) + " Tahun"
	} else if bulan > 0 {
		durasi = "durasi: " + strconv.FormatFloat(bulan, 'f', 0, 64) + " Bulan"
	} else if hari > 0 {
		durasi = "durasi: " + strconv.FormatFloat(hari, 'f', 0, 64) + " Hari"
	} else if durasiDalamJam > 0 {
		durasi = "durasi: " + strconv.FormatFloat(durasiDalamJam, 'f', 0, 64) + " Jam"
	} else {
		durasi = "durasi: 0 Hari"
	}

	var newProject = Project{
		NamaProject: projectName,
		Durasi:      durasi,
		Description: description,
		NodeJS:      nodejs,
		ReactJS:     reactjs,
		JavaScript:  javascript,
		SocketIO:    socketio,
	}

	dataProject[indexedit] = newProject

	http.Redirect(w, r, "/home", http.StatusMovedPermanently)

}
