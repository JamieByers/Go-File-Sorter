package main

import (
    "log"
    "os"
    "strings"
    "net/smtp"
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
        emailError(err)
    }
    home_dir = home_dir + "/Desktop/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        emailError(err)
    }

    return home_dir, files

}

func getDownloads() (string, []os.DirEntry) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        emailError(err)
    }
    home_dir = home_dir + "/Downloads/"

    files, err := os.ReadDir(home_dir)
    if err != nil {
        emailError(err)
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
        emailError(err)
    }

    return createFolder(folder, wd)
}


func checkFolderExists(folder string) bool {
    wd, err := os.Getwd()
    if err != nil {
        emailError(err)
    }

    files, err := os.ReadDir(wd)
    if err != nil {
        emailError(err)
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
                emailError(err)
            }
            log.Println("Moved " + file.Name() + " to /Screenshots/")
        }
    }

}

func sortLua() {
    dir, files := getDesktop()
    misc_dir := dir + "/Misc/"

    err := os.Chdir(misc_dir)
    if err != nil {
        emailError(err)
    }

    if !checkFolderExists("Lua") {
        createFolderCWD("Lua")
    }

    for _, file := range files {
        wd, err := os.Getwd()
        if err != nil {
            emailError(err)
        }

        filepath := wd + "/" + file.Name()

        if strings.Contains(filepath, ".lua") {
            err := os.Rename(filepath, wd+"/Misc/Lua/"+file.Name())
            if err != nil {
                emailError(err)
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
        emailError(err)
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
                    emailError(err)
                }
                log.Println("Moved " + oldPath + " to " + newPath)
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
            emailError(err)
        }
        log.Println("Moved " + oldPath + " to " + newPath)
    }

}


func emailError(err error) {
    sender := ""
    password := ""
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"
    recipient := ""

    subject := "There is an error with your Go Macos Sorter"
    body := "The error occuring: " + err.Error()

    message := []byte("Subject: " + subject + "\r\n" +
    "From: " + sender + "\r\n" +
    "To: " + recipient + "\r\n" +
    "\r\n" +
    body + "\r\n")


    auth := smtp.PlainAuth("", sender, password, smtpHost)

    e := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, []string{recipient}, message)

	if e != nil {
        log.Panic(e)
	}

	log.Println("Email sent successfully!")
}

