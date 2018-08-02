package comm

import (
	"bytes"
	"eInfusion/logs"
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// GetTimeStamp :获取时间戳
func GetTimeStamp() string {
	return string(time.Now().Unix())
}

// GetPureIPAddr : 获取IP中纯的地址，去除字符串中的端口数据
func GetPureIPAddr(ip string) string {
	return strings.Split(ip, ":")[0]
}

// GetCurrentTime : 获取当前时间
func GetCurrentTime() string {
	strTime := time.Now().Format("2006-01-02 15:04:05") //后面的参数是固定的 否则将无法正常输出
	return strTime
}

// GetCurrentDate :获取当前日期
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// SepLi :生成分隔行
func SepLi(rN int, rChar string) {
	s := "-"
	if rChar != "" {
		s = rChar
	}
	Msg(strings.Repeat(s, rN))
}

// TrimSpc :去除左右空格
func TrimSpc(rStr string) string {
	return strings.TrimSpace(rStr)
}

// Msg :打印到屏幕
func Msg(v ...interface{}) {
	fmt.Println(v...)
}

// GetVarType :获取变量类型
func GetVarType(rVal interface{}) string {
	return fmt.Sprint(reflect.TypeOf(rVal))
}

// CkErr :处理错误(如果有错误，返回true,无错则返回false),同时记录日志
func CkErr(rMsgTitle string, rErr error) bool {
	if rErr != nil {
		logs.LogMain.Error(rMsgTitle, rErr)
		return true
	}
	return false
}

// ConvertBasNumberToStr :把数值类型数据（仅支持int/uint）转换成指定进制数值，返回字符串
func ConvertBasNumberToStr(rBase int, rVal interface{}) string {
	//	reflect.TypeOf(ref_varContent)
	var strBaseValue string
	switch rBase {
	case 16:
		strBaseValue = "X"
	case 10:
		strBaseValue = "d"
	case 2:
		strBaseValue = "b"
	case 1:
		strBaseValue = "v"
	default:
		strBaseValue = "x"
	}
	strRetValue := fmt.Sprintf("%"+strBaseValue, rVal)
	return strRetValue
}

// ConvertOxBytesToStr :转换Bytes为十六进制字符串 (此方法没有用HEX包)
//转换过程中可能会丢失0，因此需要补0
func ConvertOxBytesToStr(rCnt []byte) string {
	var strRet string
	for i := 0; i < len(rCnt); i++ {
		strCon := ConvertBasNumberToStr(16, rCnt[i])
		if len(strCon) == 1 {
			strCon = "0" + strCon
		}
		strRet += strCon
	}
	return strRet
}

// ConvertBasStrToInt :把指定进制的字符转换成为十进制数值（int型）
// 请注意，只能还原数值的进制，并不能转换进制
func ConvertBasStrToInt(rBase int, rStrCnt string) int {
	intRetValue, err := strconv.ParseInt(rStrCnt, rBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return int(intRetValue)
}

// ConvertBasStrToUint :根据指定进制要求，把字符串转换成数字Uint8型
//  请注意，只能还原数值的进制，并不能转换进制
func ConvertBasStrToUint(rBase int, rStrCnt string) uint8 {
	intRetValue, err := strconv.ParseUint(rStrCnt, rBase, 64)
	if err != nil {
		intRetValue = 0
	}
	return uint8(intRetValue)
}

// ConvertStrIPToHexBytes :把指定十进制IP地址转换成为十六进制bytes
func ConvertStrIPToHexBytes(rIP string) []byte {
	st := strings.SplitN(rIP, ".", 4)
	var bs []byte
	for i := 0; i < len(st); i++ {
		mv, _ := strconv.ParseUint(st[i], 10, 64)
		t := ConvertBasStrToUint(16, ConvertBasNumberToStr(16, mv))
		bs = append(bs, byte(t))
	}
	return bs
}

// ConvertEvenDecToHexBytes :偶数十进制数值转换为十六进制bytes
// 请注意：只支持偶数位
func ConvertEvenDecToHexBytes(rStrCnt string) []byte {
	var bs []byte
	ms := ConvertBasNumberToStr(16, ConvertBasStrToInt(10, rStrCnt))
	if len(ms) > 2 {
		t := ConvertStrToBytesByPerTwoChar(ms)
		for i := 0; i < len(t); i++ {
			bs = append(bs, t[i])
		}
	}
	return bs
}

// ConvertDecToHexBytes :十进制数值转换为十六进制bytes
func ConvertDecToHexBytes(rIntCnt int) []byte {
	// FIXME: 转换成为十六进制还有问题
	var rbs []byte
	mbs := ConvertIntToBytes(rIntCnt)
	for i := 0; i < len(mbs); i++ {
		if mbs[i] != 0 {
			ts := ConvertBasNumberToStr(10, mbs[i])

			rbs = append(rbs, ConvertBasStrToUint(16, ts))
		}
	}
	return rbs
}

// ConvertBasStrToBytes :根据开始、结束下标返回相应的字符串内容返回bytes
func ConvertBasStrToBytes(rStrCnt string, rBegin int, rEnd int, rBase int) []byte {
	var bT []byte
	n := len(rStrCnt)
	if rEnd <= n && rBegin >= 0 {
		for i := rBegin; i <= rEnd; i++ {
			strT := string(rStrCnt[i])
			bT = append(bT, ConvertBasStrToUint(rBase, strT))
		}
	}
	return bT
}

//GetPartOfStr :根据开始、结束下标返回相应的字符串内容返回string
func GetPartOfStr(rStrCnt string, rBegin int, rEnd int) string {
	var strR string
	n := len(rStrCnt)
	if rEnd < n && rBegin >= 0 {
		for i := rBegin; i <= rEnd; i++ {
			strR += string(rStrCnt[i])
		}
	}
	return strR
}

// ConvertByteToBinaryOfBytes :转换byte内的数据为二进制的byte切片
func ConvertByteToBinaryOfBytes(rByte byte) []byte {
	var bT []byte
	s := ConvertBasNumberToStr(2, rByte)
	for i := 0; i < len(s); i++ {
		t, _ := strconv.ParseUint(string(s[i]), 10, 64)
		tt := uint8(t)
		bT = append(bT, tt)
	}
	return bT
}

// ConvertStrToBytesByPerTwoChar :把字符串内容按每两字符对应一个byte组成新的bytes，返回[]byte
// 注意：目前只支持偶数位字符转换
func ConvertStrToBytesByPerTwoChar(rStrCnt string) []byte {
	var bT []byte
	n := len(rStrCnt)
	i := 0
	for i < n-1 {
		j := i + 1
		strP := string(rStrCnt[i]) + string(rStrCnt[j])
		bT = append(bT, ConvertBasStrToUint(16, string(strP)))
		i = j + 1
	}
	return bT
}

// IsExists :判断是否存在,true表示存在
func IsExists(rPath string) bool {
	_, err := os.Stat(rPath)
	if err != nil {
		return false
	}
	return true
}

// GetCurrentDirectory :获取当然路径
func GetCurrentDirectory() string {
	var strPath string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		strPath = "\\"
	} else {
		strPath = "/"
	}
	dir, _ := os.Getwd()
	strPath = dir + strPath
	return strPath
}

// GetPathSeparator :根据操作系统环境，返回分隔符
func GetPathSeparator() string {
	var strPath string
	if os.IsPathSeparator('\\') { //前边的判断是否是系统的分隔符
		strPath = "\\"
	} else {
		strPath = "/"
	}
	return strPath
}

// ConvertIntToStr :转化整形成字符型
func ConvertIntToStr(intContent int) string {
	return strconv.Itoa(intContent)
}

// ConvertIntToBytes :整形转换成bytes
func ConvertIntToBytes(n int) []byte {
	tmp := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, tmp)
	// binary.Write(bytesBuffer, binary.LittleEndian, tmp)
	return bytesBuffer.Bytes()
}

// WrToFilWithBuffer :写入指定文件，如果没有该文件自动生成
func WrToFilWithBuffer(rFilePath string, rStrCnt string, rIsAppend bool) bool {
	//	this function for complex content to file
	var intFileOpenMode int
	if rIsAppend {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	} else {
		intFileOpenMode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	fileHandle, err := os.OpenFile(rFilePath, intFileOpenMode, 0666)
	if err != nil {
		return false
	} else {
		// write to the file
		fileHandle.WriteString("\r\n" + rStrCnt)
	}
	defer fileHandle.Close()
	return true
}
