import { Routes, Route, Navigate, Outlet } from 'react-router'
import Home from './pages/home.jsx'
import Login from './pages/login.jsx'
import Satellite from './pages/satellite.jsx'

const PrivateRoutes = () => {
  let auth = {'token':true}
  return (
    auth.token ? <Outlet/> : <Navigate to='/login'/>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />}/>
      <Route path="/login" element={<Login />}/>

      <Route element={<PrivateRoutes/>}>
        <Route path="/satellite" element={<Satellite />}/>
      </Route>
    </Routes>
  );
}

export default App