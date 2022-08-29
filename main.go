package main

import (
	"bufio"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"os"
	"strings"
)

// 文件路径切片
var fileList []string

var w fyne.Window

var content map[string]int

func main() {
	a := app.New()
	w = a.NewWindow("设备号文件合并并去重")

	showMessage := widget.NewLabel(showInfo())
	openFd := widget.NewButton("1.open file", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			fileList = append(fileList, reader.URI().Path())
			showMessage.SetText(showInfo())
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})
	mergeButton := widget.NewButton("2.merge File", mergeFiles)
	DownloadButton := widget.NewButton("3.Download file", DownloadFiles)

	w.SetContent(container.NewVBox(
		showMessage,
		widget.NewButton("clear files", func() {
			fileList = []string{}
			showMessage.SetText(showInfo())
		}),
		openFd,
		mergeButton,
		DownloadButton,
	))

	w.Resize(fyne.NewSize(500, 500))

	w.ShowAndRun()
}

func showInfo() string {
	if len(fileList) == 0 {
		return "Please select the files you need to merge:"
	}

	if len(fileList) > 0 {
		t := "The files you have selected are as follows:\n"
		t += strings.Join(fileList, "\n")
		return t
	}

	return ""
}

func mergeFiles() {
	content = make(map[string]int)
	var num int

	for _, filePath := range fileList {
		f, err := os.Open(filePath)
		if err != nil {
			dialog.ShowError(errors.New("open file["+filePath+"] err \n "+err.Error()), w)
			return
		}
		//将文件加载到内存中
		lineReader := bufio.NewReader(f)

		// 循环读取文件内容
		for {
			line, _, err := lineReader.ReadLine()
			if err == io.EOF {
				break
			}
			// 去除两边的空格
			idf := strings.TrimSpace(string(line))
			// 判断是否有重复
			if _, ok := content[idf]; ok {
				num++
				continue
			}
			// 存放到map之中
			content[idf] = 1
		}

	}
	//fmt.Println(num, content)
	msg := fmt.Sprintf("Effective number:%d \n Number of repetitions:%d", len(content), num)
	dialog.ShowInformation("result", msg, w)
}

func DownloadFiles() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if writer == nil {
			log.Println("Cancelled")
			return
		}

		fileSaved(writer, w)
	}, w)
}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window) {
	defer func(f fyne.URIWriteCloser) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	for s, _ := range content {
		_, err := f.Write([]byte(s + "\n"))
		if err != nil {
			dialog.ShowError(err, w)
		}
	}
	err := f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	log.Println("Saved to...", f.URI())
}
