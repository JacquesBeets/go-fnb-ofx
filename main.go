package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aclindsa/ofxgo"
	"github.com/joho/godotenv"
	"github.com/playwright-community/playwright-go"
)


func assertErrorToNilf(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func Unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer r.Close()

    for _, f := range r.File {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer rc.Close()

        fpath := filepath.Join(dest, f.Name)
        if f.FileInfo().IsDir() {
            os.MkdirAll(fpath, os.ModePerm)
        } else {
            var fdir string
            if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
                fdir = fpath[:lastIndex]
            }

            err = os.MkdirAll(fdir, os.ModePerm)
            if err != nil {
                log.Fatal(err)
                return err
            }
            f, err := os.OpenFile(
                fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer f.Close()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
    }



    return nil
}

func RenameFile() error {
		// Read the contents of the unzipped directory
		files, err := os.ReadDir("./unzipped")
		if err != nil {
			log.Fatalf("Failed to read directory: %v", err)
			return err
		}
	
	
		// Check if there's at least one file in the directory
		if len(files) > 0 {
			// Get the first file
			file := files[0]
	
			// Check if the file is not a directory
			if !file.IsDir() {
				// Rename the file
				oldPath := "./unzipped/" + file.Name()
				newPath := "./unzipped/fnb.ofx"
				err = os.Rename(oldPath, newPath)
				if err != nil {
					fmt.Printf("Failed to rename file: %v", err)
					return err
				}
			}
		}

		return nil
}

func ParseOFX() error {
	f, err := os.Open("./unzipped/fnb.ofx")
	if err != nil {
		fmt.Printf("could not open OFX file: %v", err) 
		return err
	}
	defer f.Close()

	resp, err := ofxgo.ParseResponse(f)
	if err != nil {
		fmt.Printf("could not parse OFX file: %v", err)
		return err
	}


	// Log the account number and balance for each bank account
	if len(resp.Bank) > 0 {
		bankMessage := resp.Bank[0]
		if stmt, ok := bankMessage.(*ofxgo.StatementResponse); ok {
			// Access the TransactionList
			transactions := stmt.BankTranList
			for _, transaction := range transactions.Transactions {
				fmt.Printf("Transaction type: %s\n", transaction.TrnType)
				fmt.Printf("Transaction date: %s\n", transaction.DtPosted.Format("2006-01-02"))
				fmt.Printf("Transaction amount: %v\n", transaction.TrnAmt)
				fmt.Printf("Transaction ID: %s\n", transaction.FiTID)
				fmt.Printf("Transaction name: %s\n", transaction.Name)
				fmt.Printf("Transaction memo: %s\n", transaction.Memo)
				fmt.Printf("Transaction check number: %s\n", transaction.CheckNum)
				fmt.Printf("Transaction ref number: %s\n", transaction.RefNum)
				fmt.Printf("Transaction sic: %v\n", transaction.SIC)
				fmt.Printf("Transaction mcc: %v\n", transaction.DtPosted.Time)
				fmt.Printf("=======================================================================================\n")
			}
		}
	}

	return nil
}

func cleanUp() error {
	// Remove unzipped directory
	err := os.RemoveAll("./unzipped")
	if err != nil {
		fmt.Printf("could not remove unzipped directory: %v", err)
		return err
	}

	// Remove downloads directory
	// err = os.RemoveAll("./downloads")
	// if err != nil {
	// 	fmt.Printf("could not remove downloads directory: %v", err)
	// 	return err
	// }

	return nil
}

func main() {
	err := godotenv.Load("./.env")
	assertErrorToNilf("Error loading .env file", err)

	var usern string
	var pass string
	var website string
	var waitForLogin string
	var waitForLogout string
	usern = os.Getenv("USERN")
	pass = os.Getenv("PASSWORD")
	website = os.Getenv("WEBSITE")
	waitForLogin = os.Getenv("WAIT_FOR_LOGIN")
	waitForLogout = os.Getenv("WAIT_FOR_LOGOUT")


	// Launch Playwright
	pw, err := playwright.Run()
	assertErrorToNilf("could not launch playwright: %w", err)

	// Launch Browser
	browser, err := pw.Chromium.Launch()

	// Luanch Browser with UI
	// browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
	// 	Headless: playwright.Bool(false),
	// })
	assertErrorToNilf("could not launch Chromium: %w", err)

	// Create New Page
	page, err := browser.NewPage()

	// Create New Page with Video Recording
	// page, err := browser.NewPage(playwright.BrowserNewPageOptions{
	// 	RecordVideo: &playwright.RecordVideo{
	// 		Dir: "videos/",
	// 	},
	// })
	assertErrorToNilf("could not create page: %w", err)

	// Goto Website
	_, err = page.Goto(website)
	assertErrorToNilf("could not goto: %w", err)

	_, err = page.Video().Path()
	assertErrorToNilf("failed to get video path: %v", err)


	// Fill in Username
	assertErrorToNilf("could not type: %v", page.Locator("input#user").Fill(usern))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Fill in Password
	assertErrorToNilf("could not type: %v", page.Locator("input#pass").Fill(pass))

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click Login
	assertErrorToNilf("could not press: %v", page.Locator("#OBSubmit").Press("Enter"))

	//WaitForLogin to complete
	frame := page.MainFrame()
	_ = frame.WaitForURL(waitForLogin)

	time.Sleep(5 * time.Second) // Wait for 5 seconds

	// Click on Accounts
	assertErrorToNilf("could not click: %v", page.Locator("#shortCutLinks > span:nth-child(1)").Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click on Balance
	assertErrorToNilf("could not click: %v", page.Locator("#tabelRow_6 .group3 .col4 a").Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click More
	assertErrorToNilf("could not click: %v", page.Locator("#footerButtonsContainer > div:nth-child(1) a").Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Click on Download
	assertErrorToNilf("could not click: %v", page.Locator("#tableActionButtons .downloadButton").Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	// Open Dropdown
	assertErrorToNilf("could not open dropdown: %v", page.Locator("#downloadFormat_dropId").Click())

	time.Sleep(3 * time.Second) // Wait for 3 seconds
	
	// Select OFX
	assertErrorToNilf("could not select OFX: %v", page.Locator(`[data-value="ofx"]`).Click())
	// assertErrorToNilf("could not select OFX: %v", page.Locator("ul.dropdown-content li:last-child").Click())
	// assertErrorToNilf("could not select OFX: %v", page.Locator("//*[@id="downloadFormat_parent"]/div[2]/div[3]/ul/li[6]").Click())  // X-PATH

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	//Download
	download, err := page.ExpectDownload(func() error {
		return page.Locator("#eziPannelButtonsWrapper #mainDownloadBtn").Click()
	})
	assertErrorToNilf("could not download: %w", err)

	// Save download to file
	err = download.SaveAs("./downloads/fnb_ofx.zip") // Save to current directory
	assertErrorToNilf("could not save download to file: %w", err)

	time.Sleep(5 * time.Second) // Wait for 5 seconds

	// Logout
	assertErrorToNilf("could not logout: %v", page.Locator("#headerButton_").Click())

	// WaitForLogout to complete
	_ = frame.WaitForURL(waitForLogout)

	time.Sleep(3 * time.Second) // Wait for 3 seconds

	assertErrorToNilf("could not close browser: %w", browser.Close())
	assertErrorToNilf("could not stop Playwright: %w", pw.Stop())


	// Unzip
	err = Unzip("./downloads/fnb_ofx.zip", "./unzipped")
	assertErrorToNilf("could not unzip: %w", err)

	// Rename File
	err = RenameFile()
	assertErrorToNilf("could not rename file: %w", err)

	// Parse OFX
	err = ParseOFX()
	assertErrorToNilf("could not parse OFX: %w", err)

	// Clean Up
	err = cleanUp()
	assertErrorToNilf("could not clean up: %w", err)

	fmt.Println("Completed")

	os.Exit(0)
}
