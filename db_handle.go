package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/sqltocsv"
	"strconv"
	"os"
	"time"
	"net/http"
	"sort"
	"strings"
	"unicode"
	//"container/list"
	"entities"
	"models"
	"config"
)


func main() {
	db, err := config.GetMySQLDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	model := models.DBModel{db}
	if err != nil {
		panic(err)
	}
	if len(os.Args) > 1 {
		if os.Args[1] == "all" {
			fmt.Println("Entering All Data Update")
			UpdateAllItemData(model)
		}
		if os.Args[1] == "dead" {
			UpdateDead(db);
			return;
		}
	}

	fmt.Println("Entering Hourly Data Update")
	UpdateHourlyItemData(model)
	WriteCSVFile(db)
}

func UpdateHourlyItemData(model models.DBModel) {
	file, err := os.Open("/var/www/go/src/hourly_ids.txt")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	avgbuyprice,	avgsellprice := 0.0, 0.0;
	// Iterate over all give item ids
	for _, idrow := range data {
		if string(idrow[0][0]) == "/" {
			fmt.Printf("At: %s\n", idrow[0])
			continue
		}
		fmt.Printf("Calculating for ID: %s\n", idrow)

		api := "http://www.gw2spidy.com/api/v0.9/csv/listings/" + idrow[0] + "/buy/1"
		prices, buyprice, length := GetLastWeeksPrices(api)
		if length == 0{
			continue
		}
		avgbuyprice = float64(prices[length/10])
		for idx, _ := range prices {
			prices[idx] = 0;
		}

		api = "http://www.gw2spidy.com/api/v0.9/csv/listings/" + idrow[0] + "/sell/1"
		prices, sellprice, length := GetLastWeeksPrices(api)
		avgsellprice = float64(prices[length - length/10])
		for idx, _ := range prices {
			prices[idx] = 0
		}
		id, _ := strconv.Atoi(idrow[0])
		item := entities.Item{id,
			"",
			buyprice,
			sellprice,
			avgbuyprice,
			avgsellprice,
			time.Now()}
	model.Update(&item)
	}
}

func UpdateAllItemData(model models.DBModel) {
	fmt.Println("Accessing Datawars2-Database")
	data, err := readCSVfromURL("https://api.datawars2.ie/gw2/v1/items/csv")
	if err != nil {
		panic(err)
	}
	fmt.Println("Begin insertion")
	for idx, row := range data {
		if idx == 0 {
			continue
		}
		id, _ := strconv.Atoi(row[0])
		name := SpaceStringsBuilder(row[1])
		buyprice, _ := strconv.ParseFloat(row[4], 64)
		sellprice, _ := strconv.ParseFloat(row[6], 64)
		buyprice = buyprice / 10000
		sellprice = sellprice / 10000
		item := entities.Item{id,
		name,
			buyprice,
			sellprice,
			0,
			0,
			time.Now()}
		model.Update(&item)
	}
	fmt.Println("Insertion complete")
	if err != nil {
		panic(err)
	}
}

func readCSVfromURL(url string) ([][]string, error) {
	resp,err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	reader.Comma = ','
	reader.LazyQuotes = true
	data, err := reader.ReadAll()
	return data, err
}
func SpaceStringsBuilder(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _,ch := range str {
		if !unicode.IsSpace(ch) && !unicode.IsPunct(ch) {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
func GetLastWeeksPrices(url string) ([]float64, float64, int) {
	mindate := time.Now().AddDate(0,0, -7).Format("2006-01-02 15:04:05 UTC")
	data, err := readCSVfromURL(url)
	var prices [300]float64
	length := 0
	lastprice := 0.0
	if err != nil {
		panic(err)
	}
	for idx, row := range data {
		if idx == 0 {
			continue
		}
		if idx == 1	{
			lastprice, _ = strconv.ParseFloat(row[1], 64)
			lastprice /= 10000
		}
		if row[0] > mindate {
			prices[idx], _ = strconv.ParseFloat(row[1], 64)
			prices[idx] /= 10000
			length = idx
		} else {
			break
		}
	}
	if length == 0 {
		return prices[0:0], 0, 0 
	}
	retprices := prices[1:length]
	sort.Float64s(retprices)
	return retprices, lastprice, length
}
func WriteCSVFile(db *sql.DB) {
	fmt.Println("Updating CSV file")
	rows, err := db.Query("select * from realItems;")
	if err != nil {
		panic(err)
	}
	err = sqltocsv.WriteFile("/var/www/go/src/data.csv", rows)
	if err != nil {
		panic(err)
	}
	fmt.Println("Updating CSV file finished")
}
func Max(a float64, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
