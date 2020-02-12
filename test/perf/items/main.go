package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/koltyakov/gosip"
	"github.com/koltyakov/gosip-sandbox/samples/dynauth"
	"github.com/koltyakov/gosip/api"
	"github.com/koltyakov/gosip/csom"
)

func main() {
	configPath := flag.String("config", "./config/private.json", "private config path")
	strategy := flag.String("strategy", "saml", "auth strategy")
	isSPO := flag.Bool("spo", false, "is SharePoint Online")
	flag.Parse()

	// Wrap auth
	authCnfg, err := dynauth.NewAuthCnfg(*strategy, *configPath)
	if err != nil {
		log.Fatal(err)
	}
	authCnfg.ReadConfig(*configPath)
	client := &gosip.SPClient{AuthCnfg: authCnfg}
	sp := api.NewSP(client)
	httpClient := api.NewHTTPClient(client)

	// Prepare list
	listName := strings.Replace(uuid.New().String(), "-", "", -1)
	list := sp.Web().Lists().GetByTitle(listName)
	if _, err := sp.Web().Lists().Add(listName, nil); err != nil {
		log.Fatal(err)
	}
	defer list.Delete()

	// Cache entity
	entType, err := list.GetEntityType()
	if err != nil {
		log.Fatal(err)
	}

	var start time.Time

	// Add series in a sequence (OData Add)
	start = time.Now()
	fmt.Println("Add 100 items in a sequence (OData Add)")
	for i := 1; i <= 100; i++ {
		metadata := make(map[string]interface{})
		metadata["__metadata"] = map[string]string{"type": entType}
		metadata["Title"] = fmt.Sprintf("Item %d", i)
		body, _ := json.Marshal(metadata)
		if _, err := list.Items().Add(body); err != nil {
			fmt.Printf("Item adding error: %v\n", err)
		}
	}
	fmt.Printf("  took: %s\n", time.Now().Sub(start))

	// Add series in parallel (OData Add)
	start = time.Now()
	var wg sync.WaitGroup
	fmt.Println("Add 100 items in parallel (OData Add)")
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			metadata := make(map[string]interface{})
			metadata["__metadata"] = map[string]string{"type": entType}
			metadata["Title"] = fmt.Sprintf("Item %d", i)
			body, _ := json.Marshal(metadata)
			if _, err := list.Items().Add(body); err != nil {
				fmt.Printf("Item adding error: %v\n", err)
			}
		}(i, &wg)
	}
	wg.Wait()
	fmt.Printf("  took: %s\n", time.Now().Sub(start))

	if *isSPO { // following methods are supported in SPO

		// Add series in a sequence (AddValidateUsingPath)
		start = time.Now()
		fmt.Println("Add 100 items in a sequence (AddValidateUsingPath)")
		for i := 1; i <= 100; i++ {
			metadata := map[string]string{"Title": fmt.Sprintf("Item %d", i)}
			if _, err := list.Items().AddValidate(metadata, nil); err != nil {
				fmt.Printf("Item adding error: %v\n", err)
			}
		}
		fmt.Printf("  took: %s\n", time.Now().Sub(start))

		// Add series in parallel (AddValidateUsingPath)
		start = time.Now()
		fmt.Println("Add 100 items in parallel (AddValidateUsingPath)")
		for i := 1; i <= 100; i++ {
			wg.Add(1)
			go func(i int, wg *sync.WaitGroup) {
				defer wg.Done()
				metadata := map[string]string{"Title": fmt.Sprintf("Item %d", i)}
				if _, err := list.Items().AddValidate(metadata, nil); err != nil {
					fmt.Printf("Item adding error: %v\n", err)
				}
			}(i, &wg)
		}
		wg.Wait()
		fmt.Printf("  took: %s\n", time.Now().Sub(start))

	}

	// Add series in a sequence (CSOM)
	start = time.Now()
	fmt.Println("Add 100 items in a sequence (CSOM)")
	for i := 1; i <= 100; i++ {
		b := csom.NewBuilder()
		b.AddObject(csom.NewObjectProperty("Web"), nil)
		b.AddObject(csom.NewObjectProperty("Lists"), nil)
		b.AddObject(csom.NewObjectMethod("GetByTitle", []string{
			fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, listName),
		}), nil)
		b.AddObject(csom.NewObject(`
			<Method Id="{{.ID}}" ParentId="{{.ParentID}}" Name="AddItem">
				<Parameters>
					<Parameter TypeId="{54cdbee5-0897-44ac-829f-411557fa11be}">
						<Property Name="FolderUrl" Type="Null" />
						<Property Name="LeafName" Type="Null" />
						<Property Name="UnderlyingObjectType" Type="Number">0</Property>
					</Parameter>
				</Parameters>
			</Method>
		`), nil)
		b.AddAction(csom.NewActionMethod("SetFieldValue", []string{
			`<Parameter Type="String">Title</Parameter>`,
			fmt.Sprintf(`<Parameter Type="String">Item %d</Parameter>`, i),
		}), nil)
		b.AddAction(csom.NewAction(`<Method Name="Update" Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), nil)

		csomXML, err := b.Compile()
		if err != nil {
			fmt.Printf("Item adding error: %v\n", err)
		}

		if _, err := httpClient.ProcessQuery(authCnfg.GetSiteURL(), bytes.NewBuffer([]byte(csomXML)), nil); err != nil {
			fmt.Printf("Item adding error: %v\n", err)
		}
	}
	fmt.Printf("  took: %s\n", time.Now().Sub(start))

	// Add series in batch (CSOM batch)
	start = time.Now()
	fmt.Println("Add 100 items (CSOM batch)")
	b := csom.NewBuilder()
	b.AddObject(csom.NewObjectProperty("Web"), nil)
	b.AddObject(csom.NewObjectProperty("Lists"), nil)
	cListObj, _ := b.AddObject(csom.NewObjectMethod("GetByTitle", []string{
		fmt.Sprintf(`<Parameter Type="String">%s</Parameter>`, listName),
	}), nil)
	for i := 1; i <= 100; i++ {
		cAddItemObj, _ := b.AddObject(csom.NewObject(`
			<Method Id="{{.ID}}" ParentId="{{.ParentID}}" Name="AddItem">
				<Parameters>
					<Parameter TypeId="{54cdbee5-0897-44ac-829f-411557fa11be}">
						<Property Name="FolderUrl" Type="Null" />
						<Property Name="LeafName" Type="Null" />
						<Property Name="UnderlyingObjectType" Type="Number">0</Property>
					</Parameter>
				</Parameters>
			</Method>
		`), cListObj)
		b.AddAction(csom.NewActionMethod("SetFieldValue", []string{
			`<Parameter Type="String">Title</Parameter>`,
			fmt.Sprintf(`<Parameter Type="String">Item %d</Parameter>`, i),
		}), cAddItemObj)
		b.AddAction(csom.NewAction(`<Method Name="Update" Id="{{.ID}}" ObjectPathId="{{.ObjectID}}" />`), cAddItemObj)
	}
	csomXML, err := b.Compile()
	if err != nil {
		fmt.Printf("Items adding error: %v\n", err)
	}
	if _, err := httpClient.ProcessQuery(authCnfg.GetSiteURL(), bytes.NewBuffer([]byte(csomXML)), nil); err != nil {
		fmt.Printf("Items adding error: %v\n", err)
	}
	fmt.Printf("  took: %s\n", time.Now().Sub(start))

}
