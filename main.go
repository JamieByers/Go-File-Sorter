package main

import (
    "log"
    "os"
    "strings"
)

func main() {
    sortDownloads()
}

func sortDesktop() {
    sortScreenshots()
    sortLua()
}

func getDesktop() (string, []os.DirEntry) {
    home_dir, _ := os.UserHomeDir()
    home_dir = home_dir + "/Desktop/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        log.Fatal("error")
    }

    return home_dir, files

}

func getDownloads() (string, []os.DirEntry) {
    home_dir, _ := os.UserHomeDir()
    home_dir = home_dir + "/Downloads/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        log.Fatal("error")
    }

    return home_dir, files

}

func createFolder(folder string, dir string) bool {
    os.Chdir(dir)
    if !checkFolderExists(folder) {
        os.Mkdir(folder, 0777)
        return true
    }

    return false
}

func createFolderCWD(folder string) bool {
    wd, err := os.Getwd()
    if err != nil {
        log.Panic(err)
    }

    return createFolder(folder, wd)
}


func checkFolderExists(folder string) bool {
    wd, _ := os.Getwd()
    files, _ := os.ReadDir(wd)

    for _, file := range files {
        if strings.ToLower(file.Name()) == folder && file.IsDir() {
            return true
        }
    }

    return false

}

func ifNotCreate(folder string) {
    if !checkFolderExists(folder) {
        createFolderCWD(folder)
    }
}

func sortScreenshots() {
    dir, files := getDesktop()

    for _, file := range files {
        if strings.Contains(strings.ToLower(file.Name()), "screenshot") {
            os.Rename(dir + file.Name(), dir + "Screenshots/" + file.Name())
            log.Println("Moved ", file, "to /screenshots/")
        }
    }

}

func sortLua() {
    dir, files := getDesktop()
    misc_dir := dir + "/Misc/"

    os.Chdir(misc_dir)

    if !checkFolderExists("Lua") {
        createFolderCWD("Lua")
    }

    for _, file := range files {
        wd, _ := os.Getwd()
        filepath := wd + "/" + file.Name()

        if strings.Contains(filepath, ".lua") {
            os.Rename(filepath, wd+"/Misc/Lua/"+file.Name())
            log.Println(filepath)
            log.Println(strings.Contains(filepath, ".lua"))
        }
    }

}

func fileTypesHasType(arr []string, item string) bool {
    for _, el := range arr {
        if el == item {
            return true
        }
    }

    return false
}

func sortDownloads() {
    dir, files := getDownloads()

    os.Chdir(dir)

    var filepaths []string
    var fileTypes []string

    for _, file := range files {
        if !strings.HasPrefix(file.Name(), ".") {
            filepath := dir + file.Name()
            filepaths = append(filepaths, filepath)

            splitFileInfo := strings.Split(file.Name(), ".")
            fileType := splitFileInfo[len(splitFileInfo)-1]

            if !fileTypesHasType(fileTypes, fileType) && !file.IsDir() {
                fileTypes = append(fileTypes, fileType)
            }

            log.Println(file.Name(), fileType)
            log.Println("filetypes: ", fileTypes)
        }
    }

    log.Println(filepaths)

}










