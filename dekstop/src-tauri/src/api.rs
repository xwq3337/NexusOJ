use rusqlite::{params, Connection, Result};

pub struct Database  {
    pub conn: Connection,
}
impl Database {
    pub fn new(conn: Connection) -> Self {
        Self { conn }
    }
    pub fn create_table(conn: &Connection) -> Result<()> {
        conn.execute(
            "CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY,
                name TEXT NOT NULL,
                email TEXT NOT NULL UNIQUE
            )",
            params![],
        )?;
        Ok(())
    }
}

