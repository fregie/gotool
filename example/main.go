package main

import (
	"fmt"
	"log"
	"unicode"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/fregie/gotool/ip138"
)

const (
	sheet = "通讯录总"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cli := ip138.NewIP138("*")
	// cli = cli
	file, err := excelize.OpenFile("./mobile.xlsx")
	if err != nil {
		log.Print(err)
		return
	}
	rows := file.GetRows(sheet)
	for i, row := range rows {
		if len(row) <= 0 {
			continue
		}
		if len(row) >= 2 && row[1] != "" {
			continue
		}
		r := make([]byte, 0, len(row[0]))
		for _, s := range row[0] {
			if !unicode.IsNumber(s) {
				continue
			}
			r = append(r, byte(s))
		}
		if len(r) < 11 || r[len(r)-11] != byte('1') {
			continue
		}
		number := string(r[len(r)-11:])
		fmt.Print(string(number) + "\t")
		info, err := cli.Mobile(number)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("%s-%s-%s", info.Province, info.City, info.Provider)
		file.SetCellValue(sheet, fmt.Sprintf("B%d", i+1), info.Province)
		file.SetCellValue(sheet, fmt.Sprintf("C%d", i+1), info.City)
		file.SetCellValue(sheet, fmt.Sprintf("D%d", i+1), info.Provider)

		fmt.Print("\n")
	}
	if err := file.Save(); err != nil {
		fmt.Println(err)
	}
}
