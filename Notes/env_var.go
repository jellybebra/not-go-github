package notmain

import (
	"fmt"
	"os"

	// "strings"
	"github.com/joho/godotenv"
)

func main() {
	// foo := os.Environ()
	// fmt.Println(strings.Join(foo, "\n"))

	// f := os.Getenv("KEY")
	// fmt.Println(f)

	dotenv := goDotEnvVariable("KEY")
	fmt.Printf("godotenv : %s = %s \n", "KEY", dotenv)

}

func goDotEnvVariable(key string) string {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	return os.Getenv(key)
}
