package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"log"
	"net/http"
	"os/exec"
)



func makeLog(errlog error){

	logfile,err:=os.OpenFile("/var/log/AsteriskAgent.log",os.O_RDWR|os.O_CREATE| os.O_APPEND, 0755)
	if err != nil{
		log.Fatal(err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)
	log.Println(errlog.Error())

}


func middleware(next http.HandlerFunc)http.HandlerFunc{

	return func(w http.ResponseWriter, r *http.Request) {
		errorAccess:=errors.New("Access Denied")
		a := r.URL.Query().Get("token")
		if a != "b52c96bea30646abf8170f333bbd42b9" {
			makeLog(errorAccess)
			fmt.Fprintln(w,"Access Denied")
			return
		}

		next(w,r)
	}


}


func handler( w http.ResponseWriter, r *http.Request){


	data:=make(map[string]map[string]string)
	Status:=make(map[string]string)
	hostname,_:=os.Hostname()
	cmd := "sudo asterisk -rx 'sip show peers'"
	res, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		makeLog(err)
	}
	cmd1 := "sudo asterisk -rx 'sip show registry'"
	res1, err := exec.Command("bash", "-c", cmd1).Output()
	if err != nil {
		makeLog(err)
	}

	cmd2 := "sudo asterisk -rx 'ooh323 show peers'"
	res2, err := exec.Command("bash", "-c", cmd2).Output()
	if err != nil {
		makeLog(err)
	}

	Status["RegisterStatus"]=string(res1)
	Status["PeerStatusSIP"]=string(res)
	Status["PeerStatusH323"]=string(res2)
	data[hostname]=Status
	w.Header().Set("Content-type","application/json")
	json.NewEncoder(w).Encode(data)
}


func main(){

	http.HandleFunc("/",middleware(handler))
	makeLog(http.ListenAndServe(":8082",nil))



}
