package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"cloud.google.com/go/translate"
	"github.com/360EntSecGroup-Skylar/excelize"
	"golang.org/x/text/language"
)

var (
	sheet1 = "移动客户端文案"
)

func main() {
	f, err := excelize.OpenFile("kite.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows := f.GetRows(sheet1)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		text := row[2]
		wg := sync.WaitGroup{}
		for index, target := range targetLang {
			wg.Add(1)
			go func(index string, t string) {
				defer wg.Done()
				translated, err := translateText(t, text)
				if err != nil {
					log.Printf("translate to %s failed: %s", t, err)
					return
				}
				log.Printf("translate %s to %s", text, translated)
				f.SetCellValue(sheet1, index+strconv.Itoa(i+1), translated)
			}(index, target)
		}
		wg.Wait()
	}
	if err := f.Save(); err != nil {
		fmt.Println(err)
	}
}

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, &translate.Options{Format: translate.Text})
	if err != nil {
		return "", fmt.Errorf("Translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("Translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}
