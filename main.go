package main

import (
	"bytes"

	"fmt"

	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode"
)

type Projector struct {
	QRCode1     string `json:"qrcode1,omitempty"`
	QRCode2     string `json:"qrcode2,omitempty"`
	Geolocation string `json:"geolocation,omitempty"`
	Video       string `json:"video,omitempty"`
}

type Projectornew struct {
	QRCode1     string `json:"qrcode1,omitempty"`
	QRCode2     string `json:"qrcode2,omitempty"`
	Geolocation string `json:"geolocation,omitempty"`
	Video       []byte `json:"video,omitempty"`
}

func ConsumeMobileDataNew(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method)
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("hello")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("hello1")
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	QRCode1 := r.FormValue("qrcode1")
	QRCode2 := r.FormValue("qrcode2")

	fmt.Println(QRCode1)

	err1 := qrcode.WriteFile(QRCode1, qrcode.Medium, 256, "1.jpg")
	err2 := qrcode.WriteFile(QRCode2, qrcode.Medium, 256, "2.jpg")

	if err1 != nil || err2 != nil {
		fmt.Println(err)
	}

	convertImagetoVideos()
}

func convertImagetoVideos() {
	cmdArguments := []string{"-loop", "1", "-i", "1.jpg", "-loop", "1", "-c:v", "libx264", "-t", "10", "-pix_fmt", "yuv420p", "-vf", "scale=320:240", "-y", "qr1.mp4"}
	cmdArguments1 := []string{"-loop", "1", "-i", "2.jpg", "-loop", "1", "-c:v", "libx264", "-t", "10", "-pix_fmt", "yuv420p", "-vf", "scale=320:240", "-y", "qr2.mp4"}
	cmdArguments2 := []string{"-i", "qr1.mp4", "-i", "qr2.mp4", "-filter_complex", "concat=n=2:v=1:a=0", "-f", "MOV", "-an", "-y", "qr.mp4"}

	cmdArguments3 := []string{"-i", "qr.mp4", "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", "-y", "temp1.ts"}
	cmdArguments4 := []string{"-i", "small.mp4", "-c", "copy", "-bsf:v", "h264_mp4toannexb", "-f", "mpegts", "-y", "temp2.ts"}
	cmdArguments5 := []string{"-i", "concat:temp2.ts|temp1.ts", "-c", "copy", "-bsf:a", "aac_adtstoasc", "-y", "output.mp4"}

	cmd := exec.Command("ffmpeg", cmdArguments...)
	cmd1 := exec.Command("ffmpeg", cmdArguments1...)
	cmd2 := exec.Command("ffmpeg", cmdArguments2...)
	cmd3 := exec.Command("ffmpeg", cmdArguments3...)
	cmd4 := exec.Command("ffmpeg", cmdArguments4...)
	cmd5 := exec.Command("ffmpeg", cmdArguments5...)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	err1 := cmd1.Run()
	err2 := cmd2.Run()
	err3 := cmd3.Run()
	err4 := cmd4.Run()
	err5 := cmd5.Run()

	if err != nil || err1 != nil || err2 != nil || err4 != nil || err5 != nil {
		log.Fatal(err, err1, err2, err3, err4, err5)
	}
	fmt.Printf("command output: %q\n", out.String())
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/projectornew/", ConsumeMobileDataNew).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
