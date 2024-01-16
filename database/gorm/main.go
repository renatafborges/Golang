package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2many:products_categories;"`
}

type Product struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Price      float64
	CategoryID int
	Categories []Category `gorm:"many2many:products_categories;"`
	gorm.Model
	Category     Category
	SerialNumber SerialNumber
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductID int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	//lock pessimista
	tx := db.Begin()
	var c Category
	err = tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, 1).Error
	if err != nil {
		panic(err)
	}
	c.Name = "Eletronicos"
	tx.Debug().Save(&c)
	tx.Commit()

	//create category
	category := Category{Name: "Cozinha"}
	db.Create(&category)

	category2 := Category{Name: "Eletronicos"}
	db.Create(&category2)

	//create product
	db.Create(&Product{
		Name:       "Panela",
		Price:      99.0,
		Categories: []Category{category, category2},
	})

	//create serial number
	db.Create(&SerialNumber{
		Number:    "1020303",
		ProductID: 3,
	})

	var categories1 []Category
	if err := db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories1).Error; err != nil {
		panic(err)
	}
	for _, category := range categories1 {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			println("- ", product.Name, category.Name, "Serial Number:", product.SerialNumber.Number)
		}
	}

	//many to many
	var categories []Category
	if err := db.Model(&Category{}).Preload("Products").Find(&categories).Error; err != nil {
		panic(err)
	}
	for _, category := range categories {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			println("- ", product.Name)
		}
	}

	// //belongs to
	var products1 []Product
	db.Preload("Category").Preload("SerialNumber").Find(&products1)
	for _, product := range products1 {
		fmt.Println(product.Name, product.Category.Name, product.SerialNumber.Number)
	}

	//create batch
	products2 := []Product{
		{Name: "Tv", Price: 2500.00},
		{Name: "Mouse", Price: 50.00},
		{Name: "Keyboard", Price: 100.00},
	}
	db.Create(&products2)

	//select one
	var product Product
	db.First(&product, 1)
	fmt.Println(product)
	db.First(&product, "name = ?", "Mouse")
	fmt.Println(product)

	//select all
	var products3 []Product
	db.Find(&products3)
	for _, product := range products3 {
		fmt.Println(product)
	}

	//limit
	var products4 []Product
	db.Limit(2).Find(&products4)
	for _, product := range products4 {
		fmt.Println(product)
	}

	//offset
	var products5 []Product
	db.Limit(2).Offset(2).Find(&products5)
	for _, product := range products5 {
		fmt.Println(product)
	}

	//where
	var products6 []Product
	db.Where("price > ?", 100).Find(&products6)
	for _, product := range products6 {
		fmt.Println(product)
	}

	//like
	var products []Product
	db.Where("name LIKE  ?", "%k%").Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}

	//edit and delete
	var p Product
	db.First(&p, 1)
	p.Name = "New Mouse"
	db.Save(&p)

	var p2 Product
	db.First(&p2, 1)
	fmt.Println(p2.Name)
	db.Delete(&p2)

}
