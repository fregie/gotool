package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"unicode/utf8"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/api/option"
)

const (
	packageName = "com.flash.svpn"
)

func main() {
	ctx := context.Background()
	service, err := androidpublisher.NewService(ctx, option.WithCredentialsJSON([]byte(cred)))
	if err != nil {
		log.Printf("init service failed:%s", err)
		return
	}
	insertCall := service.Edits.Insert(packageName, &androidpublisher.AppEdit{})
	appEdit, err := insertCall.Do()
	if err != nil {
		log.Printf("insert edits failed:%s", err)
		return
	}
	listing := service.Edits.Listings
	wg := sync.WaitGroup{}
	for listLang, transLang := range targetLang {
		wg.Add(1)
		go func(listLang, transLang string) {
			defer wg.Done()
			newListing, err := translateListing(&originList, &originListShort, transLang)
			if err != nil {
				log.Printf("translate to [%s] failed: %s", transLang, err)
				return
			}
			// fmt.Print(newListing.Title + "\n")
			// fmt.Print(newListing.ShortDescription + "\n")
			// fmt.Print(newListing.FullDescription + "\n")
			_, err = listing.Update(packageName, appEdit.Id, listLang, newListing).Do()
			if err != nil {
				log.Printf("update [%s] failed:%s", listLang, err)
				return
			}
		}(listLang, transLang)
	}
	wg.Wait()
	commitCall := service.Edits.Commit(packageName, appEdit.Id)
	_, err = commitCall.Do()
	if err != nil {
		log.Printf("commit failed: %s", err)
		return
	}
}

var ETOOLONG error = errors.New("too long")

func translateListing(origin *androidpublisher.Listing, short *androidpublisher.Listing, targetLanguage string) (new *androidpublisher.Listing, err error) {
	new = &androidpublisher.Listing{}
	new.Title, err = translateText(targetLanguage, origin.Title)
	if err != nil {
		return
	}
	if utf8.RuneCountInString(new.Title) > 50 {
		new.Title, err = translateText(targetLanguage, short.Title)
		if err != nil {
			return
		}
		if utf8.RuneCountInString(new.Title) > 50 {
			// log.Printf("Title too long: %d", utf8.RuneCountInString(new.Title))
			// fmt.Print(new.Title + "\n")
			err = errors.New("title too long")
			return
		}
	}
	new.ShortDescription, err = translateText(targetLanguage, origin.ShortDescription)
	if err != nil {
		return
	}
	if utf8.RuneCountInString(new.ShortDescription) > 80 {
		new.ShortDescription, err = translateText(targetLanguage, short.ShortDescription)
		if err != nil {
			return
		}
		if utf8.RuneCountInString(new.ShortDescription) > 80 {
			// log.Printf("ShortDescription too long: %d", utf8.RuneCountInString(new.ShortDescription))
			// fmt.Print(new.ShortDescription + "\n")
			err = errors.New("ShortDescription too long")
			return
		}
	}
	new.FullDescription, err = translateText(targetLanguage, origin.FullDescription)
	if err != nil {
		return
	}
	if utf8.RuneCountInString(new.FullDescription) > 4000 {
		new.FullDescription, err = translateText(targetLanguage, short.FullDescription)
		if err != nil {
			return
		}
		if utf8.RuneCountInString(new.FullDescription) > 4000 {
			// log.Printf("FullDescription too long: %d", utf8.RuneCountInString(new.FullDescription))
			// fmt.Print(new.FullDescription + "\n")
			err = errors.New(fmt.Sprintf("FullDescription too long [%d]", utf8.RuneCountInString(new.FullDescription)))
			return
		}
	}
	return
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
