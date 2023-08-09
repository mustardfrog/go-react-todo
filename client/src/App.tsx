import { useEffect, useState } from 'react'
import './App.css'

function App() {
  const [todos, setTodos] = useState([]);
  const [input, setInput] = useState('');
  const [change, setChange] = useState(false);
  // const [done, setDone] = useState(false);

  useEffect(() => {
    fetch("http://localhost:8080/")
      .then(res => {
        return res.json();
      })
      .then(r => {
        if (r.length < 1 || r === null) {
          setTodos([])
        }
        setTodos(r);
      })
      .catch(err => {
        console.log("Error fetching resources: " + err);
      })
  }, [change])

  const toggleTodo = async (id: any) => {
    // setDone(!done);
    setChange(!change);
        const updatedTodos = todos.map(todo => {
      if (todo.id === id) {
        return { ...todo, done: !todo.done };
      }
      return todo;
    });

    setTodos(updatedTodos);


    await fetch(`http://localhost:8080/${id}`, {
      "method": "PUT",
      "headers": { "Content-Type": "application/json" },
    })
      .then(res => {
        return res.json()
      })
      .then(json => {
        console.log(json);
      })
      .catch(err => console.log(err))

  }

  const deleteTodo = async (id: any) => {
    await fetch(`http://localhost:8080/${id}`, {
      method: "DELETE",
    })

    setChange(!change);
  }

  const handleSubmit = async (e: any) => {
    e.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/', {
        method: "POST",
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: `${input}`, done: false })
      })

      if (response.ok) {
        console.log("response was ok");
      }

    } catch (err) {
      console.log("Post request seems to have failed");
    }

    setChange(!change);
    setInput('');
  }

  function handleChange(e: any) {
    setInput(e.target.value);
  }

  return (
    <>
      <form onSubmit={handleSubmit}>
        <input type='text' value={input} onChange={handleChange} />
        <button type='submit'>Add</button>
      </form>

      {todos.length < 1 ? (
        <p>No data</p>
      ) : (
        <div>
          {todos.map(todo => (
            <li key={todo.id}>
              {/* <span onClick={() => toggleTodo(todo.id)} className={done ? "todo done" : "todo notdone"}>{todo.title}</span> */}
              <span onClick={() => toggleTodo(todo.id)} style={{ textDecoration: todo.done ? "line-through" : "none", cursor: "pointer"}}>{todo.title}</span>
              <button onClick={() => deleteTodo(todo.id)}>Delete</button>
            </li>
          ))}
        </div>
      )}
    </>
  )
}

export default App
