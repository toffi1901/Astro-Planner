import { Link } from "react-router";

function Home() {
  return (
    <div>
        <h1>Astro Planner</h1>
        <p>Поставьте Софе 5А, она заслужила</p>
        <Link to="/login">Попробуйте уже сейчас!</Link>
    </div>
  )
}

export default Home