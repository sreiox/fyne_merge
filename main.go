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
	"main/zh"
	"os"
	"strings"
)

// 文件路径切片
var fileList []string

var w fyne.Window

var content map[string]int

func main() {
	a := app.New()
	a.Settings().SetTheme(&zh.MyTheme{})
	w = a.NewWindow("设备号文件合并并去重")

	showMessage := widget.NewLabel(showInfo())
	openFd := widget.NewButton("1.选择文件", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("取消了...")
				return
			}
			fileList = append(fileList, reader.URI().Path())
			showMessage.SetText(showInfo())
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))
		fd.Show()
	})
	mergeButton := widget.NewButton("2.合并文件", mergeFiles)
	DownloadButton := widget.NewButton("3.下载文件", DownloadFiles)

	w.SetContent(container.NewVBox(
		showMessage,
		widget.NewButton("清空文件", func() {
			clearFile()
			clearContent()
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
		return "请选择你需要合并的文件:"
	}

	if len(fileList) > 0 {
		t := "所选择的文件如下:\n"
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
			dialog.ShowError(errors.New("打开文件["+filePath+"]失败: \n "+err.Error()), w)
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
	msg := fmt.Sprintf("有效数量:%d \n 重复数量:%d", len(content), num)
	dialog.ShowInformation("合并结果", msg, w)
}

func DownloadFiles() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if writer == nil {
			log.Println("取消")
			return
		}

		fileSaved(writer, w)
	}, w)
	clearContent()
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
	log.Println("保存路径为：", f.URI())
	dialog.ShowInformation("保存成功", "文件保存路径："+f.URI().String(), w)
}

func clearFile() {
	fileList = []string{} // 清空选择文件
}

func clearContent() {
	content = make(map[string]int) // 清空内容
}
