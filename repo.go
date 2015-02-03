package main

import (
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"time"
)

var todos []Todo

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
	deleteQuery = "DELETE FROM todos WHERE id = ?;"
	getAllQuery = "SELECT id, name, completed, due FROM todos;"
	getQuery    = "SELECT id, name, completed, due FROM todos WHERE id = ?;"
)

func init() {
	cluster = gocql.NewCluster("10.254.254.100")
	cluster.Keyspace = "todo"
	cluster.Consistency = gocql.Quorum
}

func RepoGetAll() {
	c, e := cluster.CreateSession()
	if e != nil {
		return
	}
	var tid, name string
	var completed bool
	var due time.Time
	todos = make([]Todo, 0)
	iter := c.Query(getAllQuery).Iter()
	for iter.Scan(&tid, &name, &completed, &due) {
		todos = append(todos, Todo{
			Id:        tid,
			Name:      name,
			Completed: completed,
			Due:       due,
		})
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}
}

func RepoFindTodo(id string) (Todo, error) {
	c, e := cluster.CreateSession()
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
		Id:        tid,
		Name:      name,
		Completed: completed,
		Due:       due,
	}, nil
}

func RepoCreateTodo(t Todo) (Todo, error) {
	c, e := cluster.CreateSession()
	if e != nil {
		return Todo{}, e
	}
	id := t.Id
	if len(id) == 0 {
		uuid, e := gocql.RandomUUID()
		if e != nil {
			fmt.Println("ID Error")
			return Todo{}, e
		}
		id = uuid.String()
	}
	due := t.Due
	if due.IsZero() {
		due = time.Now()
	}

	e = c.Query(insertQuery, id, t.Name, t.Completed, due).Exec()
	if e != nil {
		return Todo{}, e
	}

	t.Id = id
	t.Due = due
	return t, nil
}

func RepoDestroyTodo(id string) error {
	c, e := cluster.CreateSession()
	e = c.Query(deleteQuery, id).Exec()
	if e != nil {
		return fmt.Errorf("Could not find Todo with id of %d to delete", id)
	}
	return nil
}
