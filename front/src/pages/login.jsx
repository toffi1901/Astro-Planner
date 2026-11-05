import { Link, useNavigate } from "react-router";
import { useState } from "react";
import { useAuth } from "../auth.jsx";

function Login() {
  let navigate = useNavigate()
  const { login } = useAuth()

  const [email, SetEmail] = useState('')
  const [password, SetPwd] = useState('')
  const [error, setError] = useState('')

  async function HandleSubmit(e) {
    e.preventDefault()
    console.log('email:', email, 'password:', password)
    console.log('test', JSON.stringify({email, password}))
    try {
      const data = await login(email, password)
      console.log('response status:', data)
      navigate('/satellite')
    } catch (err) {
      console.log('response ERROR:', err)
      setError(err.message)
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
              required={true}
              onChange={(e)=>SetEmail(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="password">Пароль</label>
            <input type="password"
              value={password}
              required={true}
              onChange={(e)=>SetPwd(e.target.value)}
            />
          </div>
            <button type="submit">Войти</button>
          
        </form>
      {error && <p>{error}</p>}
      <Link to="/registration">Нет аккаунта? Зарегистрируйся!</Link>
    </div>
  )
}

export default Login