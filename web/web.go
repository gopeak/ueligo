package web

import (
	"fmt"
	"morego/global"
	"net/http"
	"os"

)



/**
 * 建立web server
 */
func HttpServer() {

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/web/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))

	go func() {
		http.ListenAndServe(":"+global.Config.Admin.HttpPort, nil)
	}()
}

