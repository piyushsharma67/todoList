package store

import (
	"context"
	"fmt"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	sqlitedb "go.altair.com/todolist/pkg/db"
	"go.altair.com/todolist/pkg/structs"
)

var testingT *testing.T

func TestTodoSqlStorage(t *testing.T) {
	testingT = t
	RegisterFailHandler(Fail)

	RunSpecs(t, "store suite")
}

var _ = Describe("SQL Store tests", func() {
	var tododb *sqlx.DB
	var todostore Store
	var ctx context.Context

	Context("When database created", Ordered, func() {

		BeforeAll(func() {
			var err error
			tododb, err = sqlitedb.CreateDb()
			Expect(err).NotTo(HaveOccurred())
			todostore = NewSqlStore(tododb)
			ctx = context.Background()
		})

		AfterAll(func() {
			err := tododb.Close()
			Expect(err).NotTo(HaveOccurred())
		})

		Specify("List returns empty", func() {
			var items structs.TodoItemList

			err := todostore.Update(func(tx Txn) error {
				return tx.List(ctx, &items)
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(items.Items).To(BeEmpty())
			Expect(items.Count).To(Equal(0))
		})

		Context("When todo item created", func() {
			var item structs.TodoItem
			timeNow := time.Now().UTC().Round(time.Second)
			item = structs.TodoItem{Id: "7efc0335-8da6-45f7-a9b6-d4a46ba3044b", Item: "Service motorbike",Priority: 1,Created_at: timeNow,Updated_at: timeNow}
			BeforeEach(func() {
				
				err := todostore.Update(func(tx Txn) error {
					return tx.Add(ctx, &item)
				})
				Expect(err).NotTo(HaveOccurred())
			})

			AfterEach(func() {
				// item = structs.TodoItem{Id: "7efc0335-8da6-45f7-a9b6-d4a46ba3044b", Item: "Service motorbike",Order_id: 1,Created_at: timeNow,Updated_at: timeNow}  // removing since redundant
				err := todostore.Update(func(tx Txn) error {
					return tx.Delete(ctx, item.Id)
				})
				Expect(err).NotTo(HaveOccurred())
			})

			Specify("Item is returned from get", func() {
				var gItem structs.TodoItem
				err := todostore.Update(func(tx Txn) error {
					return tx.Get(ctx, item.Id, &gItem)
				})
				
				Expect(err).NotTo(HaveOccurred())
				gItem.Created_at = gItem.Created_at.Round(time.Second)
				gItem.Updated_at = gItem.Updated_at.Round(time.Second)

				fmt.Println("****",item,gItem)
				
				Expect(gItem.Id).To(Equal(item.Id))
				Expect(gItem.Item).To(Equal(item.Item))
				Expect(gItem.Priority).To(Equal(item.Priority))
				Expect(gItem.Created_at.Equal(item.Created_at)).To(BeTrue()) // Compare the time correctly
				Expect(gItem.Updated_at.Equal(item.Updated_at)).To(BeTrue())
			})

			Specify("Item is returned from List", func() {
				var items structs.TodoItemList

				err := todostore.Update(func(tx Txn) error {
					return tx.List(ctx, &items)
				})
				Expect(err).NotTo(HaveOccurred())
				Expect(items.Count).To(Equal(1))
				for i := range items.Items {
					items.Items[i].Created_at = items.Items[i].Created_at.UTC().Round(time.Second)
					items.Items[i].Updated_at = items.Items[i].Updated_at.UTC().Round(time.Second)
				}
			
				Expect(items.Items).To(ContainElement(item))
			})

			Context("When todo item modified", func() {
				var updatedItem structs.TodoItem
				timeNow := time.Now().UTC().Round(time.Second)
				BeforeEach(func() {
					updatedItem = structs.TodoItem{
						Id:        "7efc0335-8da6-45f7-a9b6-d4a46ba3044b", 
						Item:      "Service motorbike and book MOT", 
						Priority:  1,
						Updated_at: timeNow,
						Created_at: timeNow,
					}
					err := todostore.Update(func(tx Txn) error {
						return tx.Update(ctx, &updatedItem)
					})
					Expect(err).NotTo(HaveOccurred())
				})
			
				Specify("Item is returned from get", func() {
					var gItem structs.TodoItem
					err := todostore.Update(func(tx Txn) error {
						return tx.Get(ctx, updatedItem.Id, &gItem)  // Use updatedItem.Id here
					})
					gItem.Created_at = gItem.Created_at.Round(time.Second) // rounding off the created item created_at
					
					Expect(err).NotTo(HaveOccurred())
					Expect(gItem.Id).To(Equal(item.Id))
					Expect(gItem.Item).To(Equal(updatedItem.Item))
					Expect(gItem.Priority).To(Equal(item.Priority))
					Expect(gItem.Created_at.Equal(item.Created_at)).To(BeTrue()) // Compare the time correctly
					Expect(gItem.Updated_at.Equal(item.Updated_at)).To(BeFalse()) // since we are making chnages to updated_at hence it would be falsy
					Expect(gItem).NotTo(Equal(item))      // Ensure gItem is not the same as the original item
				})
			})
			

			Context("When second todo item created", func() {
				var secondItem structs.TodoItem
				timeNow := time.Now().UTC().Round(time.Second)
				BeforeEach(func() {
					secondItem = structs.TodoItem{Id: "dac2581f-9c76-47aa-877e-6c15ddcfb064", Item: "Book holiday",Created_at: timeNow,Updated_at: timeNow,Priority: 1}
					err := todostore.Update(func(tx Txn) error {
						return tx.Add(ctx, &secondItem)
					})
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					err := todostore.Update(func(tx Txn) error {
						return tx.Delete(ctx, secondItem.Id)
					})
					Expect(err).NotTo(HaveOccurred())
				})

				Specify("Item is returned from get", func() {
					var gItem structs.TodoItem
					err := todostore.Update(func(tx Txn) error {
						return tx.Get(ctx, secondItem.Id, &gItem)
					})

					gItem.Created_at = gItem.Created_at.Round(time.Second)
					gItem.Updated_at = gItem.Updated_at.Round(time.Second)

					Expect(err).NotTo(HaveOccurred())
					// Expect(gItem).To(Equal(secondItem))

					Expect(gItem.Id).To(Equal(secondItem.Id))
					Expect(gItem.Item).To(Equal(secondItem.Item))
					Expect(gItem.Priority).To(Equal(secondItem.Priority))
					Expect(gItem.Created_at.Equal(secondItem.Created_at)).To(BeTrue()) // Compare the time correctly
					Expect(gItem.Updated_at.Equal(secondItem.Updated_at)).To(BeTrue())
				})

				Specify("Items are returned from List in ascending order of order_id and descending order of updated_at", func() {
					var items structs.TodoItemList
				
					err := todostore.Update(func(tx Txn) error {
						return tx.List(ctx, &items)
					})
					Expect(err).NotTo(HaveOccurred())
					Expect(items.Count).To(Equal(2))
				
					// Normalize time zone and round the timestamps to seconds for comparison
					for i := 1; i < len(items.Items); i++ {
						// Round timestamps to seconds for comparison 
						items.Items[i-1].Updated_at = items.Items[i-1].Updated_at.UTC().Round(time.Second)
						items.Items[i].Updated_at = items.Items[i].Updated_at.UTC().Round(time.Second)
						items.Items[i-1].Created_at = items.Items[i-1].Created_at.UTC().Round(time.Second)
						items.Items[i].Created_at = items.Items[i].Created_at.UTC().Round(time.Second)
				
						// Compare by order_id first
						Expect(items.Items[i-1].Priority).To(BeNumerically("<=", items.Items[i].Priority))
				
						// If order_id is the same, compare by updated_at (most recent updated_at should come later)
						if items.Items[i-1].Priority == items.Items[i].Priority {
							// Ensure the item with the most recent updated_at comes *last* (descending order)
							Expect(items.Items[i-1].Updated_at).To(BeTemporally(">=", items.Items[i].Updated_at))
						}
					}
				
					// Check that both items are in the list in the correct order
					Expect(items.Items).To(ContainElements(item, secondItem))
				})
			})
		})
	})
})
