# GoDBSniffer

GoDBSniffer is a CLI-based application developed in Go, designed to connect to a MySQL database and do some sniffing, which are explained in the next step in the features.

## Features

- **Database Connection Setup**: Configure and establish a connection to MySQL databases.
- **Health Checks**: 

| Check                 | Description                                         | Purpose                                       |
|-----------------------|-----------------------------------------------------|-----------------------------------------------|
| Uptime                | Measures server uptime since last restart.          | Identifies potential uptime issues.           |
| Active Connections    | Counts current active database connections.         | Prevents exceeding max connections limit.     |
| Open Tables           | Monitors open tables against a configured limit.    | Helps in adjusting `table_open_cache` setting.|
| Buffer Pool Usage     | Analyzes InnoDB buffer pool usage.                  | Ensures buffer pool is sized correctly.       |


- **Security Audits**:

| Check            | Description                                                | Purpose                                               |
|------------------|------------------------------------------------------------|-------------------------------------------------------|
| User Privileges  | Checks for users with excessive administrative privileges. | Reduces risk of unauthorized access.                  |
| Empty Passwords  | Identifies accounts with empty passwords.                  | Mitigates the risk of easy unauthorized access.       |
| Version Check    | Ensures database is running a supported version.           | Protects against vulnerabilities in older versions.   |

- **Performance Analysis**: 

| Check             | Description                                        | Purpose                                                                      |
|-------------------|----------------------------------------------------|---------------------------------------------                                 |
| Query Performance | Assesses how effectively queries are being executed, especially their use of indexes. | Optimizes query execution and index usage 
| Slow Queries      | Identifies and counts queries that execute slowly. | Helps in optimizing slow queries to improve overall performance.             |


- **Schema Review**: Examine database schemas to ensure best practices are followed.

## Getting Started

### Prerequisites

- Go (version 1.15 or higher) [Download Go](https://golang.org/dl/) 
- MySQL 5.7 or higher
- Access to a MySQL database
- A working database with its details (host: "localhost"
  port: 3306
  user: "root"   
  pass: "testpass"  
  name: "SampleDB"
)

### Installation

Clone the repository to your local machine:

```bash
git clone ....
cd GoDBSniffer

go build


Runing the Application

1. Navigate to the directory: cd path\to\GoDBSniffer
2. Execute the application: .\GoDBSniffer.exe
3. Enter these details:
    Database Host: Enter the hostname of your MySQL database (e.g., localhost).
    Database Port: Specify the port on which your MySQL database is running (e.g., 3306).
    Database User: Provide the username to access the database.
    Database Password: Enter the password associated with the username.
    Database Name: Type the name of the database you want to analyze.

