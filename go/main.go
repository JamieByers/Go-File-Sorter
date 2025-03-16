package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "net/smtp"
    "go/config"
)

func main() {
    config, err := config.HandleConfig()
    if err != nil {
        log.Fatal(err)
    }

    sortDesktop(config)
    sortDownloads(config)

    createCrontab(config.Crontab)
}

func sortDownloads(config config.Config) {
    dir, files := getDirAndFiles("Downloads")
    sortFiles(config.Downloads.FolderTypes, dir, files, config.Downloads.AllowMisc, config.Downloads.ExcludeTypes, config.Downloads.ExcludeNames)
}

func sortDesktop(config config.Config) {
    if config.Desktop.SortScreenshots {
        sortScreenshots()
    }
    dir, files := getDirAndFiles("Desktop")
    sortFiles(config.Desktop.FolderTypes, dir, files, config.Desktop.Misc.AllowMisc, config.Desktop.ExcludeTypes, config.Desktop.ExcludeNames)

    if config.Desktop.Misc.MoveAllToMisc{
        moveAllIntoMisc(dir, files, config.Desktop.Misc.ExcludeNamesFromMisc, config.Desktop.Misc.ExcludeTypesFromMisc)
    }
}

func moveAllIntoMisc(dir string, files []os.DirEntry, excluded_names []string, excluded_types []string) {
    exclude_files := []string{"Misc", ".DS_Store"}
    exclude_files = append(exclude_files, excluded_names...)

    for _, file := range files {
        if !contains(exclude_files, file.Name()) && !contains(excluded_types, file.Type().String()) {
            err := os.Rename(dir + file.Name(), dir + "Misc/" + file.Name())
            if err != nil {
                log.Fatal(err)
            }
            log.Println("Moved " + file.Name() + " to /Misc/")
        }
    }
}
func getDirAndFiles(dir string) (string, []os.DirEntry) {
    home_dir, err := os.UserHomeDir()
    if err != nil {
        emailError(err)
    }
    home_dir = home_dir + "/" + dir + "/"

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

func ifDirDoesNotExistCreate(folder string) {
    if !checkFolderExists(folder) {
        createFolderCWD(folder)
    }
}

func sortScreenshots() {
    dir, files := getDirAndFiles("Desktop")

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
    dir, files := getDirAndFiles("Desktop")
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
            splitFileInfo := strings.Split(file.Name(), ".")
            fileType := splitFileInfo[len(splitFileInfo)-1]

            // add file name to file type array
            fileTypes[fileType] = append(fileTypes[fileType], file.Name())

        }
    }

    return fileTypes
}

func sortFiles(folder_types []string, dir string, files []os.DirEntry, createMisc bool, exclude_types []string, exclude_names []string) {
    err := os.Chdir(dir)
    if err != nil {
        emailError(err)
    }

    fileTypes := getFileTypes(files)

    var file_types []string
    if contains(folder_types, "default"){
        default_types := []string{"png", "jpg", "jpeg", "html", "zip", "pdf", "csv", "docx", "app"}
        file_types = append(file_types, default_types...)
    } else if contains(folder_types, "*") {
        file_types = []string{"*"}
    } else {
        file_types = []string{"NONE"}
    }

    var misc []string


    isExcludedType := func(fileType string) bool {
        return contains(exclude_types, fileType)
    }

    isExcludedName := func(fileName string) bool {
        return contains(exclude_names, fileName)
    }

    if contains(file_types, "*") { // add all file types
        for file_type, file_names := range fileTypes {
            if !isExcludedType(file_type) {
                ifDirDoesNotExistCreate(file_type)
                moveFiles(file_names, file_type, dir, exclude_names)
            }
        }
    } else { // add only matching file types
        for file_type, file_names := range fileTypes {
            if contains(file_types, file_type) && !isExcludedType(file_type) {
                ifDirDoesNotExistCreate(file_type)
                moveFiles(file_names, file_type, dir, exclude_names)

            } else if !isExcludedType(file_type) {
                misc = append(misc, file_names...)
            }
        }

    }

    if createMisc {
        ifDirDoesNotExistCreate("Misc")

        for _, file_name := range misc {
            if !isExcludedName(file_name) {
                oldPath := dir + file_name
                newPath := dir + "Misc/" + file_name
                err := os.Rename(oldPath, newPath)
                if err != nil {
                    emailError(err)
                }
                log.Println("Moved " + oldPath + " to " + newPath)
            }
        }
    }

}

func moveFiles(file_names []string, file_type string, dir string, exclude_names []string) {
    folderDir := dir + file_type + "/"

    isExcludedName := func(fileName string) bool {
        return contains(exclude_names, fileName)
    }

    for _, file_name := range file_names {
        if !isExcludedName(file_name) {
            oldPath := dir + file_name
            newPath := folderDir + file_name

            err := os.Rename(oldPath, newPath)
            if err != nil {
                emailError(err)
            }
            log.Println("Moved: " + oldPath + " to " + newPath)
        } else {
            log.Println("Ignored: " + file_name)
        }
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

func createCrontab(crontab config.Crontab) {
    exe := "~/Documents/Code/Projects/go-file-sorter/go-file-sorter"

    crontab_line := fmt.Sprintf("%s %s %s %s %s %s",
        crontab.Minute,
        crontab.Hour,
        crontab.MonthDay,
        crontab.Month,
        crontab.DayWeek,
        exe,
    )

    log.Println("USING ", crontab_line)

    if crontabExists(crontab_line) {
        log.Println("Crontab entry already exists, skipping creation")
        return
    }

    removeExistingCronjob(exe)

    err := addCronjob(crontab_line)
    if err != nil {
        log.Fatal("Failed to add cron job:", err)
    }
}

func crontabExists(crontabLine string) bool {
    cmd := exec.Command("bash", "-c", "crontab -l")
    output, err := cmd.Output()
    if err != nil {
        return false
    }

    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        if strings.TrimSpace(line) == strings.TrimSpace(crontabLine) {
            return true
        }
    }

    return false
}

func removeExistingCronjob(exePath string) {
    cmd := exec.Command("bash", "-c", fmt.Sprintf("crontab -l | grep -v '%s' | crontab -", exePath))
    output, err := cmd.CombinedOutput()
    if err != nil {
        if strings.Contains(string(output), "no crontab") {
            log.Println("No existing crontab found")
            return
        }
        log.Printf("Error removing cron jobs: %v - %s", err, string(output))
    } else {
        log.Printf("Removed all cron jobs containing: %s", exePath)
    }
}

func addCronjob(crontab string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("(crontab -l; echo '%s') | crontab -", crontab))
	err := cmd.Run()

    log.Println("Added cronjob: ", crontab)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	return nil
}
