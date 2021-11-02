// go run ./samples/spupload/ -localFolder ./samples/spupload/upload -concurrency 10
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
	"github.com/radovskyb/watcher"

	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
)

var (
	localFolder string
	spFolder    string
	skipPublish bool
	concurrency int
)

func main() {
	strategy := flag.String("strategy", "saml", "Auth strategy")
	config := flag.String("config", "./config/private.json", "Config path")

	flag.StringVar(&localFolder, "localFolder", "", "Local folder to watch")
	flag.IntVar(&concurrency, "concurrency", 25, "Parallel upload factor")
	flag.StringVar(&spFolder, "spFolder", "Shared Documents", "SP folder to sync to")
	flag.BoolVar(&skipPublish, "skipPublish", false, "Skip publishing a major version")

	flag.Parse()

	if localFolder == "" {
		log.Fatalf("no upload source folder is provided")
	}

	// Binding auth & API client
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *config)
	if err != nil {
		log.Fatal(err)
	}
	client := &gosip.SPClient{AuthCnfg: authCnfg}

	client.Hooks = &gosip.HookHandlers{
		OnError: func(e *gosip.HookEvent) {
			if e.StatusCode == 429 {
				log.Printf("❌ : Got throttled, now waiting...\n")
			}
		},
	}

	sp := api.NewSP(client)

	web, err := sp.Web().Select("Title").Get()
	if err != nil {
		log.Fatalf("can't connect to site, %s\n", err)
	}

	log.Printf("Upload source: %s\n", localFolder)
	log.Printf("Upload target: %s, %s\n", web.Data().Title, spFolder)

	run(sp)
}

func run(sp *api.SP) {
	w := watcher.New()
	done := make(chan bool)

	go func() {
		start := time.Now()
		filesNum := 0
		slots := concurrency
		for path, file := range w.WatchedFiles() {
			if !file.IsDir() {
				filesNum++
				for slots == 0 {
					time.Sleep(10 * time.Microsecond)
				}
				slots = slots - 1
				go func(p string, s *int) {
					if err := uploadFile(sp, p); err != nil {
						log.Printf("❌ : %s: %s\n", p, err)
					}
					*s = *s + 1
				}(path, &slots)
			}
		}
		for slots != concurrency {
			time.Sleep(10 * time.Microsecond)
		}
		w.Close()
		log.Printf("📄 🏁 : Upload of %d file(s) took %s\n", filesNum, time.Since(start))
		done <- true
	}()

	if err := w.AddRecursive(localFolder); err != nil {
		log.Fatalln(err)
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

	<-done
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
		if strings.Contains(err.Error(), "System.IO.DirectoryNotFoundException") {
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
	// Check in a file if it was checked out
	if file.Data().CheckOutType != 2 {
		if _, err := sp.Web().GetFile(fileURI).CheckIn("", 2); err != nil {
			return err
		}
	}
	// Publish a file if it has minor version
	if !skipPublish && file.Data().MinorVersion != 0 {
		if _, err := sp.Web().GetFile(fileURI).Publish(""); err != nil {
			return err
		}
	}
	log.Printf("📄 ✔️ : %s (%s)\n", fileURI, time.Since(start))
	return nil
}

func getFileFolderURI(filePath string) string {
	folderPath := filepath.Dir(filePath)
	lFolder, _ := filepath.Abs(localFolder)
	relPath, _ := filepath.Rel(lFolder, folderPath)
	relPath = strings.Replace(relPath, "\\", "/", -1)
	return spFolder + "/" + relPath
}

func getFileURI(filePath string) string {
	lFolder, _ := filepath.Abs(localFolder)
	relPath, _ := filepath.Rel(lFolder, filePath)
	relPath = strings.Replace(relPath, "\\", "/", -1)
	return spFolder + "/" + relPath
}
