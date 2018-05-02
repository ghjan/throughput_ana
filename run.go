package main

import (
"flag"
"fmt"
)

func main() {
	total := flag.Int("total", 100, "how many rows created")
	filePath := flag.String("filePath", "/application/nginx/logs/dig.log", "dig log file path")
	flag.Parse()

	fmt.Println(*total)
	fmt.Println(*filePath)

}
