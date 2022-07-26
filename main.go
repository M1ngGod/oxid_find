package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var ip_file = flag.String("f","","input ip list file")
var ip_segment = flag.String("h",""," Input IP segment (eg.192.168.1.1/24，192.168.1.1-192.168.2.18)")
var ip = flag.String("ip",""," Input target host IP")



func main() {

	flag.Parse()
	//fmt.Println(*ip_file)
	//fmt.Println(*ip_segment)
	//runtime.GOMAXPROCS(2)

	//执行 -ip命令
	if *ip == "" && *ip_file == "" && *ip_segment == "" {
		fmt.Println("not input Param")
		os.Exit(0)
	}else if *ip != "" && *ip_file == "" && *ip_segment == "" {
		scan(string(*ip))

	//执行 -f 命令
	}else if *ip == "" && *ip_file != "" && *ip_segment == ""{
		buff := readAllBuff(*ip_file)
		split := strings.Split(buff, "\n")
		for i := range split {
			scan(split[i])
		}
	//执行 -h 命令
	}else if *ip == "" && *ip_file == "" && *ip_segment != ""{
		fmt.Println(string(*ip_segment))
		host := string(*ip_segment)
		if strings.Contains(host,"/24") {
			split := strings.Split(host, ".")
			for i := 1;i <= 255;i++{
				s := strconv.Itoa(i)
				target := split[0]+"."+split[1]+"."+split[2]+"."+string(s)
				scan(target)
			}
		}else if strings.Contains(host, "/16") {
			split := strings.Split(host, ".")
			for i := 1;i < 255;i++{
				s := strconv.Itoa(i)
				for b := 1;b <= 255;b++{
					f := strconv.Itoa(b)
					target := split[0]+"."+split[1]+"."+string(s)+"."+string(f)
					scan(target)
				}
			}
		//192.168.1.1-192.168.2.45
		}else if strings.Contains(host, "-"){
			split := strings.Split(host, ".")
			fmt.Println(split)
			num, err := strconv.Atoi(split[2])
			if err != nil {
				fmt.Println("ip error")
			}

			num1, err1 := strconv.Atoi(split[5])
			if err != nil && err1 != nil {
				fmt.Println("ip error")
			}

			i1 := strings.Split(host, "-")
			fmt.Println(i1)
			for i := num; i<= num1;i++{
				f := strconv.Itoa(i)
				for b := 1;b <= 255;b++{
					e := strconv.Itoa(b)
					target := split[0]+"."+split[1]+"."+string(f)+"."+string(e)
					scan(target)
				}
			}
		}

	}

}


func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func delZero(arr []byte) []byte {
	var buffer []byte
	for i := range arr {
		if arr[i] == 0{
			continue
		}else{
			//fmt.Println(arr[i])
			buffer = append(buffer, arr[i])
		}
	}
	return buffer
}

func IndexStr(arr []byte) int {
	var Num int
	for i := range arr {
		if arr[i] == 9 && arr[i+1] == 0 && arr[i+2] == 255 && arr[i+3] == 255  && arr[i+4] == 0 && arr[i+5] == 0{
			Num = i
		}
	}
	return Num

}

func scan(ip string){
	connTimeout := 1*time.Second
	conn, err := net.DialTimeout("tcp",ip+":135",connTimeout)

	//err == nil 时表示没有错误
	if err != nil {

		fmt.Println("timeout", err)
		return
	}


	//发送的数据是字符切片的格式
	buffer_v1 := []byte("\x05\x00\x0b\x03\x10\x00\x00\x00\x48\x00\x00\x00\x01\x00\x00\x00\xb8\x10\xb8\x10\x00\x00\x00\x00\x01\x00\x00\x00\x00\x00\x01\x00\xc4\xfe\xfc\x99\x60\x52\x1b\x10\xbb\xcb\x00\xaa\x00\x21\x34\x7a\x00\x00\x00\x00\x04\x5d\x88\x8a\xeb\x1c\xc9\x11\x9f\xe8\x08\x00\x2b\x10\x48\x60\x02\x00\x00\x00")
	buffer_v2 := []byte("\x05\x00\x00\x03\x10\x00\x00\x00\x18\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x05\x00")
	//向服务器发送数据
	cnt, err := conn.Write(buffer_v1)
	//接收服务器返回的数据
	buf1 := make([]byte, 1024)
	cnt, err = conn.Read(buf1)
	if err != nil {
		fmt.Println("conn.Write err:", err,cnt)
		return
	}

	cnt2, err2 := conn.Write(buffer_v2)
	//接收服务器返回的数据
	buf2 := make([]byte, 4096)
	cnt2, err2 = conn.Read(buf2)

	if err2 != nil {
		fmt.Println("conn.Write err:", err2,cnt2)
		return
	}

	subByte := BytesToString(buf2[42:])
	byteBuff := StringToBytes(subByte)

	indexStr := IndexStr(byteBuff)
	bytes := byteBuff[:indexStr]
	zero := delZero(bytes)
	toString := BytesToString(zero)

	//fmt.Println(toString)
	fmt.Println("[*] ",ip,"\t")
	split := strings.Split(toString, "\a")
	for i := range split {
		fmt.Println("\t[->]",split[i])
	}
	runtime.Gosched()
	defer conn.Close()
}



func readAllBuff(filePath string) string{
//	start1 := time.Now()
	// 打开文件
	FileHandle, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	// 关闭文件
	defer FileHandle.Close()
	// 获取文件当前信息
	fileInfo, err := FileHandle.Stat()
	if err != nil {
		log.Println(err)
		os.Exit(0)
	}
	buffer := make([]byte, fileInfo.Size())
	// 读取文件内容,并写入buffer中
	n, err := FileHandle.Read(buffer)
	if err != nil {
		log.Println(err)
	}

	return string(buffer[:n])
//	fmt.Println("readAllBuff spend : ", time.Now().Sub(start1))
}
