package main

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
)

var allCategories []Category
var categoryByID map[int]Category
var childCategories map[int][]int

func loadCategories() {
	err := dbx.Select(&allCategories, "SELECT * FROM categories")
	if err != nil {
		log.Fatal(err)
	}

	categoryByID = make(map[int]Category, len(allCategories))
	childCategories = make(map[int][]int)
	for _, c := range allCategories {
		categoryByID[c.ID] = c
	}

	for _, c := range categoryByID {
		if c.ParentID > 0 {
			c.ParentCategoryName = categoryByID[c.ParentID].CategoryName
		}
		categoryByID[c.ID] = c

		var children []int
		for _, child := range categoryByID {
			if child.ParentID == c.ID {
				children = append(children, child.ID)
			}
		}
		if len(children) > 0 {
			childCategories[c.ID] = children
		}
	}
}

func getCategoryByID(_ sqlx.Queryer, id int) (Category, error) {
	c, ok := categoryByID[id]
	if !ok {
		return Category{}, sql.ErrNoRows
	}
	return c, nil
}

func getChildCategories(parent int) []int {
	return childCategories[parent]
}

func getAllCategories() []Category {
	return allCategories
}
