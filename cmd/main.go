package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"passlocker/internal/locker"
	"passlocker/pkg/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Item struct {
	Name  string
	Count int
	Id    int
}

type Content struct {
	Items []locker.Element
}

func main() {
	fmt.Println("Hello World!")
	setPassFlag := flag.Bool("set", false, "Add new Password")
	getPassFlag := flag.Bool("get", false, "Get Password")
	justServe := flag.Bool("server", false, "Serve html")
	flag.Parse()
	passName := ""

	locker := locker.Locker{
		Key:      "test",
		Locked:   true,
		Elements: []locker.Element{},
	}
	locker.Connect()
	locker.Unlock()
	defer locker.Disconnect()

	if *justServe {
		e := echo.New()
		tmpl := template.New("index")
		var err error
		if tmpl, err = tmpl.Parse(view.Index); err != nil {
			fmt.Println(err)
		}

		if tmpl, err = tmpl.Parse(view.Items); err != nil {
			fmt.Println(err)
		}

		if tmpl, err = tmpl.Parse(view.Item); err != nil {
			fmt.Println(err)
		}
		e.Renderer = &TemplateRenderer{
			templates: tmpl,
		}

		items := Content{
			Items: locker.GetAllElements(),
		}
		e.Use(middleware.Logger())
		e.GET("/", func(c echo.Context) error {
			return c.Render(http.StatusOK, "index", items)
		})
		e.Logger.Fatal(e.Start(":3000"))

	} else {
		if *setPassFlag && !*getPassFlag {
			passValue := ""
			fmt.Println("We are setting Password")
			fmt.Print("Password Name: ")
			fmt.Scan(&passName)
			fmt.Print("Password Value: ")
			fmt.Scan(&passValue)
			if passName != "" && passValue != "" {
				locker.AddElement(passName, passValue)
				fmt.Printf("Added new Password %s\n", passName)
			} else {
				fmt.Println("You need to add both password name '-n' and value '-p'")
				os.Exit(1)
			}
		} else if *getPassFlag && !*setPassFlag {
			fmt.Print("Password Name: ")
			fmt.Scan(&passName)
			fmt.Println("We are getting Password")
			data := locker.GetElement(passName)
			fmt.Printf("Your Password: %s\n", data)
		} else {
			fmt.Println("Please use only one of 'get' or 'set' flags")
			os.Exit(1)
		}
	}
}
