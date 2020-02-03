// go run ./samples/sync/ -watch ./samples/sync/watched
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
	"github.com/koltyakov/gosip/api"
	"github.com/radovskyb/watcher"
)

var (
	wFolder  string
	spFolder string
	syncAll  bool
)

func main() {
	strategy := flag.String("strategy", "saml", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")
	flag.StringVar(&wFolder, "watch", "", "Local folder to watch")
	flag.StringVar(&spFolder, "spFolder", "SiteAssets", "SP folder to sync to")
	flag.BoolVar(&syncAll, "syncAll", false, "Sync all files on startup")

	flag.Parse()

	if wFolder == "" {
		log.Fatalf("no watched folder is provided")
	}

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatal(err)
	}
	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)

	web, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatalf("can't connect to site, %s\n", err)
	}

	log.Printf("Watching folder: %s\n", wFolder)
	log.Printf("Sync target: %s, %s\n", web.Data().Title, spFolder)

	watch(sp)
}

func watch(sp *api.SP) {
	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if err := sync(sp, event); err != nil {
					log.Printf("%s\n", err)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if syncAll {
		go func() {
			start := time.Now()
			errors := []error{}
			filesNum := 0
			for path, file := range w.WatchedFiles() {
				if len(errors) > 10 {
					log.Printf("❌ : Too many errors, skipping upload...\n")
					break
				}
				if !file.IsDir() {
					if err := uploadFile(sp, path); err != nil {
						errors = append(errors, err)
						log.Printf("%s: %s\n", path, err)
					}
					filesNum++
				}
			}
			log.Printf("📄 🏁 : Full sync of %d file(s) in %s\n", filesNum, time.Since(start))
		}()
	}

	if err := w.AddRecursive(wFolder); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func sync(sp *api.SP, event watcher.Event) error {
	// Sample shows only basic events flow
	// fmt.Printf("%s\n", event)

	// Files sync
	if !event.IsDir() {
		if event.Op.String() == "WRITE" || event.Op.String() == "CREATE" {
			return uploadFile(sp, event.Path)
		}
		if event.Op.String() == "REMOVE" {
			return deleteFile(sp, event.Path)
		}
	}

	// Folders sync
	if event.IsDir() {
		if event.Op.String() == "CREATE" {
			return createFolder(sp, event.Path)
		}
		if event.Op.String() == "REMOVE" {
			return deleteFolder(sp, event.Path)
		}
	}

	return nil
}

func uploadFile(sp *api.SP, filePath string) error {
	start := time.Now()
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil // skip 0 size files
	}
	folderURI := getFileFolderURI(filePath)
	fileURI := getFileURI(filePath)
	// Upload file to document library
	// optimistic strategy, is faster if not to check folder exists each time and create folders only on errors
	files := sp.Web().GetFolder(folderURI).Files()
	file, err := files.Add(filepath.Base(filePath), data, true)
	if err != nil {
		if strings.Index(err.Error(), "System.IO.DirectoryNotFoundException") != -1 {
			// Create remote folder
			if _, err := sp.Web().EnsureFolder(folderURI); err != nil {
				return err
			}
			log.Printf("📄 ✔️ : %s (%s)\n", folderURI, time.Since(start))
			// Another attempt after a folder(s) is/are created
			file, err = files.Add(filepath.Base(filePath), data, true)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
	}
	// Check in a file if it was checke out
	if file.Data().CheckOutType != 2 {
		if _, err := sp.Web().GetFile(fileURI).CheckIn("", 2); err != nil {
			return err
		}
	}
	log.Printf("📄 ✔️ : %s (%s)\n", fileURI, time.Since(start))
	return nil
}

func deleteFile(sp *api.SP, filePath string) error {
	start := time.Now()
	fileURI := getFileURI(filePath)
	if err := sp.Web().GetFile(fileURI).Recycle(); err != nil {
		// Ignore file does not exist errors
		if strings.Index(err.Error(), "-2146232832, Microsoft.SharePoint.SPException") == -1 {
			return err
		}
	} else {
		log.Printf("📄 ❌ : %s (%s)\n", fileURI, time.Since(start))
	}
	return nil
}

func createFolder(sp *api.SP, folderPath string) error {
	folderURI := getFolderURI(folderPath)
	if _, err := sp.Web().EnsureFolder(folderURI); err != nil {
		return err
	}
	return nil
}

func deleteFolder(sp *api.SP, folderPath string) error {
	start := time.Now()
	folderURI := getFolderURI(folderPath)
	if err := sp.Web().GetFolder(folderURI).Recycle(); err != nil {
		// Ignore folder does not exist errors
		if strings.Index(err.Error(), "404 Not Found") == -1 {
			return err
		}
	} else {
		log.Printf("📁 ❌ : %s (%s)\n", folderURI, time.Since(start))
	}
	return nil
}

func getFolderURI(folderPath string) string {
	localFolder, _ := filepath.Abs(wFolder)
	relPath, _ := filepath.Rel(localFolder, folderPath)
	relPath = strings.Replace(relPath, "\\", "/", -1)
	return spFolder + "/" + relPath
}

func getFileFolderURI(filePath string) string {
	folderPath := filepath.Dir(filePath)
	localFolder, _ := filepath.Abs(wFolder)
	relPath, _ := filepath.Rel(localFolder, folderPath)
	relPath = strings.Replace(relPath, "\\", "/", -1)
	return spFolder + "/" + relPath
}

func getFileURI(filePath string) string {
	localFolder, _ := filepath.Abs(wFolder)
	relPath, _ := filepath.Rel(localFolder, filePath)
	relPath = strings.Replace(relPath, "\\", "/", -1)
	return spFolder + "/" + relPath
}
