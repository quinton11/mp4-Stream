package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"mp4stream/internal/server"

	"github.com/gorilla/mux"

	"github.com/manifoldco/promptui"
)

type Messages struct {
	Name    string
	Details string
}

func main() {

	DirOrAbs()
}

// Prompt user to choose whether to provide absolute path of movie
// or read from directory, default directory is assets in root folder
// if assets is not found it prompts user to create asset
func DirOrAbs() {
	messages := []Messages{
		{"Movie Path", "Provide the absolute path to movie"},
		{"Use Assets", "Use assets folder in root directory"},
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002662 {{ .Name | green }}  ({{ .Details | red }})",
		Inactive: "{{ .Name | yellow }}",
		Selected: "\U00002705 {{ .Name | blue | cyan }}",
	}

	prompt := promptui.Select{
		Label:        "Mode",
		Items:        messages,
		Templates:    templates,
		Size:         2,
		HideSelected: true,
	}

	sel, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	if sel == 0 {
		//Selection 1
		EnterMovie()
	}
	if sel == 1 {
		path, err := PathCheck()
		if err != nil {
			panic(err)
		}
		movies, err := ExistingMovies(path)
		if err != nil {
			panic(err)
		}
		//fmt.Println(movies)
		if len(movies) == 0 {
			fmt.Println("Add movie to assets ")

		} else {
			mov := MoviesSelection(movies)

			//verify if path is correct
			if _, err := os.Stat(path + mov); err != nil {
				fmt.Println(err)
				panic(err)
			}
			fmt.Println("In else")
			ph := path + mov
			Run(ph)
		}
	}
}

func Run(path string) {
	router := mux.NewRouter()
	server := server.NewServer(router, path)

	server.Listen()
}

// This runs after Assets Dir is selected and list of .mp4 and .avi movies are displayed
// user selects which of them to show
func MoviesSelection(movies []string) string {
	var messages []Messages

	for _, movie := range movies {
		messages = append(messages, Messages{movie, ""})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002662 {{ .Name | green }}",
		Inactive: "{{ .Name | yellow }}",
		Selected: "\U00002705 {{ .Name | blue | cyan }}",
	}

	prompt := promptui.Select{
		Label:        "Mode",
		Items:        messages,
		Templates:    templates,
		Size:         2,
		HideSelected: true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		panic(err)
	}
	return movies[i]
}

// Check if assets directory exists and if file exists
func PathCheck() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filename := "\\assets\\"
	//check if path exists
	if _, err := os.Stat(dir + filename); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(dir+filename, os.ModeDir); err != nil {
				return "", err
			}
			return "", nil
		}
		return "", err
	}
	return dir + filename, nil
}

// returns list of movies in directory
func ExistingMovies(path string) ([]string, error) {
	var movies []string

	dir, err := os.Open(path)
	if err != nil {
		return movies, err
	}

	fileNames, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		return movies, err
	}

	for _, file := range fileNames {
		//Check if file ends in .avi or .mp4
		if strings.HasSuffix(file, ".mp4") || strings.HasSuffix(file, ".avi") {
			movies = append(movies, file)
		}
	}
	return movies, nil
}

func EnterMovie() {
	validate := func(input string) error {
		_, err := os.Stat(input)
		if err != nil {
			//fmt.Println("")
			return errors.New("invalid path")
		}
		if !(strings.HasSuffix(input, ".mp4") || strings.HasSuffix(input, ".avi")) {
			return errors.New("invalid path")
		}
		return nil
	}
	template := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}

	movieprompt := promptui.Prompt{
		Label:       "Absolute Path: ",
		HideEntered: true,
		Templates:   template,
		Validate:    validate,
	}

	movie, err := movieprompt.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(movie)
	Run(movie)

	//Check if file exists
}
