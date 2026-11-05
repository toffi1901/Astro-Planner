import { Link, useNavigate } from "react-router";
import { useState } from "react";
import { reg } from "../api";

function Registration() {
  let navigate = useNavigate()
  
  const [email, SetEmail] = useState('')
  const [password, SetPwd] = useState('')
  const [password2, SetPwd2] = useState('')
  const [error, setError] = useState('')

  async function HandleSubmit(e) {
    e.preventDefault()
    console.log('email:', email, 'password:', password)
    console.log('test', JSON.stringify({email, password}))
    try {
      const data = await reg(email, password)
      console.log('response registration:', data.status)
      console.log('response registration c:', data)
      navigate("/login")
    } catch (err) {
      console.log('response ERROR registration:', err)
      setError(err.message)
    }
  }
  return (
    <div>
        <h1>Astro Planner регистрация</h1>
        <Link to="/">На главную</Link>
        <form onSubmit={HandleSubmit}>
          <div>
            <label htmlFor="email">Эл. почта</label>
            <input type="email"
              value={email}
              placeholder="Ваша почта"
              required={true}
              onChange={(e)=>SetEmail(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="password">Пароль</label>
            <input type="password"
              value={password}
              required={true}
              placeholder="Ваш пароль [a-zA-Z0-9]"
              onChange={(e)=>SetPwd(e.target.value)}
            />
          </div>
          <div>
            <label htmlFor="password2">Подтвердите пароль</label>
            {(password!==password2)&&<span className="error-lite">Password doesnt match</span>}
            <input
              type="password"
              value={password2}
              required={true}
              placeholder="Ваш пароль [a-zA-Z0-9]"
              onChange={(e)=> SetPwd2(e.target.value)}
            />
          </div>
          <button type="submit">Зарегистрироваться</button>
          
        </form>
      {error && <p>{error}</p>}
    </div>
  )
}

export default Registration