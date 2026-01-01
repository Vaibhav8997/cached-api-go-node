const express = require("express");  //express tool to create APIs in NodeJS

const { Pool } = require("pg");  //pg library to intrect with postgres database

const Redis = require("ioredis"); //redis library

const app = express();
app.use(express.json());   //to interect with server using json 

const pool = new Pool({
    host: "postgres",
    user: "admin",
    password: "123456",
    database: "testDb",
    port: 5432
});

const redis = new Redis({
    host:"redis", port:6379
});

//POST method-- data insertion
app.post("/data", async (req, res) => {
    const { id, name, email } = req.body;

    try {
        // converted data into json string and inserted in user_queue
        const task = JSON.stringify({ id, name, email });
        await redis.lpush("user_queue", task); 
        
        res.status(201).json({ message: "Data received and queued!" });
    } catch(err) {
        res.status(500).json({ error: "Redis Error" });
    }
});

//GET method-- data retrition

app.get("/data/:id", async (req, res) => {
    try {
        const result = await pool.query(
            "SELECT * FROM test_data where id = $1",
            [req.params.id]
        );

        if (!result.rows.length) {
            return res.status(404).json({ error: "Data notfound" });
        }

        res.json(result.rows[0]);
    } catch {
        res.status(500).json({ error: "DB Error" });
    }
});

//Get all data
app.get("/data", async (req, res) => {
  const result = await pool.query(
    "SELECT * FROM test_data ORDER BY created_at DESC"
  );
  res.json(result.rows);
});


//server listen on 3000 port
app.listen(3000, () => {
  console.log("Node API running on port 3000");
});