package util

import (
	"bufio"
	"strconv"
	"os"
	"morego/golog"
	"strings"
	"crypto/rand"
	"math/big"
	"math"
	"encoding/binary"
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

//  转义json字符串
func EncodeJsonStr(str string) string {
	str = strings.Replace(str, `"`, `\"`, -1)
	return str
}

// 反解json字符串
func DecodeJsonStr(str string) string {
	str = strings.Replace(str, `\"`, `"`, -1)
	return str
}

func RandInt64(min,max int64) int64{
	maxBigInt:=big.NewInt(max)
	i,_:=rand.Int(rand.Reader,maxBigInt)
	if i.Int64()<min{
		RandInt64(min,max)
	}
	return i.Int64()
}

func Float32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)

	return bytes
}

func ByteToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)

	return math.Float32frombits(bits)
}

func Float64ToByte(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)

	return bytes
}

func ByteToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)

	return math.Float64frombits(bits)
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}