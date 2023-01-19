package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//	tele "gopkg.in/telebot.v3"

var dir string

var chats_fname string = "chats.txt"
var dir_fname string = "dir.txt"
var token_fname string = "token.txt"

var is_video_ext = map[string]bool{
	".flv": true,
	".mp4": true,
	".avi": true,
	".mov": true,
}

var is_photo_ext = map[string]bool{
	".bmp":  true,
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".tiff": true,
	".webp": true,
}

var Token string
var Chats map[int]bool = nil

var processedFiles = map[string]bool{}

func main() {
	load_data()

	go regular_func(update_n_save_chat_ids, 500)

	init_scan_folder()

	go regular_func(rescan_folder, 1)

	for {
		time.Sleep(1)
	}
}

func load_data() {
	load_token()

	load_chats()

	load_dir()
}

func botAPIUrl(method_n_stuff string) string {
	return "https://api.telegram.org/bot" + Token + "/" + method_n_stuff
}

func send_text(text string, chat_id int) {
	get_request(botAPIUrl("sendMessage?chat_id=" + strconv.Itoa(chat_id) + "&text=" + url.QueryEscape(text)))
}

func send_photo_to_all(filePath string) {
	for chat_id := range Chats {
		send_photo(filePath, chat_id)
	}
}

func send_video_to_all(filePath string) {
	for chat_id := range Chats {
		send_video(filePath, chat_id)
	}
}

func send_photo(filePath string, chat_id int) {
	send_multipart_file(filePath, botAPIUrl("sendPhoto?chat_id="+strconv.Itoa(chat_id)), "photo")
}

func send_video(filePath string, chat_id int) {
	send_multipart_file(filePath, botAPIUrl("sendVideo?chat_id="+strconv.Itoa(chat_id)), "video")
}

// https://gist.github.com/andrewmilson/19185aab2347f6ad29f5
func send_multipart_file(filePath string, requestURL string, file_variable string) {
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(file_variable, filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", requestURL, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	client.Do(r)
}

func regular_func(f func(), delay float64) {
	for {
		go f()
		time.Sleep(time.Second * time.Duration(delay))
	}
}

func update_n_save_chat_ids() {
	if update_chat_ids() {
		save_chats()
	} else {
		fmt.Println("No new chat_ids")
	}
}

func get_request(requestURL string) []byte {
	fmt.Println(requestURL)
	// requestURL := fmt.Sprintf("http://localhost:%d", serverPort)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		return nil
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("client: got response!\n")
	fmt.Printf("client: status code: %d\n", res.StatusCode)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("client: response body: %s\n", resBody)

	return resBody
}

func parse_request_body(resBody []byte) map[string]interface{} {
	var data map[string]interface{}

	err := json.Unmarshal([]byte(resBody), &data)
	fmt.Println(err)

	fmt.Printf("%#v\n", data)

	return data
}

func update_chat_ids() bool {
	// body := string(get_request("https://api.telegram.org/bot" + Token + "/getUpdates?offset=0"))

	data := parse_request_body(get_request(botAPIUrl("getUpdates?offset=0")))

	fmt.Println("*************************")
	fmt.Println("*************************")
	fmt.Println("*************************")

	fmt.Printf("ok = %#v\n", data["ok"])

	need2save := false

	for _, tres := range data["result"].([]interface{}) {
		chat_id := int(tres.(map[string]interface{})["message"].(map[string]interface{})["chat"].(map[string]interface{})["id"].(float64))
		fmt.Printf("%#v\n", chat_id)
		if !Chats[chat_id] {
			Chats[chat_id] = true
			need2save = true
		}
	}

	return need2save
}

func load_token() {
	b, err := os.ReadFile(token_fname) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	Token = string(b) // convert content to a 'string'

	fmt.Println(Token) // print the content as a 'string'
}

func load_dir() {
	b, err := os.ReadFile(dir_fname) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	dir = string(b) // convert content to a 'string'

	fmt.Println(dir) // print the content as a 'string'
}

func load_chats() {
	Chats = make(map[int]bool)
	b, err := os.ReadFile(chats_fname) // just pass the file name
	if err != nil {
		fmt.Print(err)
		return
	}

	str := string(b) // convert content to a 'string'
	strs := strings.Split(strings.ReplaceAll(str, "\r\n", "\n"), "\n")

	for _, s := range strs {
		if len(s) == 0 {
			continue
		}
		tint, err := strconv.Atoi(s)
		if err != nil {
			fmt.Println("Error reading number from s=", s, err)
			continue
		}
		Chats[tint] = true
	}

	for chat_id := range Chats {
		fmt.Println(chat_id)
	}
}

func save_chats() {
	fmt.Println("Saving chats..")
	f, err := os.Create(chats_fname)
	if err != nil {
		fmt.Print(err)
		return
	}

	for chat_id := range Chats {
		_, err := f.WriteString(strconv.Itoa(chat_id) + "\r\n")
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Chats saved")
}

func init_scan_folder() {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		fmt.Println(f.Name())
		processedFiles[f.Name()] = true
	}
}

func rescan_folder() {
	fmt.Println("Checking...")
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if processedFiles[f.Name()] {
			continue
		}
		if is_video_ext[filepath.Ext(f.Name())] {
			processedFiles[f.Name()] = true

			fmt.Println("Found new video file: " + f.Name())
			outfile := f.Name() + ".mp4" //filepath.Base(f.Name()) + ".mp4"

			processedFiles[outfile] = true
			cmd := exec.Command("ffmpeg", "-i", f.Name(), outfile)
			cmd.Dir = dir
			if err := cmd.Run(); err != nil {
				fmt.Println("Error converting file:", err)
			} else {
				fmt.Println("Conversion successful. Output file: " + outfile)

				e := os.Remove(filepath.Join(dir, f.Name()))
				if e != nil {
					fmt.Println("Error deleting file:", e)
				}

				send_video_to_all(filepath.Join(dir, outfile))
			}

		}
		if is_photo_ext[filepath.Ext(f.Name())] {
			fmt.Println("Found new photo: ", f.Name())
			send_photo_to_all(filepath.Join(dir, f.Name()))
			processedFiles[f.Name()] = true
		}
	}
}
