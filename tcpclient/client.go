package tcpclient

import (
	"clientgo/config"
	"clientgo/db"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	maxConnections     = 10
	maxRequestsPerConn = 100
)

type Job struct {
	Index int
	Key   string
}

func StartClient() {
	tps := config.AppConfig.TPSConfig["s1"]
	if tps == 0 {
		log.Println("TPS is 0. Defaulting to 100 TPS.")
		tps = 100
	}

	ports := config.AppConfig.Ports
	if len(ports) == 0 {
		log.Fatal("No TCP ports defined in config.")
	}

	totalRecords := 500000
	jobChan := make(chan Job, tps)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxConnections; i++ {
		go startWorker(i, ports[i%len(ports)], jobChan, &wg)
	}

	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	currentIndex := 1
	for range ticker.C {
		if currentIndex > totalRecords {
			break
		}
		for i := 0; i < tps && currentIndex <= totalRecords; i++ {
			key := "record:" + strconv.Itoa(currentIndex)
			wg.Add(1)
			jobChan <- Job{Index: currentIndex, Key: key}
			currentIndex++
		}
	}

	// Wait for all jobs to finish
	wg.Wait()
	close(jobChan)
	log.Println("✅ All records processed.")
}

func startWorker(workerID int, port int, jobs <-chan Job, wg *sync.WaitGroup) {
	for job := range jobs {
		record, err := db.RedisGet(job.Key)
		if err != nil {
			log.Printf("[Worker %d] Redis GET error: %v", workerID, err)
			wg.Done()
			continue
		}

		conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			log.Printf("[Worker %d] Failed to connect to port %d: %v", workerID, port, err)
			wg.Done()
			continue
		}

		payload, _ := json.Marshal(record)
		_, err = conn.Write(payload)
		if err != nil {
			log.Printf("[Worker %d] TCP write failed: %v", workerID, err)
			conn.Close()
			wg.Done()
			continue
		}

		reader := make([]byte, 1024)
		n, err := conn.Read(reader)
		if err != nil {
			log.Printf("[Worker %d] TCP read failed: %v", workerID, err)
			conn.Close()
			wg.Done()
			continue
		}

		ack := string(reader[:n])
		log.Printf("[Worker %d] Received ACK: %s", workerID, ack)

		record.Status = "processed"
		record.ServerAck = ack

		log.Printf("[Worker %d] Updating Redis key %s with value: %+v", workerID, job.Key, record)

		if err := db.RedisSet(job.Key, record, time.Duration(config.AppConfig.RedisExpiration)*time.Minute); err != nil {
			log.Printf("[Worker %d] Redis SET failed: %v", workerID, err)
		} else {
			log.Printf("[Worker %d] %s updated with ack %s", workerID, job.Key, ack)
		}

		conn.Close()
		wg.Done()
	}
}
