package main

import (
    "log"
    "os"
    "strings"
)

func main() {
    // sortScreenshots()
    sortLua()
}

func baseDesktop() (string, []os.DirEntry) {
    home_dir, _ := os.UserHomeDir()
    home_dir = home_dir + "/Desktop/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        log.Fatal("error")
    }

    return home_dir, files

}

func sortScreenshots() {
    dir, files := baseDesktop()

    for _, file := range files {
        if strings.Contains(strings.ToLower(file.Name()), "screenshot") {
            os.Rename(dir + file.Name(), dir + "Screenshots/" + file.Name())
            log.Println("Moved ", file, "to /screenshots/")
        }
    }

}

func sortLua() {
    dir, files := baseDesktop()

    os.Chdir(dir)

    for _, file := range files {
        wd, _ := os.Getwd()
        filepath := wd + "/" + file.Name()
        log.Println(filepath)
        log.Println(strings.Contains(filepath, ".lua"))
    }

    // for _, file := range files {
    //     if strings.Contains(strings.ToLower(file.Type()), "screenshot") {
    //         os.Rename(dir + file.Name(), dir + "Screenshots/" + file.Name())
    //         log.Println("Moved ", file, "to /screenshots/")
    //     }
    // }
}
