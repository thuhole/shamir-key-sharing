package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/SSSaaS/sssa-golang"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"strings"
)

func genKey() string {
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(base32.StdEncoding.EncodeToString(randomBytes))
}

func hash(user string) string {
	h := sha256.New()
	h.Write([]byte(user))
	return hex.EncodeToString(h.Sum(nil))
}

var smtpUsername string
var smtpPass string
var smtpHost string

func sendEmail(recipient string, share string) {

	m := gomail.NewMessage()
	m.SetHeader("From", smtpUsername)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "【T大树洞】管理员密钥")

	templateData := struct {
		Key string
	}{
		Key: share,
	}

	t, err := template.ParseFiles("send_key.html")
	if err != nil {
		panic(err)
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		panic(err)
	}
	m.SetBody("text/html", buf.String())
	m.AddAlternative("text/plain", "管理员您好：\n\n这是您的密钥，请您妥善保管。\n\n"+share+"\n")
	d := gomail.NewDialer(smtpHost, 465, smtpUsername, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		panic("invalid args length")
	}
	switch args[1] {
	case "gen":

		fmt.Println("please make sure that the outbox function is disabled for this smtp server.")
		fmt.Print("please input smtp username: ")
		_, err := fmt.Scanln(&smtpUsername)
		if err != nil {
			panic(err)
		}

		fmt.Print("please input smtp password: ")
		_, err = fmt.Scanln(&smtpPass)
		if err != nil {
			panic(err)
		}

		fmt.Print("please input smtp host: ")
		_, err = fmt.Scanln(&smtpHost)
		if err != nil {
			panic(err)
		}

		var n int
		fmt.Print("please input number of key shares: ")
		_, err = fmt.Scanln(&n)
		if err != nil {
			panic(err)
		}

		var adminEmails []string
		for i := 0; i < n; i++ {
			var tmp string
			fmt.Printf("please input email %d: ", i+1)
			_, err := fmt.Scanln(&tmp)
			if err != nil {
				panic(err)
			}
			adminEmails = append(adminEmails, tmp)
		}

		key := genKey()
		fmt.Println("hash of key:", hash(key))
		shares, err := sssa.Create(n/2+1, n, key)
		if err != nil {
			panic(err)
		}

		for i, email := range adminEmails {
			sendEmail(email, shares[i])
			fmt.Println("email sent to " + email + ".")
		}
	case "decrypt":
		var n int
		fmt.Print("please input number of key shares: ")
		_, err := fmt.Scanln(&n)
		if err != nil {
			panic(err)
		}
		if n < 1 {
			panic("you can't input a negative integer")
		}

		var shares []string
		for i := 0; i < n; i++ {
			var tmp string
			fmt.Printf("please key share %d: ", i+1)
			_, err := fmt.Scanln(&tmp)
			if err != nil {
				panic(err)
			}
			if sssa.IsValidShare(tmp) {
				shares = append(shares, tmp)
			} else {
				fmt.Println("invalid share! please retry.")
				i -= 1
				continue
			}
		}

		text, err := sssa.Combine(shares)
		if err != nil {
			panic(err)
		}
		fmt.Println(text)
	default:
		fmt.Println("usage: ./shamir-key-generate gen OR ./shamir-key-generate decrypt")
	}
}
