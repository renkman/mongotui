// seed.go
package testdata

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const description string = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua."

func Seed(ctx context.Context, uri string, count int) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetConnectTimeout(10*time.Second))
	if err != nil {
		log.Panicf("Connect to database %s failed: %s", uri, err.Error())
		os.Exit(1)
	}

	collection := CreateCollection(client, "test", "data")
	GenerateData(ctx, collection, count)
}

func CreateCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection {
	collection := client.Database(databaseName).Collection(collectionName)
	return collection
}

func GenerateData(ctx context.Context, collection *mongo.Collection, count int) {
	for i := 0; i < count; i++ {
		if rand.Intn(2) == 0 {
			document := generateDocument(i)
			_, err := collection.InsertOne(ctx, document)
			if err != nil {
				log.Panicf("Insert document %v to collection %s failed: %s", document, err.Error())
				os.Exit(2)
			}
			continue
		}
		document := generateListDocument(i)
		_, err := collection.InsertOne(ctx, document)
		if err != nil {
			log.Panicf("Insert document %v to collection %s failed: %s", document, err.Error())
			os.Exit(3)
		}
	}
}

func generateDocument(number int) Document {
	name := fmt.Sprintf("Document %d", number)
	subNumber := rand.Intn(1000)
	stuff := generateStuff()
	subDocument := SubDocument{subNumber, stuff}

	document := Document{name, description, time.Now(), subDocument}
	return document
}

func generateListDocument(number int) ListDocument {
	name := fmt.Sprintf("List Document %d", number)

	list := make([]SubDocument, rand.Intn(101))
	for i := 0; i < len(list); i++ {
		subNumber := rand.Intn(1000)
		stuff := generateStuff()
		subDocument := SubDocument{subNumber, stuff}
		list[i] = subDocument
	}

	document := ListDocument{name, time.Now(), list}
	log.Printf("List %d: %v", number, list)

	return document
}

func generateStuff() string {
	stuff := []string{
		"Foo",
		"Bar",
		"Baz",
		"Bla",
		"Foobar",
		"Lorem ipsum",
		"Amiga rules",
		"Apple sucks",
		"42",
		"Stuff",
	}
	index := rand.Intn(len(stuff))
	return stuff[index]
}

// func seedWindows(client *mongo.Client, ctx context.Context) {
// 	collection := client.Database("renkbench").Collection("windows")

// 	result, err := collection.InsertOne(ctx, model.Window{nil, 0, -1, model.WindowMetaInfo{"Renkbench"}, nil, &[]model.Icon{
// 		model.Icon{1, "Renkbench", model.Image{"workbench.png", 35, 30}, model.Image{"workbench_selected.png", 35, 30}},
// 		model.Icon{2, "Double click me!", model.Image{"disk.png", 32, 32}, model.Image{"disk_selected.png", 32, 32}}}})
// 	if err != nil {
// 		log.Fatal(err)
// 		return
// 	}

// 	id, ok := result.InsertedID.(primitive.ObjectID)
// 	if !ok {
// 		log.Fatalf("result.InsertedID with value %v is not of type primitive.ObjectID", result.InsertedID)
// 		return
// 	}

// 	data := []interface{}{
// 		model.Window{&id, 1, 0, model.WindowMetaInfo{"Renkbench"}, nil, &[]model.Icon{
// 			model.Icon{3, "Edit", model.Image{"notepad.png", 32, 32}, model.Image{"notepad_selected.png", 32, 32}},
// 			model.Icon{4, "Drawer", model.Image{"drawer.png", 71, 31}, model.Image{"drawer_selected.png", 34, 73}}}},
// 		model.Window{&id, 2, 0, model.WindowMetaInfo{"Note!"}, &model.Content{getRef("Hello again!"),
// 			&[]model.Article{model.Article{"Renkbench relaunch", `Next release: During 2022 I had the idea of switching the backend from Node.js to Go. Since I started my <a href="https://github.com/renkman/mongotui" target="_blank">MongoTUI MongoDB client</a> in 2020, I felt like programming more stuff in Go. So - here it is. And with this move, I replaced the static file based JSON-content with a MongoDB.<br /><br />The source code of this web app is available on my <a href="https://github.com/renkman/Renkbench" target="_blank">Github repository</a>.`}}, nil}, nil}}
// 	manyResult := insertManyCollection(data, collection, client, ctx)

