package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// Status struct untuk menyimpan data status water dan wind
type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

// UpdateStatus secara periodik memperbarui data status
func UpdateStatus(filename string) {
	for {
		status := Status{
			Water: rand.Intn(100) + 1, // angka random antara 1-100
			Wind:  rand.Intn(100) + 1, // angka random antara 1-100
		}

		data, err := json.MarshalIndent(status, "", "    ")
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			continue
		}

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			continue
		}

		fmt.Println("Data updated:", string(data))

		time.Sleep(15 * time.Second) // tunggu 15 detik sebelum memperbarui data lagi
	}
}

// Handler untuk menangani permintaan HTTP yang mengembalikan HTML template
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("status.json")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var status Status
	err = json.Unmarshal(file, &status)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	waterStatus, windStatus := getStatusMessage(status.Water, status.Wind)

	statusData := map[string]interface{}{
		"water":        status.Water,
		"water_status": waterStatus,
		"wind":         status.Wind,
		"wind_status":  windStatus,
	}

	htmlTemplate := 
	`
	<html>
	<head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Status Water and Wind</title>
	    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f0f0f0;
        }

        .container {
            max-width: 600px;
            margin: 50px auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }

        h1 {
            text-align: center;
            margin-bottom: 20px;
        }

        .block-diagram {
            display: flex;
            justify-content: space-around;
            align-items: center;
            margin-bottom: 30px;
        }

        .block {
            width: 100px;
            height: 200px;
            text-align: center;
            background-color: #f2f2f2;
            border-radius: 8px;
            padding: 10px;
        }

        .block h2 {
            margin: 0;
        }

        .percentage {
            font-size: 24px;
            color: #333;
        }

        .status {
            font-weight: bold;
        }

        .bar {
   	    	 width: 100%;
    	    /* Tentukan tinggi blok sesuai dengan nilai maksimum */
        	height: 200px; /* Tinggi maksimum */
        	background-color: #f2f2f2;
        	border-radius: 8px;
        	padding: 10px;
        	box-sizing: border-box; /* Agar padding tidak menambah ukuran blok */
        	position: relative; /* Diperlukan untuk menambahkan lapisan di atas blok */
    	}

    	.bar.water {
        	background-color: #3498db; /* Warna biru untuk water */
    	}

    	.bar.wind {
        	background-color: #2ecc71; /* Warna hijau untuk wind */
    	}
    </style>
	</head>
	<body>
	    <div class="container">
	        <h1>Status Water and Wind</h1>
	        <div class="block-diagram">
	            <div class="block">
	                <h2>Water</h2>
					<div class="bar water" style="height: %water%;"></div>
	                <div class="percentage">%water% m</div>
	                <div class="status">%water_status%</div>
	            </div>
	            <div class="block">
	                <h2>Wind</h2>
					<div class="bar wind" style="height: %wind%;"></div>
	                <div class="percentage">%wind% m/s</div>
	                <div class="status">%wind_status%</div>
	            </div>
	        </div>
	    </div>
	</body>
	</html>
	`

	for key, value := range statusData {
		htmlTemplate = strings.ReplaceAll(htmlTemplate, "%"+key+"%", fmt.Sprint(value))
	}

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, htmlTemplate)
}


// getStatusMessage mengembalikan status berdasarkan nilai water atau wind
func getStatusMessage(water, wind int) (waterStatus, windStatus string) {
    // Tentukan status untuk air
    switch {
    case water < 5:
        waterStatus = "Aman"
    case water >= 6 && water <= 8:
        waterStatus = "Siaga"
    default:
        waterStatus = "Bahaya"
    }

    // Tentukan status untuk angin
    switch {
    case wind < 6:
        windStatus = "Aman"
    case wind >= 7 && wind <= 15:
        windStatus = "Siaga"
    default:
        windStatus = "Bahaya"
    }

    return waterStatus, windStatus
}


func main() {
	go UpdateStatus("status.json")

	http.HandleFunc("/status", StatusHandler)

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
