import { useEffect } from "react"
import { useAppDispatch, useAppSelector } from "../store/hooks"
import { fetchLinks } from "../features/links/linksSlice"

const Dashboard = () => {
  const dispatch = useAppDispatch()
  const { links, loading, error } = useAppSelector((state) => state.links)

  useEffect(() => {
    dispatch(fetchLinks())
  }, [])

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Mes liens</h1>
      {loading && <p>Chargement...</p>}
      {error && <p className="text-red-500">{error}</p>}
      <ul>
        {links.map((link) => (
          <li key={link.id} className="mb-2">
            <a href={link.url} className="text-blue-600 underline" target="_blank">
              {link.title}
            </a>{" "}
            — Tags: {link.tags.join(", ")}
          </li>
        ))}
      </ul>
    </div>
  )
}

export default Dashboard

