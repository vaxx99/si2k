package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

//Flags struct
type Flags struct {
	f01 int64
	f02 int64
	f03 int64
	f04 int64
	f05 int64
	f06 int64
	f07 int64
	f08 int64
	f09 int64
	f10 int64
	f11 int64
	f12 int64
	f13 int64
	f14 int64
	f15 int64
	f16 int64
	f17 int64
	f18 int64
	f19 int64
}

//Frec 200 struct
type Frec struct {
	IDR  int64
	IDC  int64
	FLG  Flags
	SLS  int64
	CHS  int64
	ZCL  int64
	SLL  int64
	ZCD  string
	SNC  string
	P100 P100
	P102 P102
	P103 P103
	P104 P104
	P105 P105
	P106 P106
	P107 P107
	P108 P108
	P109 P109
	P110 P110
	P111 P111
	P112 P112
	P113 P113
	P114 P114
	P115 P115
	P116 P116
	P119 P119
	P121 P121
}

//P100 struct
type P100 struct {
	IDI int64
	CNL int64
	CNC string
}

//P102 struct
type P102 struct {
	IDI int64
	DTS string
	F1  int64
}

//P103 struct
type P103 struct {
	IDI int64
	DTE string
	F1  int64
}

//P104 struct
type P104 struct {
	IDI int64
	CNT int64
}

//P105 struct
type P105 struct {
	IDI int64
	SVC int64
	TVC int64
}

//P106 struct
type P106 struct {
	IDI int64
	SVC int64
}

//P107 struct
type P107 struct {
	IDI int64
	SVC int64
}

//P108 struct
type P108 struct {
	IDI int64
	TPE int64
	SVC int64
}

//P109 struct
type P109 struct {
	IDI int64
	CNL int64
	CNC string
}

//P110 struct
type P110 struct {
	IDI int64
	CAT int64
}

//P111 struct
type P111 struct {
	IDI int64
	DIR int64
}

//P112 struct
type P112 struct {
	IDI int64
	CFC int64
}

//P113 struct
type P113 struct {
	IDI int64
	TGN int64
	SLN int64
	MDN int64
	PTN int64
	CHN int64
}

//P114 struct
type P114 struct {
	IDI int64
	TGN int64
	SLN int64
	MDN int64
	PTN int64
	CHN int64
}

//P115 struct
type P115 struct {
	IDI int64
	DUR int64
}

//P116 struct
type P116 struct {
	IDI int64
	BTL int64
	CRC int64
}

//P119 struct
type P119 struct {
	IDI int64
	BTL int64
	CNL int64
	CNC string
}

//P121 struct
type P121 struct {
	IDI int64
	BTL int64
	COI int64
	CNC string
}

var wd, sp string

