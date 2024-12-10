# Find in MySQL

 Find in MySQL is a command-line tool that allows you to search for a specific keyword across all columns of a specified table in a MySQL database. The tool automatically stores or retrieves MySQL connection details in a config.ini file, making it easier to reuse.

# Features
* Easy to Use: Specify the table name and keyword as command-line parameters.
* Automatic Configuration: Prompts for MySQL connection details if not already saved and stores them in a config.ini file.
* Default Values: Offers default values for host (127.0.0.1) and port (3306) if not specified.
* Comprehensive Search: Searches all columns of the specified table for the given keyword.
* Reusable Configurations: Stores MySQL connection details for future use.

# Requirements
* Go: Version 1.19 or higher
* MySQL: A running MySQL server
* Go Modules: Install the following dependencies:
  * go-sql-driver/mysql
  * gopkg.in/ini.v1

# Installation
1 Clone the repository:
```bash
bash git clone https://github.com/brooqs/find-in-mysql.git
cd find-in-mysql
```
2. Install the required Go modules:
```bash
go mod tidy
```
3. Build project:
``` bash
go build -o findmysql
```
# Usage
Run the program with the required parameters:
```bash
./findmysql -t <table_name> -w <keyword>
```

# First Run (Prompts for Connection Details)
If the config.ini file does not exist, the program will prompt you to enter the following MySQL connection details:

* Host (default: 127.0.0.1)
* Port (default: 3306)
* Username
* Password
* Database name
These details will then be saved in the `config.ini` file for future use.
Example:
```bash
./findmysql -t users -w admin
```
```bash
config.ini not found. Please enter MySQL connection details:
Host (default: 127.0.0.1): 
Port (default: 3306): 
Username: root
Password: ****
Database name: my_database
Config saved to config.ini.
```

# Subsequent Runs (Uses `config.ini`)
For subsequent runs, the program will automatically use the saved details from `config.ini`:
 ```bash
./findmysql -t users -w admin
```
Output:
```bash
Table: users, Keyword: admin
Results (3 columns):
id: 1    username: admin    email: admin@example.com
```

# Contributing
Contributions are welcome! If you’d like to contribute, please fork the repository and submit a pull request.

# License
This project is licensed under the MIT License.
