package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	is_video_ext := map[string]bool{
		".flv":  true,
		".mp4":  true,
		".avi":  true,
		".mov":  true,
		".webp": true,
	}

	dir := "D:\\EP\\Git\\test_folder"
	files, _ := ioutil.ReadDir(dir)
	processedFiles := map[string]bool{}
	for _, f := range files {
		fmt.Println(f.Name())
		processedFiles[f.Name()] = true
	}
	for {
		time.Sleep(3 * time.Second)
		fmt.Println("Checking...")
		files, _ := ioutil.ReadDir(dir)
		for _, f := range files {
			if processedFiles[f.Name()] {
				continue
			}
			if is_video_ext[filepath.Ext(f.Name())] {
				fmt.Println("Found new video file: " + f.Name())
				outfile := f.Name() + ".mp4" //filepath.Base(f.Name()) + ".mp4"
				cmd := exec.Command("ffmpeg", "-i", f.Name(), outfile)
				cmd.Dir = dir
				if err := cmd.Run(); err != nil {
					fmt.Println("Error converting file:", err)
				} else {
					fmt.Println("Conversion successful. Output file: " + outfile)

					processedFiles[outfile] = true
					processedFiles[f.Name()] = true

					e := os.Remove(filepath.Join(dir, f.Name()))
					if e != nil {
						fmt.Println("Error deleting file:", e)
					}
				}

			}
		}
	}
}
