use std::collections::HashMap;
use std::env;
use std::process::Command;

fn main() {
    let mut args: Vec<String> = env::args().collect();
    for arg in args.iter_mut() {
        *arg = arg.trim().to_lowercase();
    }

    let command = args[1].as_str();

    if args.len() > 1 {
        match command {
            "run" => run(),
            "update" | "-u" => update(),
            "add" | "-a" => add(args),
            "remove" | "-r" => remove(args),
            "path" | "-p" => println!("{}", file!()),
            "list" | "-l" => list_crontabs(),
            "home" => println!("{}", home_path()),
            "help" | "-h" => help(),
            _ => panic!("Command not suitable: {} \n      - gfs [command] <paramater>", command),
        }
    }
}

fn help() {
    println!("Usage: gfs [command] [options]");
    let path = format!("shows you the path of where the software is stored. In this case: {}", file!());
    let commands = HashMap::from([
        ("update | -h", "updates the file sorter based on the config file"),
        ("add | -a \"[cronjob timing]\"", "add a new cronjob based off inputted timing, this will add one on top of the config cronjob. Example timing: \"0 18 * * *\" - surround timing with speech marks"),
        ("remove | -r <cronjob timing>", "Remove all gfs crontabs or remove one single existing cronjob using its timing, Example timing: \"0 18 * * *\" - surround timing with speech marks"),
        ("path | -p", path.as_str()),
        ("list | -l", "lists all of the active running crontabs"),
        ("help | -h", "provides help"),
    ]);

    println!("Commands: ");
    for command in commands.clone().keys() {
        println!("{}", format!("    {}: {}", command, commands[command]))
    }
}


fn get_exe_path() -> String {
    let binding = env::current_dir().expect("Could not get current directory");
    let cwd = binding.to_str().unwrap();
    let exe = "go-logic";

    format!("{}/{}", cwd, exe)
}

fn home_path() -> String {
    let file_dir = env::current_exe().unwrap();
    let file_dir_string = file_dir.to_str().unwrap();
    let dirs: Vec<_> = file_dir_string.split("/").collect();
    let home_dir = &dirs[0..dirs.len()-4];
    let home_dir_string = home_dir.join("/");

    home_dir_string
}

fn run() {
    let home_path = home_path();
    let exe_folder = "/go-logic";
    // let exe = "/go-logic";
    let exe_path = home_path + exe_folder;

    env::set_current_dir(exe_path.clone()).unwrap();
    println!("PATH {:?}", env::current_dir().unwrap());


    let run_exe_command = Command::new("./go-logic")
        .output()
        .unwrap();

    println!("{:?}", run_exe_command);
    let run_exe_command_output_string = String::from_utf8_lossy(&run_exe_command.stderr);
    println!("Output: \n{}", run_exe_command_output_string)
}

fn update() {
    let _cmd = Command::new("go")
        .arg("build")
        .output()
        .expect("Couldn't go build");

}

fn list_crontabs() {
    let crontabs = Command::new("crontab")
        .arg("-l")
        .output()
        .expect("Couldn't get crontabs list");

    println!("Crontabs: ");

    let crontabs_output = String::from_utf8(crontabs.stdout).unwrap();

    for crontab in crontabs_output.lines() {
        println!("{}", crontab);
    }
}

fn check_if_cronjob_exists(cronjob: String) -> bool {

    let crontabs = Command::new("crontab")
        .arg("-l")
        .output()
        .expect("Couldn't get crontabs list");

    let current_crontab = String::from_utf8(crontabs.stdout)
        .expect("Failed to parse crontab output");

    let new_crontab: Vec<&str> = current_crontab
        .lines()
        .filter(|line| !line.contains(&cronjob))
        .collect();

    if new_crontab.len() != current_crontab.lines().count() {
        true
    } else {
        false
    }
}

fn remove_all_crontabs(exe_path: String) {
    let _remove_crontabs = Command::new("bash")
        .arg("-c")
        .arg(format!(
            "crontab -l | grep -v '{}' | crontab -",
            exe_path
        ))
        .output()
        .expect("Could not get crontabs list");
}

fn remove_individual_crontab(crontab: String) {
    let _remove_crontab = Command::new("bash")
        .arg("-c")
        .arg(format!(
            "crontab -l | grep -v '{}' | crontab -",
            crontab
        ))
        .output()
        .expect("Could not get crontabs list");
}


fn remove(args: Vec<String>) {
    if args.len() == 2 {
        let exe_path = get_exe_path();
        if check_if_cronjob_exists(exe_path.clone()) {
            remove_all_crontabs(exe_path);
        }

    } else if args.len() == 3{
        let timings = args[2].clone();
        let exe_path = get_exe_path();

        let cronjob = format!("{} {}", timings, exe_path);

        if check_if_cronjob_exists(cronjob.clone()) {
            remove_individual_crontab(cronjob);
        }
    }
}

fn add(args: Vec<String>) {
    if args.len() > 3 {
       panic!("Please surround the crontab in quotes: \"0 18 * * *\"");
    }
    let exe_path = get_exe_path();

    let timings = args[2].clone();
    let cronjob = format!("{} {}", timings, exe_path);

    let _cronjob_creation = Command::new("bash")
        .arg("-c")
        .arg(format!("(crontab -l; echo '{}') | crontab -", cronjob))
        .output()
        .expect("Failed at creating cronjob");

    println!("Created cronjob: {}", cronjob);
}
