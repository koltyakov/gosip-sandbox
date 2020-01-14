// go run ./samples/sync/ -watched ./samples/sync/watched
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip/api"
	strategy "github.com/koltyakov/gosip/auth/saml"
	"github.com/radovskyb/watcher"
)

var (
	wFolder  string
	spFolder string
	syncAll  bool
)

func main() {
	flag.StringVar(&wFolder, "watched", "", "Local folder to watch")
	flag.StringVar(&spFolder, "spFolder", "SiteAssets", "SP folder to sync to")
	flag.BoolVar(&syncAll, "syncAll", false, "Sync all files on startup")
	flag.Parse()

	if wFolder == "" {
		log.Fatalf("no watched folder is provided")
	}

	// Binding auth & API client
	configPath := "./config/private.saml.json"
	authCnfg := &strategy.AuthCnfg{}
	if err := authCnfg.ReadConfig(configPath); err != nil {
		log.Fatalf("unable to get config: %v", err)
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
			for path := range w.WatchedFiles() {
				uploadFile(sp, path)
			}
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
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil // skip 0 size files
	}
	// Ensure remote folder exists
	folderURI := getFileFolderURI(filePath)
	fileURI := getFileURI(filePath)
	if _, err := sp.Web().EnsureFolder(folderURI); err != nil {
		return err
	}
	// Upload file to document library
	file, err := sp.Web().GetFolder(folderURI).Files().Add(filepath.Base(filePath), data, true)
	if err != nil {
		return err
	}
	if file.Data().CheckOutType != 2 {
		if _, err := sp.Web().GetFile(fileURI).CheckIn("", 2); err != nil {
			return err
		}
	}
	log.Printf("📄 ✔️: %s\n", fileURI)
	return nil
}

func deleteFile(sp *api.SP, filePath string) error {
	fileURI := getFileURI(filePath)
	if err := sp.Web().GetFile(fileURI).Recycle(); err == nil {
		log.Printf("📄 ❌: %s\n", fileURI)
	}
	return nil
}

func createFolder(sp *api.SP, folderPath string) error {
	folderURI := getFolderURI(folderPath)
	if _, err := sp.Web().EnsureFolder(folderURI); err != nil {
		return err
	}
	// log.Printf("Folder is created: %s\n", folderURI)
	return nil
}

func deleteFolder(sp *api.SP, folderPath string) error {
	folderURI := getFolderURI(folderPath)
	if err := sp.Web().GetFolder(folderURI).Recycle(); err != nil {
		return err
	}
	log.Printf("📁 ❌: %s\n", folderURI)
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
