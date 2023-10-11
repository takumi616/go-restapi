package repository 

import (
	"context"
	"log"
	"go-restapi/mysql"
	"go-restapi/entity"
)

func SelectAllDateStockIndex(ctx context.Context) ([]entity.StockIndex, error) {
	query := "SELECT * FROM nasdaq100"
	
	rows, err := mysql.Db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Failed to fetch data. %v", err)
		return nil, err
	}

	var historicalStockIndex []entity.StockIndex
	for rows.Next() {
		stockIndex := entity.StockIndex{}
		err = rows.Scan(&stockIndex.Date, &stockIndex.Open, &stockIndex.High, &stockIndex.Low, &stockIndex.Close)
		if err != nil {
			log.Printf("Failed to scan data into go struct. %v", err)
			return nil, err
		}
		historicalStockIndex = append(historicalStockIndex, stockIndex)
	}

	return historicalStockIndex, nil
}

func SelectStockIndex(ctx context.Context, date string) (entity.StockIndex, error) {
	query := "SELECT * FROM nasdaq100" + " WHERE date = ?"

	var stockIndex entity.StockIndex

	err := mysql.Db.QueryRowContext(ctx, query, date).Scan(&stockIndex.Date, &stockIndex.Open, &stockIndex.High, &stockIndex.Low, &stockIndex.Close)
	if err != nil {
		log.Printf("Failed to fetch and scan data into go struct. %v", err)
		return stockIndex, err	
	}
	
	return stockIndex, nil
}

func InsertStockIndex(ctx context.Context, stockIndex *entity.StockIndex) (int64, error) {
	query := "INSERT INTO nasdaq100" + "(date, open_price, high_price, low_price, close_price)VALUES(?, ?, ?, ?, ?)"
	result, err := mysql.Db.ExecContext(ctx, query, stockIndex.Date, stockIndex.Open, stockIndex.High, stockIndex.Low, stockIndex.Close)
	if err != nil {
		log.Printf("Failed to insert data. %v", err)
		return 0, err
	}

	var rowsNumber int64
	rowsNumber, err = result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get RowsAffected number. %v", err)
		return rowsNumber, err
	}
	
	return rowsNumber, nil
}

func UpdateStockIndex(ctx context.Context, stockIndex *entity.StockIndex, date string) (int64, error) {
	query :="UPDATE nasdaq100" + " SET date = ?, open_price = ?, high_price = ?, low_price = ?, close_price = ? WHERE date = ?"
	result, err := mysql.Db.ExecContext(ctx, query, stockIndex.Date, stockIndex.Open, stockIndex.High, stockIndex.Low, stockIndex.Close, date)
	if err != nil {
		log.Printf("Failed to update data. %v", err)
		return 0, err
	}

	var rowsNumber int64
	rowsNumber, err = result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get RowsAffected number. %v", err)
		return rowsNumber, err
	}

	return rowsNumber, nil
}

func DeleteStockIndex(ctx context.Context, date string) (int64, error) {
	query := "DELETE FROM nasdaq100" + " WHERE date = ?" 
	
	result, err := mysql.Db.ExecContext(ctx, query, date)
	if err != nil {
		log.Printf("Failed to delete data. %v", err)
		return 0, err
	}   

	var rowsNumber int64
	rowsNumber, err = result.RowsAffected()
	if err != nil {
		log.Printf("Failed to get RowsAffected number. %v", err)
		return rowsNumber, err
	}
	
	return rowsNumber, nil
}

