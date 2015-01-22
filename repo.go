package main

import (
    "fmt"
    "log"
    "time"
    "github.com/gocql/gocql"

)

var todos Todos

var (
    cluster *gocql.ClusterConfig
)
/*
create keyspace todo with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
create table todo.todos(id UUID, name text, completed boolean, due timestamp, PRIMARY KEY(id));
create index on todo.todos(????);
*/

var (
    insertQuery = "INSERT INTO todos(id, name, completed, due) VALUES(?, ?, ?, ?);"
    getAllQuery = "SELECT id, name, completed, due FROM todos;"
    getQuery = "SELECT id, name, completed, due FROM todos WHERE id = ?;"
)

// Give us some seed data
func init() {
    cluster = gocql.NewCluster("10.254.254.100")
    cluster.Keyspace = "todo"
    cluster.Consistency = gocql.Quorum
    RepoCreateTodo(Todo{Name: "todo 1"})
    RepoCreateTodo(Todo{Name: "todo 2"})
    RepoGetAll()
}

func RepoGetAll() {
    c, e := GetDataAccessSession()
    if e != nil {
        return
    }
    var tid, name string 
    var completed bool
    var due time.Time
    iter := c.Query(getAllQuery).Iter()
    for iter.Scan(&tid, &name, &completed, &due) {
        fmt.Println(&tid)
        todos = append(todos, Todo{
                 Id: tid,
                Name: name,
                Completed: completed,
                Due: due,
            })
    }
    if err := iter.Close(); err != nil {
        log.Fatal(err)
    }
}

func RepoFindTodo(id int) (Todo, error) {
    c, e := GetDataAccessSession()
    if e != nil {
        return Todo{}, e
    }
    var tid, name string 
    var completed bool
    var due time.Time
    e = c.Query(getQuery, id).Scan(&tid, &name, &completed, &due)
    if e != nil {
        return Todo{}, e
    }
    return Todo{
        Id: tid,
        Name: name,
        Completed: completed,
        Due: due,
    }, nil
}

func RepoCreateTodo(t Todo) (string, error) {
    c, e := GetDataAccessSession()
    if e != nil {
        return "", e
    }
    id := t.Id
    if len(id) == 0 {
        uuid, e := gocql.RandomUUID()
        if e != nil {
            fmt.Println("ID Error")
            return "", e
        }
        id = uuid.String()
    }
    due := t.Due
    if due.IsZero() {
        due = time.Now()
    }

    e = c.Query(insertQuery, id, t.Name, t.Completed, due).Exec()
    if e != nil {
        return "", e
    }
    return id, nil
}

// func RepoDestroyTodo(id int) error {
//     for i, t := range todos {
//         if t.Id == id {
//             todos = append(todos[:i], todos[i+1:]...)
//             return nil
//         }
//     }
//     return fmt.Errorf("Could not find Todo with id of %d to delete", id)
// }

func GetDataAccessSession()  (*gocql.Session, error)  {
    session, e := cluster.CreateSession()
    // defer session.Close()
    return session, e
}

func GetData() {
    // connect to the cluster
    cluster := gocql.NewCluster("cassandra")
    cluster.Keyspace = "example"
    cluster.Consistency = gocql.Quorum
    session, _ := cluster.CreateSession()
    defer session.Close()

    // insert a tweet
    if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`,
        "me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
        log.Fatal(err)
    }

    var id gocql.UUID
    var text string

    /* Search for a specific set of records whose 'timeline' column matches
     * the value 'me'. The secondary index that we created earlier will be
     * used for optimizing the search */
    if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`,
        "me").Consistency(gocql.One).Scan(&id, &text); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Tweet:", id, text)

    // list all tweets
    iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
    for iter.Scan(&id, &text) {
        fmt.Println("Tweet:", id, text)
    }
    if err := iter.Close(); err != nil {
        log.Fatal(err)
    }
}