func main() {
	switch len(os.Args) {
	case 2:
		wd = os.Args[1]
	case 3:
		wd = os.Args[1]
		sp = os.Args[2]
	default:
		wd = "."
		sp = ""
	}
	os.Chdir(wd)

	if info, err := os.Stat(wd); err == nil && info.IsDir() {
		f, _ := ioutil.ReadDir(wd)
		for _, fn := range f {
			if !fn.IsDir() && issi(fn.Name()) {
				cnt, rec := si2k(fn.Name())
				f, err := os.OpenFile(fn.Name()[1:], os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				defer f.Close()
				if err != nil {
					panic(err)
				}

				for _, sr := range rec {
					var ama string
					ama = Frec2Str(sr, sp)
					_, err = f.WriteString(ama + "\n")
				}
				fmt.Printf("%20s %8d %8d\n", fn.Name(), fn.Size(), cnt)
			}
		}
	}
}

func Frec2Str(sr Frec, sp string) string {
	var ama string
	idr := strconv.Itoa(int(sr.IDR))
	idc := strconv.Itoa(int(sr.IDC))
	cat := strconv.Itoa(int(sr.P110.CAT))
	dir := strconv.Itoa(int(sr.P111.DIR))
	itg := strconv.Itoa(int(sr.P113.TGN))
	otg := strconv.Itoa(int(sr.P114.TGN))
	dur := strconv.Itoa(int(sr.P115.DUR))
	ama = idr + sp + idc + sp + cat + sp + sr.ZCD + sp + sr.SNC + sp + sr.P100.CNC + sp + sr.P102.DTS + sp + sr.P103.DTE + sp + dir + sp + itg + sp + otg + sp + dur
	return ama
}

func issi(fn string) bool {
	f, _ := os.Open(fn)
	defer f.Close()
	data, _ := Read(f, 1)
	ad := H2c(data)
	if ad == "C8" {
		return true
	}
	return false
}

//Read bytes from file
func Read(file *os.File, bt int) ([]byte, error) {
	data := make([]byte, bt)
	_, e := file.Read(data)
	if e != nil {
		log.Println("File open error:", e)
	}
	return data, e
}

func b2i(b []byte) int64 {
	hd := strings.ToUpper(hex.EncodeToString(b))
	res, _ := strconv.ParseInt(hd, 16, 64)
	return res
}

func bc2i(b string) int64 {
	a, _ := strconv.ParseInt(b, 2, 64)
	return a
}

func check(e error) {
	if e != nil {
		log.Println("BCD error:", e)
		return
	}
}

//H2c hex to string
func H2c(dt []byte) string {
	hd := strings.ToUpper(hex.EncodeToString(dt))
	return hd
}

//Oct byte to string
func Oct(b byte) string {
	return fmt.Sprintf("%08b", b)
}

//Bts byte to int
func Bts(d int64) int {
	return int(float64(d)/2 + 0.5)
}

func flags(s string) Flags {
	var f []int64
	for _, j := range s {
		res, _ := strconv.ParseInt(string(j), 10, 64)
		f = append(f, res)
	}
	return Flags{f[7], f[6], f[5], f[4], f[3], f[2], f[1], f[7], f[15], f[14], f[13], f[12], f[11], f[10], f[9], f[8], f[18], f[17], f[16]}
}

func dates(b []byte) string {
	rd := ""
	if len(b) > 0 {
		rd = "20" + dd(int(b[0])) + dd(int(b[1])) + dd(int(b[2])) + dd(int(b[3])) +
			dd(int(b[4])) + dd(int(b[5])) + dd(int(b[6]))
	}
	return rd
}

func dd(d int) string {
	if d < 10 {
		return "0" + strconv.Itoa(d)
	}
	return strconv.Itoa(d)
}

func si2k(fn string) (int64, []Frec) {
	f, _ := os.Open(fn)
	var cnt int64
	var rec Frec
	var Rec []Frec
	defer f.Close()

	for {
		head := make([]byte, 3)
		_, e := f.Read(head)

		i := b2i(head[1:3])
		switch head[0] {
		case 200:
			b := make([]byte, i-3)
			_, e = f.Read(b)
			rec = s200(b, i-3)
			Rec = append(Rec, rec)
			cnt++

		case 210:
			b := make([]byte, 13)
			_, e = f.Read(b)

		case 211:
			b := make([]byte, 13)
			_, e = f.Read(b)

		case 212:
			b := make([]byte, 6)
			_, e = f.Read(b)
		}
		if e != nil {
			break
		}
	}
	return cnt, Rec
}

func s200(b []byte, bs int64) Frec {
	var srec Frec
	//Индекс записи
	srec.IDR = b2i(b[0:4])
	//Идентификатор вызова
	srec.IDC = b2i(b[4:8])
	//Flags
	fb := Oct(b[8]) + Oct(b[9]) + Oct(b[10])
	srec.FLG = flags(fb)
	bc := Oct(b[11])
	//Последовательность
	a, _ := strconv.ParseInt(bc[:4], 2, 8)
	//Состояние учета	стоимости
	c, _ := strconv.ParseInt(bc[4:], 2, 8)
	srec.SLS = a
	srec.CHS = c
	bc = Oct(b[12])
	//Длина кода зоны
	d, _ := strconv.ParseInt(bc[:3], 2, 8)
	//Длина списочного номера
	f, _ := strconv.ParseInt(bc[3:], 2, 8)
	srec.ZCL = d
	srec.SLL = f
	//Bytes count
	btz := Bts(d)
	btn := Bts(f)
	//Код зоны
	srec.ZCD = H2c(b[13 : 13+btz])[0:d]
	//Списочный номер абонента
	srec.SNC = H2c(b[13 : 13+btz+btn])[0 : d+f]
	//dynamic part start byte
	nb := 13 + int((float64(d)+float64(f))/2+0.5)
	for nb < int(bs) {
		id, _ := strconv.ParseInt(Oct(b[nb]), 2, 64)
		nb = dynp(id, nb, b, &srec)
	}

	return srec
}

//Переменная часть записи
func dynp(id int64, nb int, b []byte, rec *Frec) int {
	switch id {
	//Номер второго абонента (Called number) (100)
	case 100:
		var st P100
		//Идентификатор информационного элемента (100)
		st.IDI = id
		dc, _ := strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		//Длина списочного номера второго абонента
		st.CNL = dc
		btb := Bts(dc)
		//Списочный номер второго абонента
		st.CNC = H2c(b[nb+2 : nb+2+btb])[0:int(st.CNL)]
		nb = nb + 2 + btb
		rec.P100 = st
		return nb
		//Дата и время начала (Start date and time) (102)
	case 102:
		var st P102
		//Идентификатор информационного элемента (102)
		st.IDI = id
		bc := Oct(b[nb+9])
		//Флаг
		f1, _ := strconv.ParseInt(bc[7:], 2, 64)
		st.F1 = f1
		//Дата и время
		st.DTS = dates(b[nb+1 : nb+8])
		nb = nb + 9
		rec.P102 = st
		return nb
		//Дата и время завершения вызова (End date and time) (103)
	case 103:
		var st P103
		//Идентификатор информационного элемента (103)
		st.IDI = id
		bc := Oct(b[nb+9])
		f1, _ := strconv.ParseInt(bc[7:], 2, 64)
		//Флаг
		st.F1 = f1
		//Дата и время
		st.DTE = dates(b[nb+1 : nb+8])
		nb = nb + 9
		rec.P103 = st
		return nb
		//Количество тарифных импульсов (Number of charging units) (104)
	case 104:
		var st P104
		bc := Oct(b[nb+1]) + Oct(b[nb+2]) + Oct(b[nb+3])
		//Идентификатор информационного элемента (104)
		st.IDI = id
		cnt, _ := strconv.ParseInt(bc, 2, 64)
		//Количество тарифных импульсов
		st.CNT = cnt
		nb = nb + 4
		rec.P104 = st
		return nb
		//Базовая услуга (Basic service) (105)
	case 105:
		var st P105
		//Идентификатор информационного элемента (105)
		st.IDI = id
		//Базовая услуга
		st.SVC, _ = strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		//Телеслужбы
		st.TVC, _ = strconv.ParseInt(Oct(b[nb+2]), 2, 64)
		nb = nb + 3
		rec.P105 = st
		return nb
		//Дополнительная услуга у инициатора вызова (Supplementary service used by calling subscriber) (106)
	case 106:
		var st P106
		//Идентификатор информационного элемента (106)
		st.IDI = id
		//Дополнительная услуга
		st.SVC, _ = strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		nb = nb + 2
		rec.P106 = st
		return nb
		//Дополнительная услуга у вызванного абонента (Supplementary service used by called subscriber) (107)
	case 107:
		var st P107
		//Идентификатор информационного элемента (107)
		st.IDI = id
		//Дополнительная услуга
		st.SVC, _ = strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		nb = nb + 2
		rec.P107 = st
		return nb
		//Администрирование услуги абонентом (Subscriber's control input)(108)
	case 108:
		var st P108
		//Идентификатор информационного элемента (108)
		st.IDI = id
		//Тип ввода
		st.TPE = b2i(b[nb+1 : nb+2])
		//Дополнительная услуга
		st.SVC = b2i(b[nb+2 : nb+3])
		nb = nb + 3
		rec.P108 = st
		return nb
		//Последовательность символов (Dialed digits) (109)
	case 109:
		var st P109
		//Идентификатор информационного элемента (109)
		st.IDI = id
		//Число символов в последовательности
		st.CNL = b2i(b[nb+1 : nb+2])
		btb := Bts(st.CNL)
		//Последовательность символов
		st.CNC = H2c(b[nb+2 : nb+2+btb])[0:int(st.CNL)]
		nb = nb + 2 + btb
		rec.P109 = st
		return nb
		//Исходящая категория (Origin category) (110)
	case 110:
		var st P110
		//Идентификатор информационного элемента (110)
		st.IDI = id
		//Исходящая категория
		st.CAT, _ = strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		nb = nb + 2
		rec.P110 = st
		return nb
		//Тарифное направление (111)
	case 111:
		var st P111
		//Идентификатор информационного элемента (111)
		st.IDI = id
		//Тарифное направление
		st.DIR, _ = strconv.ParseInt(Oct(b[nb+1]), 2, 64)
		nb = nb + 2
		rec.P111 = st
		return nb
		//Причина безуспешного вызова (Failure cause) (112)
	case 112:
		var st P112
		//Идентификатор информационного элемента (112)
		st.IDI = id
		//Причина безуспешного вызова
		st.CFC = b2i(b[nb+1 : nb+3])
		nb = nb + 2
		rec.P112 = st
		return nb
		//Идентификация входящей соединительной линии (Incoming trunk data) (113)
	case 113:
		var st P113
		//Идентификатор информационного элемента (113)
		st.IDI = id
		//Идентификатор группы соединительных линий
		st.TGN = b2i(b[nb+1 : nb+3])
		//Идентификатор соединительной линии
		st.SLN = b2i(b[nb+3 : nb+5])
		//Идентификатор модуля
		st.MDN = b2i(b[nb+5 : nb+6])
		//Идентификатор порта
		st.PTN = b2i(b[nb+6 : nb+8])
		//Идентификатор канала
		st.CHN = b2i(b[nb+8 : nb+9])
		rec.P113 = st
		nb = nb + 9
		return nb
		//Идентификация исходящей соединительной линии (Outgoing trunk data) (114)
	case 114:
		var st P114
		//Идентификатор информационного элемента (114)
		st.IDI = id
		//Идентификатор группы соединительных линий
		st.TGN = b2i(b[nb+1 : nb+3])
		//Идентификатор соединительной линии
		st.SLN = b2i(b[nb+3 : nb+5])
		//Идентификатор модуля
		st.MDN = b2i(b[nb+5 : nb+6])
		//Идентификатор порта
		st.PTN = b2i(b[nb+6 : nb+8])
		//Идентификатор канала
		st.CHN = b2i(b[nb+8 : nb+9])
		rec.P114 = st
		nb = nb + 9
		return nb
		//Длительность вызова или использования дополнительной услуги(Call / Service duration) (115)
	case 115:
		var st P115
		//Идентификатор информационного элемента (115)
		st.IDI = id
		//Длительность вызова или использования дополнительной услуги в мс
		st.DUR = b2i(b[nb+1 : nb+5])
		rec.P115 = st
		nb = nb + 5
		return nb
		//Контрольная сумма (Checksum) (116)
	case 116:
		var st P116
		//Идентификатор информационного элемента (116)
		st.IDI = id
		//Длина информационного элемента в байтах
		st.BTL = b2i(b[nb+1 : nb+2])
		//Контрольная сумма
		st.CRC = b2i(b[nb+2 : nb+4])
		rec.P116 = st
		nb = nb + 4
		return nb
		//Оригинальный номер вызывающего абонента (Original calling party number) (119)
	case 119:
		var st P119
		//Идентификатор информационного элемента (119)
		st.IDI = id
		//Длина информационного элемента в байтах
		st.BTL = b2i(b[nb+1 : nb+2])
		//Длина списочного номера абонента
		st.CNL = b2i(b[nb+2 : nb+3])
		btb := Bts(st.CNL)
		//Оригинальный номер вызывающего абонента
		st.CNC = H2c(b[nb+3 : nb+3+btb])[0:int(st.CNL)]
		nb = nb + 3 + btb
		rec.P119 = st
		return nb
		//Причина разъединения вызова (Call release cause) (121)
	case 121:
		var st P121
		//Идентификатор информационного элемента (121)
		st.IDI = id
		//Длина информационного элемента в байтах
		st.BTL = b2i(b[nb+1 : nb+2])
		//Код причины (Cause value)
		st.COI = b2i(b[nb+2 : nb+4])
		//Локация
		st.CNC = Oct(b[nb+5])
		nb = nb + 5
		rec.P121 = st
		return nb
	}
	//Anything else
	fmt.Println("Unknown ID:", id)
	os.Exit(0)
	return 0
}
