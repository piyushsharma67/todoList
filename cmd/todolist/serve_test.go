package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	sqlitedb "go.altair.com/todolist/pkg/db"
	"go.altair.com/todolist/pkg/structs"
	"go.altair.com/todolist/pkg/todolist"
	"go.altair.com/todolist/pkg/todolist/store"
)

var testingT *testing.T

func TestTodoServe(t *testing.T) {
	testingT = t
	RegisterFailHandler(Fail)

	RunSpecs(t, "serve suite")
}

func testRequest(ts *httptest.Server, method, path string, requestBody interface{}, decodedRespBody interface{}) *http.Response {

	var body io.Reader
	if requestBody != nil {
		jsonData, err := json.Marshal(requestBody)
		Expect(err).NotTo(HaveOccurred())
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, ts.URL+path, body)
	Expect(err).NotTo(HaveOccurred())

	resp, err := http.DefaultClient.Do(req)
	Expect(err).NotTo(HaveOccurred())

	if decodedRespBody != nil {
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&decodedRespBody)
		Expect(err).NotTo(HaveOccurred())
	}
	defer resp.Body.Close()
	return resp
}

var _ = Describe("Todo Serve tests", func() {
	Context("When serving", Ordered, func() {
		var ts *httptest.Server
		BeforeAll(func() {
			tododb, err := sqlitedb.CreateDb()
			Expect(err).NotTo(HaveOccurred())
			todostore := store.NewSqlStore(tododb)
			todoService := todolist.NewItemsService(todostore)
			handler := &todolist.ItemsHandlers{
				ItemsService: todoService,
			}
			router := newRouter()
			handler.ConfigureRoutes(router)
			ts = httptest.NewServer(router)
		})

		AfterAll(func() {
			ts.Close()
		})

		Specify("List returns empty", func() {
			var items structs.TodoItemList
			resp := testRequest(ts, "GET", "/todolist", nil, &items)
			Expect(resp.StatusCode).To(Equal(200))
			Expect(items.Count).To(Equal(0))
		})

		Context("When todo item created", func() {
			var item structs.TodoItem
			BeforeEach(func() {
				item = structs.TodoItem{Id: "7efc0335-8da6-45f7-a9b6-d4a46ba3044b", Item: "Service motorbike",Priority: 1,Created_at: time.Now(),Updated_at: time.Now()}
				resp := testRequest(
					ts,
					"POST",
					"/todolist",
					&item,
					nil)
				Expect(resp.StatusCode).To(Equal(202))
			})

			AfterEach(func() {
				resp := testRequest(
					ts,
					"DELETE",
					"/todolist/7efc0335-8da6-45f7-a9b6-d4a46ba3044b",
					nil,
					nil)
				Expect(resp.StatusCode).To(Equal(204))
			})

			Specify("Item is returned from get", func() {
				var gItem structs.TodoItem
				resp := testRequest(
					ts,
					"GET",
					"/todolist/7efc0335-8da6-45f7-a9b6-d4a46ba3044b",
					nil,
					&gItem)
			
				Expect(resp.StatusCode).To(Equal(200))
			
				// Round the created_at and updated_at to seconds
				gItem.Created_at = gItem.Created_at.Round(time.Second)
				gItem.Updated_at = gItem.Updated_at.Round(time.Second)
				item.Created_at = item.Created_at.Round(time.Second)
				item.Updated_at = item.Updated_at.Round(time.Second)
			
				Expect(item).To(Equal(gItem))
			})

			
			Specify("Item is returned from List", func() {
				var items structs.TodoItemList
				resp := testRequest(ts, "GET", "/todolist", nil, &items)
				Expect(resp.StatusCode).To(Equal(200))
				Expect(items.Count).To(Equal(1))

				// Round the Created_at and Updated_at fields of each item in the list
				for i := range items.Items {
					items.Items[i].Created_at = items.Items[i].Created_at.Round(time.Second)
					items.Items[i].Updated_at = items.Items[i].Updated_at.Round(time.Second)
				}

				// Round the Created_at and Updated_at fields of the item to compare with the list
				item.Created_at = item.Created_at.Round(time.Second)
				item.Updated_at = item.Updated_at.Round(time.Second)

				// Now compare
				Expect(items.Items).To(ContainElement(item))
			})

			Context("When todo item modified", func() {
				var updatedItem structs.TodoItem
				BeforeEach(func() {
					updatedItem = structs.TodoItem{Id: "7efc0335-8da6-45f7-a9b6-d4a46ba3044b", Item: "Service motorbike and book MOT",Priority: 1,Created_at: time.Now(),Updated_at: time.Now()}
					resp := testRequest(ts, "PUT", "/todolist/7efc0335-8da6-45f7-a9b6-d4a46ba3044b", updatedItem, nil)
					Expect(resp.StatusCode).To(Equal(202))
				})

				
				Specify("Item is returned from get", func() {
					var gItem structs.TodoItem
					resp := testRequest(
						ts,
						"GET",
						"/todolist/7efc0335-8da6-45f7-a9b6-d4a46ba3044b",
						nil,
						&gItem)

					Expect(resp.StatusCode).To(Equal(200))

					// Round the Created_at and Updated_at fields of gItem and updatedItem
					gItem.Created_at = gItem.Created_at.UTC().Round(time.Second)
					
					updatedItem.Created_at = updatedItem.Created_at.UTC().Round(time.Second)
					

					// Now compare
					Expect(gItem.Id).To(Equal(updatedItem.Id)) // since after updation it will remain same
					Expect(gItem.Created_at).To(Equal(updatedItem.Created_at)) // since after updation it will remain same
					Expect(gItem).NotTo(Equal(item)) //since it will not remain same as updated_at will change
				})
			})

			Context("When second todo item created", func() {
				var secondItem structs.TodoItem
				BeforeEach(func() {
					secondItem = structs.TodoItem{Id: "dac2581f-9c76-47aa-877e-6c15ddcfb064", Item: "Book holiday",Priority: 1,Created_at: time.Now(),Updated_at: time.Now()}
					resp := testRequest(
						ts,
						"POST",
						"/todolist",
						&secondItem,
						nil)
					Expect(resp.StatusCode).To(Equal(202))
				})

				AfterEach(func() {
					resp := testRequest(
						ts,
						"DELETE",
						"/todolist/dac2581f-9c76-47aa-877e-6c15ddcfb064",
						nil,
						nil)
					Expect(resp.StatusCode).To(Equal(204))
				})

				Specify("Item is returned from get", func() {
					var gItem structs.TodoItem
					resp := testRequest(
						ts,
						"GET",
						"/todolist/dac2581f-9c76-47aa-877e-6c15ddcfb064",
						nil,
						&gItem)
				
					Expect(resp.StatusCode).To(Equal(200))
				
					// Convert Created_at and Updated_at to UTC and round to the nearest second for gItem and secondItem
					gItem.Created_at = gItem.Created_at.UTC().Round(time.Second)
					gItem.Updated_at = gItem.Updated_at.UTC().Round(time.Second)
					secondItem.Created_at = secondItem.Created_at.UTC().Round(time.Second)
					secondItem.Updated_at = secondItem.Updated_at.UTC().Round(time.Second)
				
					// Now compare gItem with secondItem
					Expect(secondItem).To(Equal(gItem))
				})

				Specify("Item is returned from List", func() {
					var items structs.TodoItemList
					resp := testRequest(ts, "GET", "/todolist", nil, &items)
					Expect(resp.StatusCode).To(Equal(200))
					Expect(items.Count).To(Equal(2))
				
					// Round Created_at and Updated_at to the nearest second and convert to UTC for each item in the list
					for i := range items.Items {
						items.Items[i].Created_at = items.Items[i].Created_at.UTC().Round(time.Second)
						items.Items[i].Updated_at = items.Items[i].Updated_at.UTC().Round(time.Second)
					}
				
					// Round Created_at and Updated_at to the nearest second and convert to UTC for item and secondItem
					item.Created_at = item.Created_at.UTC().Round(time.Second)
					item.Updated_at = item.Updated_at.UTC().Round(time.Second)
					secondItem.Created_at = secondItem.Created_at.UTC().Round(time.Second)
					secondItem.Updated_at = secondItem.Updated_at.UTC().Round(time.Second)
				
					// Now compare the list to ensure it contains both item and secondItem
					Expect(items.Items).To(ContainElements(item, secondItem))
				})
			})
		})
	})
})
