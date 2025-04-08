package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Customer struct {
	id      int
	cutDone chan struct{}
}

type addRequest struct {
	customer Customer
	resp     chan bool
}

type getRequest struct {
	resp chan Customer
}

type waitingRoom struct {
	addCustomer chan addRequest
	getCustomer chan getRequest
}

func (wr *waitingRoom) run() {
	var queue []Customer
	var pendingGets []chan Customer

	for {
		select {
		case req := <-wr.addCustomer:
			if len(pendingGets) > 0 {
				pendingGets[0] <- req.customer
				pendingGets = pendingGets[1:]
				req.resp <- true
			} else if len(queue) < 6 {
				queue = append(queue, req.customer)
				req.resp <- true
			} else {
				req.resp <- false
			}

		case req := <-wr.getCustomer:
			if len(queue) > 0 {
				customer := queue[0]
				queue = queue[1:]
				req.resp <- customer
			} else {
				pendingGets = append(pendingGets, req.resp)
			}
		}
	}
}

func receptionist(wrAddChan chan<- addRequest, receptionistChan <-chan addRequest) {
	for req := range receptionistChan {
		wrResp := make(chan bool)
		wrAddChan <- addRequest{customer: req.customer, resp: wrResp}
		accepted := <-wrResp
		req.resp <- accepted
	}
}

func barber(waitingRoomGetChan chan<- getRequest) {
	for {
		req := getRequest{resp: make(chan Customer)}
		waitingRoomGetChan <- req
		customer := <-req.resp

		fmt.Printf("Barber starts cutting hair of customer %d\n", customer.id)
		cutTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
		time.Sleep(cutTime)
		fmt.Printf("Barber finishes cutting hair of customer %d\n", customer.id)

		close(customer.cutDone)
	}
}

func customer(id int, receptionistChan chan<- addRequest) {
	cutDone := make(chan struct{})
	req := addRequest{
		customer: Customer{
			id:      id,
			cutDone: cutDone,
		},
		resp: make(chan bool),
	}

	receptionistChan <- req
	accepted := <-req.resp

	if accepted {
		fmt.Printf("Customer %d is waiting for a haircut.\n", id)
		<-cutDone
		fmt.Printf("Customer %d has left the barber shop after a haircut.\n", id)
	} else {
		fmt.Printf("Customer %d left because the waiting room was full.\n", id)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	waitingRoomAdd := make(chan addRequest)
	waitingRoomGet := make(chan getRequest)
	receptionistChan := make(chan addRequest)

	wr := &waitingRoom{
		addCustomer: waitingRoomAdd,
		getCustomer: waitingRoomGet,
	}
	go wr.run()

	go receptionist(waitingRoomAdd, receptionistChan)
	go barber(waitingRoomGet)

	customerID := 1
	for {
		time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
		go customer(customerID, receptionistChan)
		customerID++
	}
}
