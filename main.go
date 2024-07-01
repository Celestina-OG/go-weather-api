package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net"
        "net/http"
        "os"
        "io/ioutil"
)

type Response struct {
        ClientIP  string `json:"client_ip"`
        Location  string `json:"location"`
        Greeting  string `json:"greeting"`
}

type WeatherResponse struct {
        Location struct {
                Name string `json:"name"`
        } `json:"location"`
        Current struct {
                TempC float32 `json:"temp_c"`
        } `json:"current"`
}

func getIP(r *http.Request) string {
        ip := r.Header.Get("X-Real-IP")
        if ip == "" {
                ip = r.Header.Get("X-Forwarded-For")
        }
        if ip == "" {
                ip, _, _ = net.SplitHostPort(r.RemoteAddr)
        }
        return ip
}

func getLocation(ip string) (string, error) {
        // Use weatherapi.com to get the location based on IP
        apiKey := "9a010c30ebc249b58b5175544240107"
        url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, ip)
        resp, err := http.Get(url)
        if err != nil {
                return "", err
        }
        defer resp.Body.Close()

        var weatherResponse WeatherResponse
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return "", err
        }
        json.Unmarshal(body, &weatherResponse)

        return weatherResponse.Location.Name, nil
}

func getWeather(ip string) (float32, error) {
        // Use weatherapi.com to get the temperature based on IP
        apiKey := "9a010c30ebc249b58b5175544240107"
        url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, ip)
        resp, err := http.Get(url)
        if err != nil {
                return 0, err
        }
        defer resp.Body.Close()

        var weatherResponse WeatherResponse
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return 0, err
        }
        json.Unmarshal(body, &weatherResponse)

        return weatherResponse.Current.TempC, nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
        visitorName := r.URL.Query().Get("visitor_name")
        clientIP := getIP(r)

        location, err := getLocation(clientIP)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        temperature, err := getWeather(clientIP)
        if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
        }

        response := Response{
                ClientIP: clientIP,
                Location: location,
                Greeting: fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", visitorName, temperature, location),
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
}

func main() {
        http.HandleFunc("/api/hello", helloHandler)
        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
        }
        log.Printf("Server starting on port %s\n", port)
        log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
