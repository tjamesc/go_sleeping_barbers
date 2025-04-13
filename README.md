# Sleeping Barber Problem in Go

#### Thomas Carriero - 730510525 | Minjae Kung - 730570794

**Customer Goroutine:** Each customer is created at random intervals. They send a request to the receptionist, who checks if the waiting room has space. If accepted, the customer waits until the barber closes their cutDone channel, indicating the haircut is done.

**Receptionist Goroutine:** Acts as a middleman, forwarding customer requests to the waiting room and relaying the response back.

**Waiting Room Goroutine:** Manages a FIFO queue of customers and pending barber requests. It ensures customers are added if space exists or directs them to pending barber requests immediately.

**Barber Goroutine:** Continuously processes customers from the waiting room. After cutting hair (simulated with a random sleep), the barber notifies the customer by closing their cutDone channel.