// 	id, ok = manyResult.InsertedIDs[0].(primitive.ObjectID)
// 	if !ok {
// 		log.Fatalf("result.InsertedIDs[0] with value %v is not of type primitive.ObjectID", manyResult.InsertedIDs[0])
// 		return
// 	}

// 	data = []interface{}{
// 		model.Window{&id, 3, 1, model.WindowMetaInfo{"Edit"}, &model.Content{nil, nil, getRef(`<div class="textbox" tabindex="0"></div>`)}, nil},
// 		model.Window{&id, 4, 1, model.WindowMetaInfo{"Drawer"}, nil, &[]model.Icon{
// 			model.Icon{5, "Document", model.Image{"document.png", 40, 40}, model.Image{"document_selected.png", 40, 40}}}}}
// 	manyResult = insertManyCollection(data, collection, client, ctx)

// 	id, ok = manyResult.InsertedIDs[0].(primitive.ObjectID)
// 	if !ok {
// 		log.Fatalf("result.InsertedIDs[0] with value %v is not of type primitive.ObjectID", manyResult.InsertedIDs[0])
// 		return
// 	}

// 	result, err = collection.InsertOne(ctx, model.Window{&id, 5, 4, model.WindowMetaInfo{"Document"}, &model.Content{getRef("Hello again!"),
// 		&[]model.Article{model.Article{"Document", `This is just a document containing this text.`}}, nil}, nil})
// }

// func seedMenu(client *mongo.Client, ctx context.Context) {

// 	data := []interface{}{
// 		model.Menu{"Renkbench", []model.MenuEntry{
// 			model.MenuEntry{"Open", "openWindow", `[ {"property" : "isSelected", "value" : true } ]`},
// 			model.MenuEntry{"Close", "closeWindow", `[ {"property" : "isSelected", "value" : true }, {"property" : "isOpened", "value" : true } ]`},
// 			model.MenuEntry{"Duplicate", "duplicate", `[ {"property" : "isSelected", "value" : true } ]`},
// 			model.MenuEntry{"Rename", "rename", `[ {"property" : "isSelected", "value" : true } ]`},
// 			model.MenuEntry{"Info", "info", `[ {"property" : "isSelected", "value" : true } ]`},
// 			model.MenuEntry{"Discard", "discard", `[ {"property" : "isSelected", "value" : true }, {"property" : "pid", "value" : 0, "operand": "greaterThan" }]`}}},
// 		model.Menu{"Disk", []model.MenuEntry{
// 			model.MenuEntry{"Empty trash", "emptyTrash", `[ {"property" : "isSelected", "value" : true },  {"property" : "isTrashcan", "value" : true } ]`},
// 			model.MenuEntry{"Initialize", "initialize", `[ {"property" : "isSelected", "value" : true },  {"property" : "pid", "value" : 0 } ]`}}},
// 		model.Menu{"Special", []model.MenuEntry{
// 			model.MenuEntry{"Clean up", "cleanUp", `[ {"property" : "isSelected", "value" : true },  {"property" : "id", "value" : 1 } ]`},
// 			model.MenuEntry{"Last error", "lastError", "true"},
// 			model.MenuEntry{"Redraw", "redraw", "true"},
// 			model.MenuEntry{"Snapshot", "snapshot", `[ {"property" : "isSelected", "value" : true } ]`},
// 			model.MenuEntry{"Version", "version", "true"}}}}

// 	insertManyName(data, "menu", client, ctx)
// }

// func insertManyName(data []interface{}, collectionName string, client *mongo.Client, ctx context.Context) *mongo.InsertManyResult {
// 	collection := client.Database("renkbench").Collection(collectionName)
// 	return insertManyCollection(data, collection, client, ctx)
// }

// func insertManyCollection(data []interface{}, collection *mongo.Collection, client *mongo.Client, ctx context.Context) *mongo.InsertManyResult {
// 	result, err := collection.InsertMany(ctx, data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return result
// }

// func getRef(value string) *string {
// 	return &value
// }
