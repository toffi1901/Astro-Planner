import { Link } from "react-router";
import { useState } from "react";

function Login() {
  
  const [email, SetEmail] = useState('')
  const [password, SetPwd] = useState('')

  async function HandleSubmit(e) {
    e.preventDefault()
    console.log('email:', email, 'password:', password)
    console.log('test', JSON.stringify({email, password}))
    try {
      const response = await fetch('http://localhost:8080/api/auth/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({email, password})
      })
      const text = await response.text()
      console.log('response status:', response.status)
      console.log('response:', text)
    } catch (err) {
      console.log('response error:', err)
    }
  }
  return (
    <div>
        <h1>Astro Planner авторизация</h1>
        <Link to="/">На главную</Link>
        <form onSubmit={HandleSubmit}>
          <div>
            <label htmlFor="email">Эл. почта</label>
            <input type="email"
              value={email}
              onChange={(e)=>SetEmail(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="password">Пароль</label>
            <input type="password"
              value={password}
              onChange={(e)=>SetPwd(e.target.value)}
            />
          </div>
            <button type="submit">Войти</button>
          
        </form>

    </div>
  )
}

export default Login