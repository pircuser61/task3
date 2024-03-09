package queries

const QueryCreate = "INSERT INTO employee (Name) VALUES ($1) RETURNING empl_id"
const QueryGet = "SELECT * FROM employee WHERE empl_id = $1"
const QueryUpdate = "UPDATE employee set name=$1 WHERE empl_id = $2"
const QueryDelete = "DELETE FROM employee WHERE empl_id = $1"
const QueryList = "SELECT * FROM employee"
