package handle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"../model"
)

type OrderReq struct {
	Symbol   string `json:"symbol"`
	Side     string `json:"side"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type OrderResp struct {
	OrderID string `json:"order_id"`
}

type BookReq struct {
	Symbol string `json:"symbol"`
}

type BookItem struct {
	Depth    int    `json:"depth"`
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
}

type BookResp struct {
	Sell []BookItem `json:"sell"`
	Buy  []BookItem `json:"buy"`
}

func Order(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		w.Write([]byte("please use POST method"))
		return
	}
	db := model.GetDB()
	if db == nil {
		log.Println("connect to mysql error")
		w.WriteHeader(404)
		w.Write([]byte("connect to mysql error"))
		return
	}
	defer db.Close()
	//取出post的body
	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	log.Println("body:", body)

	//解析body
	orderReq, err := getOrderReq(body)
	if err != nil {
		log.Println("get order parameters failed")
		w.WriteHeader(402)
		w.Write([]byte("get order parameters failed"))
		return
	}
	orderId, err := getOrderResp(orderReq, db)
	if err != nil {
		log.Println("get order resp failed:", err)
		w.WriteHeader(402)
		w.Write([]byte("get order resp failed"))
		return
	}

	w.WriteHeader(200)
	rsp := fmt.Sprintf(`{"order_id":%s}`, strconv.Itoa(int(orderId)))
	w.Write([]byte(rsp))
	return
}

func Book(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(404)
		w.Write([]byte("please use POST method"))
		return
	}
	db := model.GetDB()
	if db == nil {
		log.Println("connect to mysql error")
		w.WriteHeader(404)
		w.Write([]byte("connect to mysql error"))
		return
	}
	defer db.Close()

	body, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	log.Println("body:", body)

	//解析body
	bookReq, err := getBookReq(body)
	if err != nil {
		log.Println("get book parameters failed")
		w.WriteHeader(402)
		w.Write([]byte("get book parameters failed"))
		return
	}

	//获取book的数据
	rsp, err := getBookResp(bookReq, db)
	if err != nil {
		log.Println("getBookResp error")
		w.WriteHeader(404)
		w.Write([]byte("getBookResp error"))
		return
	}
	w.WriteHeader(200)
	data, _ := json.Marshal(rsp)
	w.Write(data)
	return
}

func getOrderReq(data []byte) (OrderReq, error) {
	orderReq := OrderReq{}
	err := json.Unmarshal(data, &orderReq)
	if err != nil {
		return orderReq, err
	}
	return orderReq, nil
}

func getOrderResp(orderReq OrderReq, db *sql.DB) (int64, error) {
	quantity, _ := strconv.Atoi(orderReq.Quantity)
	log.Println("quantity:", quantity)
	//上传order消息到数据库
	res, err := db.Exec(`insert into test.order(symbol,side,price,quantity) values(?,?,?,?)`, orderReq.Symbol, orderReq.Side, orderReq.Price, quantity)
	if err != nil {
		log.Println("Insert failed,err:", err)
		return 0, err
	}
	orderId, err := res.LastInsertId()
	return orderId, err
}

func getBookReq(data []byte) (BookReq, error) {
	bookReq := BookReq{}
	err := json.Unmarshal(data, &bookReq)
	if err != nil {
		return bookReq, err
	}
	return bookReq, nil
}

func getBookResp(bookReq BookReq, db *sql.DB) (BookResp, error) {

	bookResp := BookResp{}
	sides := []string{"sell", "buy"}

	for _, side := range sides {
		var rows *sql.Rows
		books := make([]BookItem, 0)
		rows, err := db.Query(`select count(id) as depth,price,sum(quantity) from test.order where symbol=? and side=? group by price`, bookReq.Symbol, side)
		if err != nil {
			fmt.Println(err)
			return bookResp, err
		}

		for rows.Next() {
			var depth, quantity int
			var price string

			rows.Scan(&depth, &price, &quantity)
			book := BookItem{
				Depth:    depth - 1,
				Price:    price,
				Quantity: strconv.Itoa(quantity),
			}
			books = append(books, book)
		}
		if side == "sell" {
			bookResp.Sell = books
		} else if side == "buy" {
			bookResp.Buy = books
		}
		rows.Close()
	}

	return bookResp, nil

}
