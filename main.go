package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "math"
)
	
var colorDict = map[string]string{
    "White":  "#FFFFFF",
    "Black":  "#000000",
    "Red":    "#FF0000",
    "Green":  "#00FF00",
    "Blue":   "#0000FF",
    "Yellow": "#FFFF00",
}

// hex в RGB
func hexToRGB(hex string) (int, int, int) {
    r, _ := strconv.ParseInt(hex[1:3], 16, 8)
    g, _ := strconv.ParseInt(hex[3:5], 16, 8)
    b, _ := strconv.ParseInt(hex[5:7], 16, 8)
    return int(r), int(g), int(b)
}

// RGB в hex
func rgbToHex(r, g, b int) string {
    return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func closestColorName(r, g, b int) string {
    minDist := math.MaxFloat64
    closestName := "Unknown"

    for name, hex := range colorDict {
        r1, g1, b1 := hexToRGB(hex)
        dist := math.Sqrt(float64((r-r1)*(r-r1) + (g-g1)*(g-g1) + (b-b1)*(b-b1)))
        if dist < minDist {
            minDist = dist
            closestName = name
        }
    }

    return closestName
}

func combineHexColors(colors []string) (string, string) {
    var rSum, gSum, bSum int
    count := len(colors)

    for _, color := range colors {
        r, g, b := hexToRGB(color)
        rSum += r
        gSum += g
        bSum += b
    }

    rCombined := rSum / count
    gCombined := gSum / count
    bCombined := bSum / count

    combinedColor := rgbToHex(rCombined, gCombined, bCombined)

    colorName := closestColorName(rCombined, gCombined, bCombined)

    return combinedColor, colorName
}

func combineColors(w http.ResponseWriter, r *http.Request) {
    var colorReq struct {
        Colors []string `json:"colors"`
    }

    err := json.NewDecoder(r.Body).Decode(&colorReq)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    if len(colorReq.Colors) == 0 {
        http.Error(w, "No colors provided", http.StatusBadRequest)
        return
    }

    for _, color := range colorReq.Colors {
        if len(color) != 7 || color[0] != '#' {
            http.Error(w, "Invalid color format", http.StatusBadRequest)
            return
        }
        _, err := strconv.ParseInt(color[1:], 16, 64)
        if err != nil {
            http.Error(w, "Invalid hex color code", http.StatusBadRequest)
            return
        }
    }

    combinedColor, colorName := combineHexColors(colorReq.Colors)

    json.NewEncoder(w).Encode(map[string]string{"combinedColor": combinedColor, "colorName": colorName})
}

func main() {
    http.HandleFunc("/colors", combineColors)
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}
