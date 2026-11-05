import { Routes, Route, Navigate, Outlet } from 'react-router'
import { useAuth } from './auth.jsx'
import Home from './pages/home.jsx'
import Login from './pages/login.jsx'
import Satellite from './pages/satellite.jsx'
import Registration from './pages/registration.jsx'
import Profile from './pages/profile.jsx'

const PrivateRoutes = () => {
  const { isAuth, loading } = useAuth()
  if (loading) return <p>Шестеренки скрипят...</p>
  return (
    isAuth ? <Outlet/> : <Navigate to='/login'/>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />}/>
      <Route path="/login" element={<Login />}/>
      <Route path="/registration" element={<Registration />}/>

      <Route element={<PrivateRoutes/>}>
        <Route path="/satellite" element={<Satellite />}/>
        <Route path="/profile" element={<Profile />}/>
      </Route>
    </Routes>
  );
}

export default App