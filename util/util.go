package util

import (
	"bufio"
	"strconv"
	"os"
	"morego/golog"
)

func saveFile(str string, n int) {
	f, err := os.Create("./output" + strconv.Itoa(n) + ".txt") //创建文件

	if err != nil {
		golog.Error("os.Create Error:", err.Error())
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f) //创建新的 Writer 对象
	_, errw := w.WriteString(str)
	if errw != nil {
		golog.Error("WriteString Error:", errw.Error())
		return
	}
	//fmt.Printf("写入 %d 个字节\n", n4)
	w.Flush()
	f.Close()

}


