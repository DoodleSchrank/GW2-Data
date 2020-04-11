package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)


func UpdateDead(db *sql.DB) {
	results, err := db.Query("select id from items;")
	if err != nil {
		panic(err)
	}
	defer results.Close()
	id := 0
	for results.Next() {
		results.Scan(&id)
		fmt.Printf("Scanning id: %d for sales\n", id)
		apilink := fmt.Sprintf("%s%d%s", "http://www.gw2spidy.com/api/v0.9/csv/listings/", id, "/buy/1")

		changed := GetDeadItem(apilink)
		db.Exec("update items set deaditem = ? where id = ?", changed, id)
	}
}
func GetDeadItem(url string) (changed bool) {
	mindate := time.Now().AddDate(0,0, -7).Format("2006-01-02 15:04:05 UTC")
	listings := 0
	quantity := 0
	price := 0
	changed = false

	data, err := readCSVfromURL(url)
	if err != nil {
		panic(err)
	}
	for idx, row := range data {
		if idx == 0 {
			continue
		}
		if idx == 1	{
			price, _ = strconv.Atoi(row[1])
			quantity, _ = strconv.Atoi(row[2])
			listings, _ = strconv.Atoi(row[3])
		}
		currPrice, _ := strconv.Atoi(row[1])
		currQuantity, _ := strconv.Atoi(row[2])
		currListings, _ := strconv.Atoi(row[3])
		if row[0] > mindate {
			if currListings == listings && currQuantity == quantity  && currPrice == price{
				continue
			} else {
				changed = true
				return
			}
		} else {
			break;
		}
	}
	return
}
