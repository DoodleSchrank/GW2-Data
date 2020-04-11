package models

import("database/sql"
       "entities")

type DBModel struct {
	Db *sql.DB
}

func (dbmodel DBModel) Update(item *entities.Item) {
	_, err := dbmodel.Db.Exec("update items set buyprice = ?, sellprice = ?, avgbuyprice = ?, avgsellprice = ?, lastupdate = ? where id = ?", item.Buyprice, item.Sellprice, item.Avgbuyprice, item.Avgsellprice, item.Lastupdate, item.Id)
	if err != nil {
		panic(err)
	}
}
