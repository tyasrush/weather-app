package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tyarus/weather-app/internal/config"
	"tyarus/weather-app/internal/infra"
	"tyarus/weather-app/pkg/utils"
)

func main() {
	cfg := config.Load()

	db, err := infra.InitDatabase(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("failed to connect MySQL: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultDBTimeout)
	defer cancel()

	// 	1. Number of orders and total amount per **status** in the last 30 days
	err = getNumberOrderAndTotalAmount(ctx, db)
	if err != nil {
		log.Fatalf("failed to get number of orders and total amount: %v", err)
	}

	fmt.Println("-----------------------------------------------------------------------")
	// 2. Top 5 customers by total spend
	err = getTopUserSpent(ctx, db)
	if err != nil {
		log.Fatalf("failed to get number of orders and total amount: %v", err)
	}
}

func getNumberOrderAndTotalAmount(ctx context.Context, db *sql.DB) error {
	query := "SELECT status, COUNT(*) as number_orders, SUM(amount) as total_spent FROM orders WHERE created_at >= NOW() - INTERVAL 30 DAY GROUP BY status"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("failed to get number of orders and total amount: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		status, numberOrders, totalAmount := "", 0, float64(0)
		if err := rows.Scan(&status, &numberOrders, &totalAmount); err != nil {
			return err
		}
		fmt.Printf("status: %s  number of orders: %d  total amount: %d\n", status, numberOrders, totalAmount)
	}

	return nil
}

func getTopUserSpent(ctx context.Context, db *sql.DB) error {
	query := "SELECT customer_id, status, COUNT(*) as number_orders, SUM(amount) as total_spent FROM orders WHERE status = 'PAID' GROUP BY customer_id ORDER BY total_spent DESC LIMIT 5"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Fatalf("failed to get number of orders and total amount: %v", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		customerID, status, numberOrders, totalAmount := "", "", 0, float64(0)
		if err := rows.Scan(&customerID, &status, &numberOrders, &totalAmount); err != nil {
			return err
		}
		fmt.Printf("customer_id: %s  status: %s  number of orders: %d  total amount: %d\n", customerID, status, numberOrders, totalAmount)
	}

	return nil
}
