package main

import (
    "log"
    "os"
    "strings"
)

func main() {
    sortDesktop()
    sortDownloads()
}

func sortDesktop() {
    sortScreenshots()
    sortLua()
}

func getDesktop() (string, []os.DirEntry) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        log.Panic("UserHomeDir Error: ", err)
    }
    home_dir = home_dir + "/Desktop/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        log.Fatal("error")
    }

    return home_dir, files

}

func getDownloads() (string, []os.DirEntry) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        log.Panic("UserHomeDir Error: ", err)
    }
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
    wd, err := os.Getwd()
    if err != nil {
        log.Panic("Getwd error: ", err)
    }

    files, err := os.ReadDir(wd)
    if err != nil {
        log.Panic("Read dir error: ", err)
    }

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
        if strings.Contains(strings.ToLower(file.Name()), "screenshot") && !file.IsDir() {
            err := os.Rename(dir + file.Name(), dir + "Screenshots/" + file.Name())
            if err != nil {
                log.Panic("Rename err: ", err)
            }

            log.Println("Moved ", file, "to /screenshots/")
        }
    }

}

func sortLua() {
    dir, files := getDesktop()
    misc_dir := dir + "/Misc/"

    err := os.Chdir(misc_dir)
    if err != nil {
        log.Panic("Chdir err: ", err)
    }

    if !checkFolderExists("Lua") {
        createFolderCWD("Lua")
    }

    for _, file := range files {
        wd, err := os.Getwd()
        if err != nil {
            log.Panic("Getwd error: ", err)
        }

        filepath := wd + "/" + file.Name()

        if strings.Contains(filepath, ".lua") {
            err := os.Rename(filepath, wd+"/Misc/Lua/"+file.Name())
            if err != nil {
                log.Panic("Rename err: ", err)
            }
        }
    }

}

func contains(arr []string, item string) bool {
    for _, el := range arr {
        if el == item {
            return true
        }
    }
    return false
}

func getFileTypes(files []os.DirEntry) map[string][]string {
    fileTypes := make(map[string][]string)
    for _, file := range files {
        if !strings.HasPrefix(file.Name(), ".") && !file.IsDir() {
            // filepath := dir + file.Name()

            splitFileInfo := strings.Split(file.Name(), ".")
            fileType := splitFileInfo[len(splitFileInfo)-1]

            fileTypes[fileType] = append(fileTypes[fileType], file.Name())

        }
    }

    return fileTypes
}


func sortDownloads() {
    dir, files := getDownloads()

    err := os.Chdir(dir)
    if err != nil {
        log.Panic("Chdir err: ", err)
    }

    fileTypes := getFileTypes(files)

    manualTypes := []string{"png", "jpg", "jpeg", "html", "zip", "pdf", "csv", "docx", "app"}
    var misc []string

    for key, values := range fileTypes {
        if contains(manualTypes, key) {
            ifNotCreate(key)

            folderDir := dir + key + "/"
            for _, value := range values {
                oldPath := dir + value
                newPath := folderDir + value

                err := os.Rename(oldPath, newPath)
                if err != nil {
                    log.Panic("Rename err: ", err)
                }
            }


        } else {
            misc = append(misc, values...)
        }
    }

    ifNotCreate("Misc")

    for _, value := range misc {
        oldPath := dir + value
        newPath := dir + "Misc/" + value
        err := os.Rename(oldPath, newPath)
        if err != nil {
            log.Panic("Rename err: ", err)
        }
    }

}










