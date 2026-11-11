import { createContext, useContext, useEffect, useState } from "react"
import { refresh, login as apiLogin } from './api.js'

const AuthContext = createContext(null)

function AuthProvider({ children }) {
  const [isAuth, setIsAuth] = useState(false)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    refresh().then(ok => { setIsAuth(ok); setLoading(false) })
  }, [])

  const login = async (email, password) => {
    await apiLogin(email, password)
    setIsAuth(true)
  }

  return (
    <AuthContext.Provider value={{ isAuth, loading, login }}>
      { children }
    </AuthContext.Provider>
  )
}

export default AuthProvider
export const useAuth = () => useContext(AuthContext)